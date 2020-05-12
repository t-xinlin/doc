package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	UA = "Golang Downloader from Kejibo.com"
)

func main() {
	f, err := os.OpenFile("./file.exe", os.O_RDWR, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
	if err != nil {
		panic(err)
	}
	stat, err := f.Stat() //获取文件状态
	if err != nil {
		panic(err)
	}
	f.Seek(stat.Size(), 0) //把文件指针指到文件末，当然你说为何不直接用 O_APPEND 模式打开，没错是可以。我这里只是试验。
	url1 := "http://dl.google.com/chrome/install/696.57/chrome_installer.exe"
	var req http.Request
	req.Method = "GET"
	req.Close = true
	req.URL, err = url.Parse(url1)
	if err != nil {
		panic(err)
	}
	header := http.Header{}
	header.Set("Range", "bytes="+strconv.Itoa(int(stat.Size()))+"-")
	req.Header = header
	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		panic(err)
	}
	written, err := io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}
	println("written: ", written)
}

var dir string
var port int
var staticHandler http.Handler

// 初始化参数
func init() {
	dir = path.Dir(os.Args[0])
	flag.IntVar(&port, "port", 800, "服务器端口")
	flag.Parse()
	fmt.Println("dir:", http.Dir(dir))
	staticHandler = http.FileServer(http.Dir(dir))
}

func main11() {
	http.HandleFunc("/", StaticServer)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 静态文件处理
func StaticServer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("path:" + req.URL.Path)
	if req.URL.Path != "/down/" {
		staticHandler.ServeHTTP(w, req)
		return
	}

	io.WriteString(w, "hello, world!\n")
}

//上传（client）

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}
	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))

	return nil
}

// sample usage
func main12() {
	target_url := "http://localhost:9090/upload"
	filename := "./astaxie.pdf"
	postFile(filename, target_url)
}

// 初始化参数
func init() {
	dir = path.Dir(os.Args[0])
	flag.IntVar(&port, "port", 800, "服务器端口")
	flag.Parse()
	fmt.Println("dir:", http.Dir(dir))
}

func main13() {
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
