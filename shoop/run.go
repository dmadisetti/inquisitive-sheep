package shoop

import "time"

func Catch(session Session){
    if r := recover(); r != nil {
        session.Context.Warningf("Badness: %v", r)
        session.Settings.Error += 1
        if session.Settings.Error >= 5{
            session.Settings.Fatal = true;
            session.Settings.Over = true;
        }
        session.Save()
    }else{
        if -session.Settings.Start.Sub(time.Now()).Hours() > 24{
            session.Settings.Over = true;
        }
        session.Save()
    }
}