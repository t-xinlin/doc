package workpool

import "net/http"
import (
	_ "net/http/pprof"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	"path"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/pkg/errors"
	"log"
)

const MaxWorkers = 10


func init()  {
	log.Printf("workpool init...")
}


func init_work() {
	JobQueue = make(chan Job, 5)
	dispatcher := NewDispatcher(MaxWorkers)
	dispatcher.Run()
}

func TestWorkPool() {
	ConfigLocalFilesystemLogger("/opt/helloworldserver/log", "access.log", time.Duration(time.Minute*15), time.Duration(time.Minute*5))
	//init_logrus()
	init_work()
	http.HandleFunc("/testCase01", payloadHandler)
	log.Printf("ListenAndServe on :8888\n")
	http.ListenAndServe(":8888", nil)

}

// config logrus log to local filesystem, with file rotation
func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{})

	logrus.AddHook(lfHook)
}

var temlog = logrus.New()

func init_logrus() {
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	logrus.SetOutput(os.Stdout)
	//设置最低loglevel
	logrus.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)

	if err == nil {
		//temlog.Out = file
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}

	//temlog.WithFields(logrus.Fields{
	//	"filename": "123.txt",
	//}).Info("打开文件失败")

	logrus.SetOutput(file)
}
