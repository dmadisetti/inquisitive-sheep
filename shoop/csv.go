package shoop

import (
    "net/http"
    "fmt"
    "time"
    "reflect"
    "strings"
    "strconv"
)

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

// Spit out datastore with reflection
func WriteCSV(w http.ResponseWriter, table reflect.Type, tables reflect.Value, session Session){
    if !session.Settings.Over{
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, Response{"Woops":"Not done yet"})
        return
    }
    w.Header().Set("Content-Type", "text/csv")
    fmt.Fprint(w, File{table:table,tables:tables})
}