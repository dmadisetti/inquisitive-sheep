package shoop

import "appengine/datastore"

// Kill it
func Delete(table string, session Session) bool{
    q := datastore.NewQuery(table)
    keys, err := q.KeysOnly().GetAll(session.Context,nil)
    if err == nil {
        err := datastore.DeleteMulti(session.Context, keys)
        if err != nil{
            session.Context.Warningf("Badness: %v", err)
            return false
        }
    }else {
        return false
    }
    return true
}