package shoop

import (
	"encoding/json"
	"time"
)

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