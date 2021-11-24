package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Conf struct {
	At      map[string]string
	Url     string
	Keyword string
}

var conf = Conf{}

func init() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	json, err := json.Marshal(conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(string(json))
}

func dealGitlab(writer http.ResponseWriter, request *http.Request) {
	bytes, _ := ioutil.ReadAll(request.Body)
	data := string(bytes)
	fmt.Println("receive:" + data)
	kind := gjson.Get(data, "object_kind").Str
	if "issue" != kind {
		return
	}
	action := getAction(data)
	if action == 0 {
		return
	}
	title := gjson.Get(data, "object_attributes.title").Str
	giturl := gjson.Get(data, "object_attributes.url").Str
	if action == 1 {
		title = "[new bug] " + title
		fmt.Println(title)
	} else {
		title = "[bug 待回归] " + title
		fmt.Println(title)
	}

	ats := getAt(data, action)
	postUrl := conf.Url
	body := getBody(title, giturl, ats)

	req, err := http.NewRequest("POST", postUrl, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func getAt(data string, action int) string {
	prjName := gjson.Get(data, "project.name").Str

	if action == 2 {
		prjName = "test-group"
	}

	for k, v := range conf.At {
		if k == prjName {
			return v
		}
	}
	return ""
}

//
func getBody(title, url, ats string) string {
	body := `{
    "msgtype": "text", 
    "text": {
        "content": "_TITLE_(_KEYWORD_)\r_URL_"
    },
    "at": {
        "atMobiles": [
            `
	for i, at := range strings.Split(ats, ",") {
		body += "\"" + at + "\""
		if i < len(strings.Split(ats, ","))-1 {
			body += ","
		}
	}
	body += `
        ], 
        "isAtAll":false
    }
}`
	body = strings.ReplaceAll(body, "_TITLE_", title)
	body = strings.ReplaceAll(body, "_URL_", url)
	body = strings.ReplaceAll(body, "_KEYWORD_", conf.Keyword)

	fmt.Println("return body:" + body)
	return body
}

// 0过滤, 1打开, 2回归
func getAction(data string) int {
	state := gjson.Get(data, "object_attributes.action").Str
	if "close" == state {
		return 0
	}

	// 是否有回归
	rs := gjson.Get(data, "object_attributes.labels").Array()
	if len(rs) > 0 {
		for _, r := range rs {
			title := r.Get("title").Str
			if strings.Contains(title, "回归") {
				return 2
			}
		}
	}

	if "open" == state || "reopen" == state {
		return 1
	}
	return 0
}

func main() {
	http.HandleFunc("/", dealGitlab)
	http.ListenAndServe(":11223", nil)
	fmt.Println("listen:11223")
}
