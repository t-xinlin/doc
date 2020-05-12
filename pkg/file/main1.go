package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main1() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":1789", nil)
}

func upload1(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Fprintln(w, "upload ok!")
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(tpl))
}

const tpl = `<html>
<head>
<title>上传文件</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
 <input type="file" name="uploadfile" />
 <input type="hidden" name="token" value="{...{.}...}"/>
 <input type="submit" value="upload" />
</form>
</body>
</html>`

//https://github.com/songtianyi/rrframework

func main2() {
	f, err := os.OpenFile("K:/file.mp3", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	stat, err := f.Stat() //获取文件状态
	if err != nil {
		panic(err)
	} //把文件指针指到文件末，当然你说为何不直接用 O_APPEND 模式打开，没错是可以。我这里只是试验。
	url := "http://127.0.0.1:3000/assets/37-02.mp3"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", "bytes="+strconv.FormatInt(stat.Size(), 10)+"-")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	written, err := io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}
	println("written: ", written)
}

//服务器的话就更简单了，这个是忽略url中的/assets/，直接找到对应的raido目录

var staticHandler1 http.Handler

// 静态文件处理
func StaticServer1(w http.ResponseWriter, req *http.Request) {
	fmt.Println("path:" + req.URL.Path)
	staticHandler.ServeHTTP(w, req)
}
func init() {
	staticHandler = http.StripPrefix("/assets/", http.FileServer(http.Dir("radio")))
}

func main3() { //已经有静态文件了
	http.HandleFunc("/assets/", StaticServer)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//if r.Method == "POST" {
//
//        r.ParseMultipartForm(8 << 20)
//        title := r.ParseFormValue["title"]
//        fhs := r.MultipartForm.File["radio[]"]
//        options := r.MultipartForm.Value["options[]"]
//        answers := r.MultipartForm.Value["answers[]"]
//
//        l := len(options)
//        optionDirs := make([]string, l)
//        t := time.Now()
//        for i := 0; i < l; i++ {
//            file, err := fhs[i].Open()
//            if err != nil {
//                panic(err)
//            }
//            filename := fhs[i].Filename
//            f, err := os.OpenFile("statics/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
//            if err != nil {
//                panic(err)
//            }
//            defer f.Close()
//            io.Copy(f, file)
//            optionDirs = append(optionDirs, filename)
//        }
//        db.InsertHomework(&db.HomeWork{
//            Title:      title,
//            Options:    options,
//            OptionDirs: optionDirs,
//            Answers:    answers,
//            Time:       t,
//        })
//        sess := session.GlobalSessionManager.SessionStart(w, r)
//        if sess != nil {
//            sess.Set("flash", true)
//        }
//        defer sess.SessionRelease()
//        http.Redirect(w, r, "/homeworks", http.StatusFound)
//    }
