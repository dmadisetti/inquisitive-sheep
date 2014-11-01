Inquisitive Sheep - Someone left their API open. Uhoh
=========

[![Build Status](https://travis-ci.org/dmadisetti/inquisitive-sheep.png)](https://travis-ci.org/dmadisetti/inquisitive-sheep)

Put together over the night. Starting to like go, but still athe point where the API is my best friend. Data collection for app api. Not including app incase I get sued. Benign script, but still. Might be able to guess what service if you checkout the structs for json.

---
Manage settings from the GAE datastore viewer

```
    Start time.Time // Jobs run for 24 from this point 
    Long string // Longitude
    Lat string // Latitude
    Error int // Error count incase cron goes sour
    Fatal bool // Cron go sour?
    Over bool // 24 period over
    Name string // Just for visual recognition
    Host string // api url
```

---

Objectives:

- More Go - `Check`
- Lots of Data - `Check`
- Witty remarks - `Check`

Todo:

- Tests
- Clean!!!!
- Learn from mistakes
- Document???
