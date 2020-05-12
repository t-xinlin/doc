package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fmt.Println("run")
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
