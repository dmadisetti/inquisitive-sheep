package shoop

import (
    "net/http"
    "fmt"
    "time"
    "reflect"
    "encoding/json"
    "html/template"
    "net/url"
    "strconv"
    "crypto/sha1"
    "appengine/datastore"
    "appengine/memcache"
    "encoding/base64"
)

// Headers stuffs
const protocol = "https:"
const userAgent = "Dalvik/1.6.0 (Linux; U; Android 4.4.2; SPH-L710 Build/KOT49H)"
const connection = "Keep-Alive"
const acceptEncoding = "gzip"

var t *template.Template
var get = url.URL{
    Host:"forgot-to-close-my-api.net",
    Path:"/api/getMessages",
}

// Start er up!
func init(){
    NewHandler("/", mainHandle)
    NewHandler("/run", runHandle)
    NewHandler("/flush", flushHandle)
    NewHandler("/instance.csv", instanceHandle)
    NewHandler("/message.csv", messageHandle)
}

// Handles
func mainHandle(w http.ResponseWriter, r *http.Request, session Session){
    t, err := template.ParseGlob("templates/the.html")
    if err != nil {
        fmt.Fprint(w, err)        
        return
    }
    err = t.Execute(w, session.Settings)
    if err !=nil{
        panic(err)
    }
}

func runHandle(w http.ResponseWriter, r *http.Request, session Session){
    defer Catch(session)

    // Make sure no fatal and good to go
    if !session.Check(){
        fmt.Fprint(w,  Response{"Woops":"Has't run"})
        return
    }

    // Grab time and hash
    now := time.Now()
    t := strconv.FormatInt(now.Unix(),10)
    hash := (sha1.Sum([]byte(t)))
    params := url.Values{
        "lat":{session.Settings.Lat},
        "long":{session.Settings.Long},
        "userID":{"ddc6728f-bd0d-11e3-ab67-0401130aa601"},
        "version":{"2.1.003"},
        "salt": {t},
        "hash": {base64.StdEncoding.EncodeToString(hash[:])},
    }

    // Set host, headers and call
    get.Host = session.Settings.Host
    request,err := http.NewRequest("GET", protocol + get.String() + "?" + params.Encode(), nil)
    if err != nil {
        panic(err)
        return
    }
    request.Header.Set("User-Agent", userAgent)
    request.Header.Set("Connection", connection)
    request.Header.Set("Accept-Encoding", acceptEncoding)
    response,err := session.Client.Do(request)
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
    // For message in messages save if not in memcache
    // Create instance for each call
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
        if err = memcache.Add(session.Context, item); err == nil {
            // Doesn't matter if fails, nothing we can do
            if message.FormattedTime,err = time.Parse("2006-01-02 15:04:05", message.Time); err == nil{
                datastore.Put(session.Context,datastore.NewKey(session.Context,"Message",id,0, nil),&message)
            }
        }
        datastore.Put(session.Context,datastore.NewIncompleteKey(session.Context,"Instance",nil),&instance)
    }

    fmt.Fprint(w,  Response{"ran":true})
}

// Kill everything and clear errors
func flushHandle(w http.ResponseWriter, r *http.Request, session Session){
    message := "Not done yet"
    if session.Settings.Over{
        message = "Done thanks"
        if !(Delete("Instance",session) && Delete("Message",session)){
            message = "Broke"
        }
        session.Settings.Fatal = false
        session.Settings.Error = 0
        session.Save()
    }
    fmt.Fprint(w,  Response{"Thanks":message})
}

// CSV handlers. 
// Redundant, but reflection gets weird 
// without an initial concrete value
// Write out instances to csv
func instanceHandle(w http.ResponseWriter, r *http.Request, session Session){
    var tables []Instance
    // Look it all up
    keys,err := datastore.NewQuery("Instance").GetAll(session.Context,&tables)
    if err != nil  || len(keys) == 0{
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, Response{"Woops":"No data"})
        session.Context.Warningf("Datasotre fail: %v", err)
        return
    }
    WriteCSV(w,reflect.TypeOf(Instance{}),reflect.ValueOf(tables),session)
}

// Write out messages to csv
func messageHandle(w http.ResponseWriter, r *http.Request, session Session){
    var tables []Message
    // Look it all up
    keys,err := datastore.NewQuery("Message").GetAll(session.Context,&tables)
    if err != nil  || len(keys) == 0{
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, Response{"Woops":"No data"})
        session.Context.Warningf("Datasotre fail: %v", err)
        return
    }
    WriteCSV(w,reflect.TypeOf(Message{}),reflect.ValueOf(tables),session)
}