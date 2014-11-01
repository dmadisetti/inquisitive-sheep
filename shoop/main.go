package shoop

import(
    "net/http"
    "log"
    "fmt"
    "time"
    "reflect"
    "strings"
    "encoding/json"
    "html/template"
    "net/url"
    "strconv"
    "crypto/sha1"
    "appengine"
    "appengine/urlfetch"
    "appengine/datastore"
    "appengine/memcache"
    "encoding/base64"
)

const protocol = "https:"
const userAgent = "Dalvik/1.6.0 (Linux; U; Android 4.4.2; SPH-L710 Build/KOT49H)"
const connection = "Keep-Alive"
const acceptEncoding = "gzip"

var t *template.Template
var c appengine.Context
var client *http.Client
var get = url.URL{
    Host:"forgot-to-close-my-api.net",
    Path:"/api/getMessages",
}
var params = url.Values{
    "lat":{"33.9800272"},
    "long":{"-80.9842865"},
    "userID":{"ddc6728f-bd0d-11e3-ab67-0401130aa601"},
    "version":{"2.1.003"},
}

type Settings struct {
    Start time.Time
    Long string
    Lat string
    Error int
    Fatal bool
    Over bool
    Name string
    Host string
}

func defaultSettings() *Settings{   
    return &Settings{
        Long: "-80.9842865",
        Lat: "33.9800272",
        Fatal: false,
        Over: true,
        Name: "USC",        
        Host: "forgot-to-close-my-api.net",
    }
}

func save(settings *Settings){
    datastore.Put(c,datastore.NewKey(c,"Settings","",1, nil),settings)
}

// Structs to match json
type Message struct {
    MessageID string
    Message string
    NumberOfLikes int
    PosterID string
    Handle string
    Time string
    FormattedTime time.Time
}
type Package struct {
    Messages [] Message
}
// Struct to store time data
type Instance struct {
    MessageID string
    NumberOfLikes int
    Time time.Time
}

// Json hack from some blog (will update if I find again)
type Response map[string]interface{}
func (r Response) String() (s string) {
    b, err := json.Marshal(r)
    if err != nil {
            s = ""
            return
    }
    s = string(b)
    return
}

// Write out as CSV
type File struct {
    table reflect.Type
    tables reflect.Value
}
func (f File) String() (s string) {
    delimiter := "|"
    lbreak  := "\n"

    // Set headers
    for i := 0; i < f.table.NumField(); i++ {
        s += "\"" + f.table.Field(i).Name + "\""
        if i + 1 < f.table.NumField(){
            s += delimiter
        }else{
            s += lbreak
        }
    }

    // Populate fields
    var field reflect.Value
    for i := 0; i < f.tables.Len(); i++ {
        for j := 0; j < f.table.NumField(); j++ {
            field = f.tables.Index(i).Field(j)

            s += "\""
            switch field.Kind() {
                case reflect.String:
                    s += strings.Replace(field.String(),"\"","'",-1)
                    break;
                case reflect.Int:
                    s += strconv.FormatInt(field.Int(),10)
                    break
                case reflect.ValueOf(time.Time{}).Kind():
                    s += strconv.FormatInt(field.MethodByName("Unix").Call([]reflect.Value{})[0].Int(),10)
                    break
            }
            s += "\""

            if j + 1 < f.table.NumField(){
                s += delimiter
            }else{
                s += lbreak
            }
        }
    }
    return
}

type Handler struct{
    handle func(http.ResponseWriter, *http.Request, *Settings)
}
func (handle *Handler)preHandle(w http.ResponseWriter, r *http.Request){
    c = appengine.NewContext(r)

    settings := defaultSettings()
    err := datastore.Get(c,datastore.NewKey(c,"Settings","",1, nil),settings)
    if err !=nil{
        log.Println(err)
        save(settings)
    }

    if client == nil{
        client = urlfetch.Client(c)
    }
    handle.handle(w,r, settings)
}
func newHandle(path string, handle func(http.ResponseWriter, *http.Request, *Settings)) {
    handler := &Handler{handle:handle}
    http.HandleFunc(path, handler.preHandle)
}


// Start er up!
func init(){
    log.Println("Init")
    newHandle("/", mainHandle)
    newHandle("/run", runHandle)
    newHandle("/flush", flushHandle)
    newHandle("/instance.csv", instanceHandle)
    newHandle("/message.csv", messageHandle)
}

// Handles
func mainHandle(w http.ResponseWriter, r *http.Request, settings *Settings){
    t, e := template.ParseGlob("templates/the.html")
    if e != nil {
        fmt.Fprint(w, e)        
        return
    }
    err := t.Execute(w, settings)
    if err !=nil{
        panic(err)
    }
}

func runHandle(w http.ResponseWriter, r *http.Request, settings *Settings){
    defer func() {
        if r := recover(); r != nil {
            settings.Error += 1
            if settings.Error >= 5{
                settings.Fatal = true;
            }
            save(settings)
        }else{
            if -settings.Start.Sub(time.Now()).Hours() > 24{
                settings.Over = true;
            }
            save(settings)
        }
    }()
    if settings.Start.Sub(time.Now()).Hours() > 0{
        fmt.Fprint(w,  Response{"Woops":"Hasn't started"})
        return
    }
    if settings.Over || settings.Fatal {
        fmt.Fprint(w,  Response{"Woops":"Over"})
        return
    }
    params.Set("long",settings.Long)
    params.Set("lat",settings.Lat)
    get.Host = settings.Host
    now := time.Now()
    t := strconv.FormatInt(now.Unix(),10)
    hash := (sha1.Sum([]byte(t)))
    params.Add("salt",t)
    params.Add("hash",base64.StdEncoding.EncodeToString(hash[:]))
    request,err := http.NewRequest("GET", protocol + get.String() + "?" + params.Encode(), nil)
    if err != nil {
        panic(err)
        return
    }
    request.Header.Set("User-Agent", userAgent)
    request.Header.Set("Connection", connection)
    request.Header.Set("Accept-Encoding", acceptEncoding)
    response,err := client.Do(request)
    if err != nil {
        panic(err)        
        return
    }
    //Decode request
    var p Package
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&p)
    if err != nil {
        panic(err)
    }

    var item *memcache.Item
    var instance Instance
    var message Message
    var id string
    for i := 0; i < len(p.Messages); i++ {
        message = p.Messages[i]
        id = message.MessageID
        instance = Instance{
            MessageID: id,
            NumberOfLikes: message.NumberOfLikes,    
            Time: now,
        }
        item = &memcache.Item{
            Key:   id,
            Value: make([]byte,1),
        }
        // If not in memcache, Create in datastore
        if err = memcache.Add(c, item); err == nil {
            // Doesn't matter if fails, nothing we can do
            if message.FormattedTime,err = time.Parse("2006-01-02 15:04:05", message.Time); err == nil{
                datastore.Put(c,datastore.NewKey(c,"Message",id,0, nil),&message)
            }
        }
        datastore.Put(c,datastore.NewIncompleteKey(c,"Instance",nil),&instance)
    }

    fmt.Fprint(w,  Response{"ran":true})
}

func delete(table string) bool{
    q := datastore.NewQuery(table)
    keys, err := q.KeysOnly().GetAll(c,nil)
    if err == nil {
        err := datastore.DeleteMulti(c, keys)
        if err != nil{
            return false
        }
    }else {
        return false
    }
    return true
}

func flushHandle(w http.ResponseWriter, r *http.Request, settings *Settings){
    message := "Not done yet"
    if settings.Over{
        message = "Done thanks"
        if !(delete("Instance") && delete("Message")){
            message = "Broke"
        }
    }
    fmt.Fprint(w,  Response{"Thanks":message})
}

func writeCSV(w http.ResponseWriter, table reflect.Type, tables reflect.Value, settings *Settings){
    if !settings.Over{
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, Response{"Woops":"Not done yet"})
        return
    }
    w.Header().Set("Content-Type", "text/csv")
    fmt.Fprint(w, File{table:table,tables:tables})
}

func instanceHandle(w http.ResponseWriter, r *http.Request, settings *Settings){
    var tables []Instance

    // Look it all up
    keys,err := datastore.NewQuery("Instance").GetAll(c,&tables)
    if err != nil  || len(keys) == 0{
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, Response{"Woops":"No data"})
        return
    }

    writeCSV(w,reflect.TypeOf(Instance{}),reflect.ValueOf(tables),settings)
}

func messageHandle(w http.ResponseWriter, r *http.Request, settings *Settings){
    var tables []Message

    // Look it all up
    keys,err := datastore.NewQuery("Message").GetAll(c,&tables)
    if err != nil  || len(keys) == 0{
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, Response{"Woops":"No data"})
        return
    }

    writeCSV(w,reflect.TypeOf(Message{}),reflect.ValueOf(tables),settings)
}