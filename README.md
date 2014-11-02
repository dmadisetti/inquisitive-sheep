Inquisitive Sheep - Someone left their API open. Uhoh
=========

[![Build Status](https://travis-ci.org/dmadisetti/inquisitive-sheep.png)](https://travis-ci.org/dmadisetti/inquisitive-sheep)

Put together over the night. Starting to like Go, but still at the point where the API is my best friend. 

Data collection for mobile app api. Not discussing which app incase I get sued (benign script, but still). You might be able to guess what service if you checkout the structs for json. Cool exercise in packet sniffing and investigation

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
Inteferface for status and downloading aggregate data after run
![alt tag](https://raw.github.com/dmadisetti/inquisitive-sheep/master/screenshot.jpg "Screenshot")

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
