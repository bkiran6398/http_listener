package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var counter int
var counter_mtx sync.RWMutex = sync.RWMutex{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		counter_mtx.Lock()
		counter++
		fmt.Printf("=%d===================================================================%s=\n", counter, time.Now().Format("2006-01-02 15:04:05"))
		counter_mtx.Unlock()
		fmt.Println("Proto:", r.Proto)
		fmt.Println("Method:", r.Method)
		fmt.Printf("URL: %s%s\n", r.Host, r.URL)
		fmt.Println("RequestFrom:", r.RemoteAddr)

		fmt.Println("-------------------------------Headers-")
		for k, v := range r.Header {
			fmt.Printf("%s: %+v\n", k, v)
		}

		fmt.Println("-------------------------------Body-")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("error reading body: ", err)
		}
		r.Body.Close()
		var b bytes.Buffer
		if err := json.Indent(&b, body, "", "\t"); err != nil {
			log.Fatal("error marshaling body: ", err)
		}
		fmt.Println(b.String())

		fmt.Println("-------------------------------FormValues-")
		if err := r.ParseForm(); err != nil {
			log.Fatal("failed to parse form: ", err)
		}
		for k, v := range r.Form {
			fmt.Printf("%s: %+v\n", k, v)
		}

		fmt.Println("-------------------------------PostFormValues-")
		for k, v := range r.PostForm {
			fmt.Printf("%s: %+v\n", k, v)
		}

		fmt.Println("-------------------------------MultipartForm-")
		if r.Header.Get("Content-Type") == "multipart/form-data" {
			if err := r.ParseMultipartForm(1024); err != nil {
				log.Fatal("failed to parse multipart form: ", err)
			}
			for k, v := range r.MultipartForm.Value {
				fmt.Printf("%s: %+v", k, v)
			}
		}

		if r.Response != nil {
			fmt.Println("---------------------------Response-")
			fmt.Println(r.Response)
		}
	})
	fmt.Println("Listening on http://localhost:10101")
	if err := http.ListenAndServe(":10101", nil); err != nil {
		log.Fatal("error starting server: ", err)
	}
}
