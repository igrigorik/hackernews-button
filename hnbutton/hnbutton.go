package button

import (
    "appengine"
    "appengine/urlfetch"
    "encoding/json"
    "fmt"
    "html/template"
    "io/ioutil"
    "net/http"
    "net/url"
)

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

    if hnreply.Hits == 0 {
        c.Infof("No hits, rendering submit template")
        params := map[string]interface{}{"Url": req_url[0], "Title": req_title[0]}
        if err := buttonTemplate.Execute(w, params); err != nil {
            panic("Cannot execute template")
        }

    } else {
        c.Infof("Hits: %f, Points: %f, ID: %i \n",
            hnreply.Hits,
            hnreply.Results[0].Item.Points,
            hnreply.Results[0].Item.Id)

        if err := buttonTemplate.Execute(w, hnreply.Results[0].Item); err != nil {
            panic("Cannot execute template")
        }
    }

}

var buttonTemplate = template.Must(template.New("button").Parse(buttonTemplateHTML))

const buttonTemplateHTML = `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">

<head>
<style type="text/css">
.hn_bt {font-family:'lucida grande',verdana,sans-serif;font-size:11px;text-decoration:none;cursor:pointer;line-height:14px;display:block;padding:0;margin:0;}.hn_cn{float:left;position:relative;z-index:2;height:0;width:5px;top:-5px;left:2px;}
.hn_cn s,.hn_cn i {border:solid transparent;border-right-color:#d7d7d7;border-width:4px 5px 4px 0;top:1px;display:block;position:relative;}
.hn_cn i {left:2px;top:-7px;border-right-color:#fff;}
.hn_cc{background:#fff;border:1px solid #d1d1d1;float:left;font-weight:normal;height:14px;margin-left:1px;min-width:15px;padding:1px 2px 1px 2px;text-align:center;line-height:14px;}
.hn_a {padding:2px;border:1px solid;border-radius:3px;border-color:#ffdac2;background-color:#ff6500;display:block;color:white;}.hn_lo{width:25px;height:14px;margin-left:-1px;padding-left:17px;}
.hn_s {width:10px;height:14px;padding-left:13px;margin-left:1px;}.hn_pi{background:url("data:image/gif;base64,R0lGODlhDAAMAKIAAPihV/aKLf/jyvmxc/////9mAAAAAAAAACH5BAAAAAAALAAAAAAMAAwAAAMiWLrMRGE9NshQ9TqhBGkFYFUAWHgEZ1aWqUxuAbuzWcdMAgA7") 0 1px no-repeat;}
.hn_ar {background:url("data:image/gif;base64,R0lGODlhCgAKALMJANPT06enp/b29r+/v52dnfn5+bq6usLCwpqamv///wAAAAAAAAAAAAAAAAAAAAAAACH5BAEAAAkALAAAAAAKAAoAAAQcMMlJq712GIzQDV1QFV0nUETZTYDaAdIhz3ISAQA7") 1px 0px no-repeat;}
.hn-button:before { content:'Y';color: #FFF;background: #ff6600;width: 12px;height:12px;margin-top: 3px;padding-right: 3px;margin-right:4px;padding-left: 3px;}
.hn-button {
    -webkit-border-radius: 3px;
    -moz-border-radius: 3px;
    border-radius: 3px;
    color: #ff6600;
    text-decoration: none;
    background: rgb(250,250,250);
    background: -moz-linear-gradient(top, rgba(250,250,250,1) 0%, rgba(234,234,234,1) 100%);
    background: -webkit-gradient(linear, left top, left bottom, color-stop(0%,rgba(250,250,250,1)), color-stop(100%,rgba(234,234,234,1)));
    background: -webkit-linear-gradient(top, rgba(250,250,250,1) 0%,rgba(234,234,234,1) 100%);
    background: -o-linear-gradient(top, rgba(250,250,250,1) 0%,rgba(234,234,234,1) 100%);
    background: -ms-linear-gradient(top, rgba(250,250,250,1) 0%,rgba(234,234,234,1) 100%);
    filter: progid:DXImageTransform.Microsoft.gradient( startColorstr='#fafafa', endColorstr='#eaeaea',GradientType=0 );
    background: linear-gradient(top, rgba(250,250,250,1) 0%,rgba(234,234,234,1) 100%);
    font: 12px "Helvetica Neue", Arial, Helvetica, Geneva, sans-serif;
    border: 1px solid #d2d0ce;

    padding: 2px 4px 2px 4px;

    -webkit-box-shadow: inset 0px 1px 0px 0px rgba(255, 255, 255, 1);
    -moz-box-shadow: inset 0px 1px 0px 0px rgba(255, 255, 255, 1);
    box-shadow: inset 0px 1px 0px 0px rgba(255, 255, 255, 1);
}
.hn-button:hover {border: 1px solid rgba(66,66,66,0.52);}
.hn-button:active {
    background: rgb(234,234,234);
    background: -moz-linear-gradient(top, rgba(234,234,234,1) 0%, rgba(250,250,250,1) 100%);
    background: -webkit-gradient(linear, left top, left bottom, color-stop(0%,rgba(234,234,234,1)), color-stop(100%,rgba(250,250,250,1)));
    background: -webkit-linear-gradient(top, rgba(234,234,234,1) 0%,rgba(250,250,250,1) 100%);
    background: -o-linear-gradient(top, rgba(234,234,234,1) 0%,rgba(250,250,250,1) 100%);
    background: -ms-linear-gradient(top, rgba(234,234,234,1) 0%,rgba(250,250,250,1) 100%);
    filter: progid:DXImageTransform.Microsoft.gradient( startColorstr='#eaeaea', endColorstr='#fafafa',GradientType=0 );
    background: linear-gradient(top, rgba(234,234,234,1) 0%,rgba(250,250,250,1) 100%);
    border: 1px solid rgba(255,102,0,0.74);
}
</style>
</head>

<body style="background:transparent;">
  <div>
  {{if .Id}}
    <a href="http://news.ycombinator.com/item?id={{.Id}}" target="_blank">
  {{else}}
    <a href="http://news.ycombinator.com/submitlink?u={{.Url}}&t={{.Title}}">
  {{end}}
    <span class="hn_bt">
    <table cellspacing="0" cellpadding="0">
      <tbody>
        <tr>
          <td><a class="hn_a"><span class="hn_s hn_pi"></span></a></td>
          <td>
            <table cellspacing="0" cellpadding="0">
              <tbody>
                <tr>
                  <td><div class="hn_cn"><s></s><i></i></div></td>
                  <td><div class="hn_cc"><span style="color:#333;">
                  {{if .Id}}
                    {{.Points}}
                  {{else}}
                    submit
                  {{end}}
                  </span></div></td>
                </tr>
              </tbody>
            </table>
          </td>
        </tr>
      </tbody>
    </table>
    </span>
  </a>
  </div>
</body>
</html>
`

func Redirect(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://github.com/igrigorik", http.StatusFound)
}

func init() {
    http.HandleFunc("/button", Button)
    http.HandleFunc("/", Redirect)
}
