package button

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
  	"io/ioutil"
	"net/http"
	"net/url"
	"fmt"
)

type hnapireply struct {
	Hits int
    Results []Result
}

type Result struct {
    Item Hit
}

type Hit struct {
	Id int
	Points int
	Num_comments int
	Username string
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

    if hnreply.Hits == 0 {
    	// RENDER SUBMIT TEMPLATE
    	// panic("Nope!")
    	c.Infof("No hits")
    } else {
   		c.Infof("Hits: %f, Points: %f, ID: %i \n", hnreply.Hits, hnreply.Results[0].Item.Points, hnreply.Results[0].Item.Id)
    }

	//w.Header().Set("Location", req_url[0])
	w.WriteHeader(http.StatusFound)
  	w.Write(body)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/igrigorik", http.StatusFound)
}

func init() {
	http.HandleFunc("/button", Button)
	http.HandleFunc("/", Redirect)
}