package shoop

import (
    "net/http"
    "log"
    "appengine"
    "appengine/urlfetch"
    "appengine/datastore"
)

// Handler just to hold everything together
type Handler struct{
    handle func(http.ResponseWriter, *http.Request, Session)
}

// Constructor
func NewHandler(path string, handle func(http.ResponseWriter, *http.Request, Session)) {
    handler := &Handler{handle:handle}
    http.HandleFunc(path, handler.preHandle)
}

// Passed into http for all handlers
func (h *Handler)preHandle(w http.ResponseWriter, r *http.Request){
    c := appengine.NewContext(r)
    // Create session
    session := Session{
        Context: c,
        Client: urlfetch.Client(c),
        Settings: NewSettings(),
    }

    // Set transport to allow for https
    session.Client.Transport = &urlfetch.Transport{
        Context:                       c,
        Deadline:                      0,
        AllowInvalidServerCertificate: false,
    }

    // get session from data store or create
    err := datastore.Get(session.Context,datastore.NewKey(session.Context,"Settings","",1, nil),session.Settings)
    if err !=nil{
        log.Println(err)
        session.Save()
    }

    // Call handler set earlier
    h.handle(w,r, session)
}