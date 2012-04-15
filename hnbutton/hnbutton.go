package button

import (
    "appengine"
    "appengine/urlfetch"
    "encoding/json"
    "html/template"
    "io/ioutil"
    "net/http"
    "net/url"
    "fmt"
)

var buttonTemplate, _ = template.New("page").ParseFiles("hnbutton/button.html")

type hnapireply struct {
    Hits    int
    Results []Result
}

type Result struct {
    Item Hit
}

type Hit struct {
    Id           int
    Points       int
    Num_comments int
    Username     string
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

    c.Infof("Fetching HN data for: %s, %s\n", req_title, req_url)

    _, err := url.Parse(req_url[0])
    if err != nil {
        panic("Invalid URL: " + err.Error())
    }

    pageData := "http://api.thriftdb.com/api.hnsearch.com/items/_search?filter[fields][url][]=" + req_url[0]

    client := urlfetch.Client(c)
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

    // Cache the response in the HTTP edge cache, if possible
    // http://code.google.com/p/googleappengine/issues/detail?id=2258
    w.Header().Set("Cache-Control", "public, max-age=61;")

    if hnreply.Hits == 0 {
        c.Infof("No hits, rendering submit template")
        params := map[string]interface{}{"Url": req_url[0], "Title": req_title[0]}
        if err := buttonTemplate.ExecuteTemplate(w, "button", params); err != nil {
            panic("Cannot execute template")
        }

    } else {
        c.Infof("Hits: %f, Points: %f, ID: %i \n",
            hnreply.Hits,
            hnreply.Results[0].Item.Points,
            hnreply.Results[0].Item.Id)

        if err := buttonTemplate.ExecuteTemplate(w, "button", hnreply.Results[0].Item); err != nil {
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
