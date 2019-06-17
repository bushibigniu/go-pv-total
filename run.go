package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type resource struct {
	url string
	target string
	start int
	end int
}

func main()  {

	//flag 获取命令行参数 go run run.go --total=18 --filepath=/tmp/a.txt
	total := flag.Int("total",10,"how many num")
	filePath := flag.String("filepath","/tmp/test.log","file path")
	flag.Parse()

	//返会的是指针，所以加*
	log.Println(*total, *filePath)


	//构造真实url
	//url := []string
	res := ruleResource()
	list := buildUrl(res)

	//生成total 行日志
	var logStr string
	for i:=0;i< *total;i++  {

		currentUrl := list[ rand.Intn(len(list)-1) ]
		refer := list[rand.Intn(len(list)-1)]
		ua := list[rand.Intn(len(list)-1)]

		//logStr := makeLog(currentUrl, refer, ua)
		logStr = logStr + makeLog(currentUrl, refer, ua)+ "\n"

		//写入文件 ioutil.WriteFile 会覆盖写
		//ioutil.WriteFile(*filePath, []byte(logStr), 0644)
	}

	fd, err := os.OpenFile(*filePath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil{
		panic(err)
	}
	_, err = fd.Write([]byte(logStr))
	if err != nil{
		panic(err)
	}
	defer fd.Close()
	fmt.Println(logStr)
}



func makeLog(current, refer, ua string) string {

	u:=  url.Values{}
	u.Set("refer",refer)
	u.Set("time","1")
	u.Set("ua",ua)
	u.Set("url",current)

	// Encode encodes the values into ``URL encoded'' form
	// ("bar=baz&foo=quux") sorted by key.
	paramStr := u.Encode()
	fmt.Println( paramStr )

	logTemp := "127.0.0.1 - - [16/Jun/2019:00:37:59 +0800] \"" +
		"GET " +
		"{paramStr}" +
		"{ua}"

	logStr := strings.Replace(logTemp,"{paramStr}", paramStr, -1)
	logStr = strings.Replace(logStr,"{ua}", ua, -1)

	return logStr

}

func buildUrl(res []resource) []string {

	var urlList []string

	for _, item := range res{
		if len(item.target) == 0 {
			urlList = append(urlList, item.url)
		} else {
			for i:=item.start; i< item.end; i++ {
				urlStr := strings.Replace(item.url,item.target,strconv.Itoa(i), -1)
				urlList = append(urlList, urlStr)
			}
		}

	}

	return urlList
}

func ruleResource() []resource {
	var res []resource
	indexRes := resource{ //index
		url:"http://localhost:9901",
		target:"",
		start:0,
		end:0,
	}
	listRes := resource{ //list {$id} 是一个字符串占位符，不是变量
		url:"http://localhost:9901/list/{$id}.html",
		target:"{$id}",
		start:0,
		end:21,
	}
	detailRes := resource{
		url:"http://localhost:9901/movie/{$id}.html",
		target: "{$id}",
		start:0,
		end:1200,
	}
	res = append(res, indexRes)
	res = append(res, listRes)
	res = append(res, detailRes)

	return res
}
