package shoop

import (
    "net/http"
    "time"
    "appengine"
    "appengine/datastore"
)

type Session struct {
    Context appengine.Context
    Client *http.Client
    Settings *Settings
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

func (session *Session) Save(){
    datastore.Put(session.Context,datastore.NewKey(session.Context,"Settings","",1, nil),session.Settings)
}

func (session *Session) Check() bool{
    if session.Settings.Start.Sub(time.Now()).Hours() > 0{
        return false
    }
    if session.Settings.Over || session.Settings.Fatal {        
        return false
    }
    return true
}

func NewSettings()*Settings{   
    return &Settings{
        Long: "-80.9842865",
        Lat: "33.9800272",
        Fatal: false,
        Over: true,
        Name: "USC",        
        Host: "forgot-to-close-my-api.net",
    }
}
