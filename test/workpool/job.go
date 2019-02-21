package workpool

import (
	"net/http"
	"os"
	"time"
	//"log"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

var (
	MaxWorker       = os.Getenv("MAX_WORKERS")
	MaxQueue        = os.Getenv("MAX_QUEUE")
	MaxLength int64 = 1024
)

type PayloadCollection struct {
	WindowsVersion string    `json:"version"`
	Token          string    `json:"token"`
	Payloads       []Payload `json:"data"`
}

type Payload struct {
	// [redacted]
	UserId string `json:"user_id"`
	Pwd    string `json:"pwd"`
}

func (p *Payload) UploadToS3() error {
	// the storageFolder method ensures that there are no name collision in
	// case we get same timestamp in the key name

	//storage_path := fmt.Sprintf("%v/%v", p.storageFolder, time.Now().UnixNano())
	//
	//bucket := S3Bucket
	//
	//b := new(bytes.Buffer)
	//encodeErr := json.NewEncoder(b).Encode(payload)
	//if encodeErr != nil {
	//	return encodeErr
	//}
	//
	//// Everything we post to the S3 bucket should be marked 'private'
	//var acl = s3.Private
	//var contentType = "application/octet-stream"
	//
	//return bucket.PutReader(storage_path, b, int64(b.Len()), contentType, acl, s3.Options{})

	time.Sleep(time.Second * 1)
	return nil
}

// Job represents the job to be run
type Job struct {
	Payload Payload
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				if err := job.Payload.UploadToS3(); err != nil {
					log.Printf("Error uploading to S3: %s", err.Error())
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

//我们修改了HTTP请求处理函数来创建一个含有载荷（payload）的Job结构，然后将它送到一个叫JobQueue的channel。worker会对它们进行处理。

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	v := r.PostForm.Get("v")
	log.Printf(" 【%s】-------%s", r.URL.Path, v)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var content = &PayloadCollection{}

	// Read the body into a string for json decoding
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ReadAll error: %+v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	defer r.Body.Close()

	err = json.Unmarshal(bytes, &content)
	if err != nil {
		fmt.Printf("Unmarshal error: %+v", err)
		return
	}

	//err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
	//err = json.NewDecoder(r.Body).Decode(&content)
	//if err != nil {
	//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//defer r.Body.Close()
	//fmt.Printf(">>>>>>>>>>content: %+v\n", content)

	// Go through each payload and queue items individually to be posted to S3

	for _, payload := range content.Payloads {

		// let's create a job with the payload
		work := Job{Payload: payload}

		// Push the work onto the queue.
		JobQueue <- work
	}

	w.WriteHeader(http.StatusOK)
}
