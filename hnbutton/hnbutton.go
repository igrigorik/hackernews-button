package button

import (
    "appengine"
    "appengine/urlfetch"
    "appengine/memcache"
    "encoding/json"
    "html/template"
    "crypto/md5"
    "io/ioutil"
    "net/http"
    "net/url"
    "time"
    "hash"
    "fmt"
)

var buttonTemplate, _ = template.New("page").ParseFiles("hnbutton/button.html")

type hnapireply struct {
    NbHits int
    Hits []Hit
}

type Hit struct {
    Story_id     int
    Points       int
    Hits         int
    Num_comments int
    Author     string
}

func Button(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    c.Infof("Requested URL: %v", r.URL)

    defer func() {
        if err := recover(); err != nil {
            c.Errorf("%s", err)
            c.Errorf("%s", "Traceback: %s", r)

            reply := map[string]string{
                "error": fmt.Sprintf("%s", err),
            }

            resp, _ := json.Marshal(reply)
            w.WriteHeader(500)
            w.Write(resp)
        }
    }()

    query, _ := url.ParseQuery(r.URL.RawQuery)
    req_url, ok_url := query["url"]
    req_title, ok_title := query["title"]

    if !ok_url || !ok_title {
        panic("required parameters: url, title")
    }

    _, err := url.Parse(req_url[0])
    if err != nil {
        panic("Invalid URL: " + err.Error())
    }

    var h hash.Hash = md5.New()
    h.Write([]byte(req_url[0]))
    var hkey string = fmt.Sprintf("%x", h.Sum(nil))

    c.Infof("Fetching HN data for: %s, %s\n", req_title, req_url)

    var item Hit
    if cachedItem, err := memcache.Get(c, hkey); err == memcache.ErrCacheMiss {
        pageData := "http://hn.algolia.com/api/v1/search_by_date?query=" + url.QueryEscape(req_url[0])

        client := &http.Client{
            Transport: &urlfetch.Transport{
            Context: c,
            Deadline: time.Duration(15)*time.Second,
          },
        }

        resp, err := client.Get(pageData)
        if err != nil {
            panic("Cannot fetch HN data: " + err.Error())
        }

        defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        var hnreply hnapireply
        if err := json.Unmarshal(body, &hnreply); err != nil {
            panic("Cannot unmarshall JSON data")
        }

        if hnreply.NbHits == 0 {
            item.Hits = 0
        } else {
            item.Hits = hnreply.NbHits;
            item.Story_id = hnreply.Hits[0].Story_id;
            item.Points = hnreply.Hits[0].Points;
            item.Num_comments = hnreply.Hits[0].Num_comments;
            item.Author = hnreply.Hits[0].Author;
        }


        var sdata []byte
        if sdata, err = json.Marshal(item); err != nil {
            panic("Cannot serialize hit to JSON")
        }

        c.Debugf("Saving to memcache: %s", sdata)

        data := &memcache.Item{
            Key: hkey,
            Value: sdata,
            Expiration: time.Duration(60)*time.Second,
        }

        if err := memcache.Set(c, data); err != nil {
            c.Errorf("Cannot store hit to memcache: %s", err.Error())
        }

    } else if err != nil {
        panic("Error getting item from cache: %v")

    } else {
        if err := json.Unmarshal(cachedItem.Value, &item); err != nil {
            panic("Cannot unmarshall hit from cache")
        }
        c.Infof("Fetched from memcache: %i", item.Story_id)
    }

    // Cache the response in the HTTP edge cache, if possible
    // http://code.google.com/p/googleappengine/issues/detail?id=2258
    w.Header().Set("Cache-Control", "public, max-age=300")
    w.Header().Set("Pragma", "Public")

    if item.Hits == 0 {
        c.Infof("No hits, rendering submit template")
        params := map[string]interface{}{"Url": req_url[0], "Title": req_title[0]}
        if err := buttonTemplate.ExecuteTemplate(w, "button", params); err != nil {
            panic("Cannot execute template")
        }

    } else {
        c.Infof("Points: %f, ID: %i \n", item.Points, item.Story_id)

        if err := buttonTemplate.ExecuteTemplate(w, "button", item); err != nil {
            panic("Cannot execute template")
        }
    }
}

func Redirect(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://github.com/igrigorik/hackernews-button", http.StatusFound)
}

func init() {
    http.HandleFunc("/button", Button)
    http.HandleFunc("/", Redirect)
}
