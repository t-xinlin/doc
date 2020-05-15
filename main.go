package main

import (
	"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"log"
	"net/http"
)

// Version is the build version
var Version string

// GitTag is the git tag of the build
var GitTag string

// BuildDate is the date when the build was created
var BuildDate string

func main() {
	cast.ToBool(true)
	fmt.Printf("Version: %s, GitTag: %s, BuildDate: %s\n", Version, GitTag, BuildDate)
	http.HandleFunc("/spans", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("ReadAll error: %+v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		if r != nil {
			defer r.Body.Close()
		}
		log.Printf(">>>>>>>>>>>>>>>>>>>>>%+v", string(bytes))
		var status int = 200
		w.WriteHeader(status)
		log.Printf("Rec: : %+v  -> Back[  version: 1.0.0   httpStatusCode: %+v]", string(bytes), status)
		w.Write([]byte("Svc version: 1.0.0"))
	})
	http.ListenAndServe("127.0.0.1:9898", nil)
}
