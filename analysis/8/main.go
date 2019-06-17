package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"github.com/mgutz/str"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

const HANDLE_DIG  = "/dig?"

type cmdParams struct {
	filePath string
	goroutineNum int
}

type digData struct {
	time string
	url string
	referer string
	ua string
}
type urlData struct {
	data digData
	uid string
}

type storageBlock struct {
	counterType string
	storageModel string
	urlNode string
}
var log = logrus.New()
func init()  {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}

func main()  {

	//1.获取参数， flag
	goroutineNum := flag.Int("total",10, "consumer number for goroutine")
	//消费日志存储路径
	filePath := flag.String("filepath", "/tmp/dig.log", "the project runtime log")
	//打日志存储路径，就是打日志存储到哪里
	l := flag.String("", "/tmp", "the project runtime log")
	flag.Parse()

	params := cmdParams{
		*filePath,
		*goroutineNum,
	}

	//2.打日志
	file, err := os.OpenFile(*l, os.O_CREATE|os.O_RDONLY,0644)
	if err != nil{
		panic(err)
	}
	defer file.Close()

	log.Out = file
	log.Infoln("exec param start")
	log.Infoln("params filePath=%s; goroutineNum=%s; l=%l",filePath)

	//3.初始化channel,用于数据传递
	var logChannel = make(chan string, 3*params.goroutineNum)
	var pvChannel = make(chan urlData, params.goroutineNum)
	var uvChannel = make(chan urlData, params.goroutineNum)
	//var storageChannel = make(chan storageBlock, params.goroutineNum)


	//4.日志消费
	err = readFileLine(params, logChannel)
	if err != nil {
		panic(err)
	}

	//创建一组日志处理
	for i:=0;i< params.goroutineNum ;i++  {
		go logConsumer(logChannel, pvChannel, uvChannel )
	}

	//创建pv uv 统计器

	//创建存储器
}

func logConsumer(logchan chan string, pvchan, uvchan chan urlData)  {
	for logstr := range logchan{
		//切割日志，抠出打点上报日志
		data := cutLogFetchData(logstr)

		//uid 模拟生成uid
		hasher := md5.New()
		hasher.Write([]byte(data.url+data.ua))
		uid := hex.EncodeToString(hasher.Sum(nil))


		//这边完成解析工作
		//...TODO

		uData := urlData{
			data:data,
			uid: uid,
		}

		log.Infoln(uData)
		pvchan <- uData
		uvchan <- uData
	}

}

//日志消费
func readFileLine(params cmdParams, logChan chan string) error {
	fd, err := os.Open(params.filePath)
	if err != nil {
		log.Warningln("err")
		return err
	}
	defer fd.Close()

	count := 0
	reader := bufio.NewReader(fd)
	for  {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Infoln("")
				time.Sleep(time.Millisecond)
			} else {
				log.Warning("")
			}
		}
		logChan <- line
		count++

		if count % (1000*params.goroutineNum) == 0  {
			log.Infoln("")
		}
	}

	return nil
}

func cutLogFetchData(logstr string) digData {
	//两边去空格
	logstr = strings.TrimSpace(logstr)

	//position := str.IndexOf(logstr, HANDLE_DIG, 0)
	position := getIndexOf(logstr, HANDLE_DIG, 0)
	if position == -1 {
		return digData{}
	}

	position += len(HANDLE_DIG)
	positionHttp := getIndexOf(logstr, "HTTP/", position)

	str := str.Substr(logstr, position, positionHttp-position)

	url, err := url.Parse("http://localhost/?"+str)
	if err != nil {
		return  digData{}
	}
	data := url.Query()


	return digData{
		data.Get("time"),
		data.Get("url"),
		data.Get("refer"),
		data.Get("ua"),
	}
}

func getIndexOf(s string, needle string, start int) int {
	l := len(s)
	if needle == "" {
		if start < 0 {
			return 0
		} else if start < l {
			return start
		} else {
			return l
		}
	}
	if start < 0 || start > l-1 {
		return -1
	}
	pos := strings.Index(s[start:], needle)
	if pos == -1 {
		return -1
	}
	return start + pos
}
