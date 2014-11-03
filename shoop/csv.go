package shoop

import (
    "net/http"
    "fmt"
    "time"
    "reflect"
    "strings"
    "strconv"
    "appengine/datastore"
)

// Write out as CSV
type File struct {
    table string
    session Session
}
func (f File) String() (s string) {
    delimiter := "|"
    lbreak  := "\n"

    table := f.getReflection()

    // Set headers
    for i := 0; i < table.NumField(); i++ {
        s += "\"" + table.Field(i).Name + "\""
        if i + 1 < table.NumField(){
            s += delimiter
        }else{
            s += lbreak
        }
    }

    iterator := datastore.NewQuery(f.table).Run(f.session.Context)

    // Populate fields
    var field reflect.Value
    for {

        // Grab next if next
        tables,err := f.next(iterator)
        if err == datastore.Done {
                break
        }

        for j := 0; j < table.NumField(); j++ {
            field = tables.Field(j)

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

            if j + 1 < table.NumField(){
                s += delimiter
            }else{
                s += lbreak
            }
        }
    }
    return
}

func (f File) getReflection() (v reflect.Type){
    switch f.table{ 
        case "Message":
            v = reflect.TypeOf(Message{})
        default:
            v = reflect.TypeOf(Instance{})
    }
    return
}

func (f File) next(iterator *datastore.Iterator) (v reflect.Value, err error){
    switch f.table{ 
        case "Message":
            var group Message
            _,err = iterator.Next(&group)
            v = reflect.ValueOf(group)
            break
        default:
            var group Instance
            _,err = iterator.Next(&group) 
            v = reflect.ValueOf(group)
    }
    return
}

// Spit out datastore with reflection
func WriteCSV(w http.ResponseWriter, table string, session Session){
    if !session.Settings.Over{
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, Response{"Woops":"Not done yet"})
        return
    }
    w.Header().Set("Content-Type", "text/csv")
    fmt.Fprint(w, File{table:table,session:session})
}