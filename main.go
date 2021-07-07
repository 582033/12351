package main

import (
	"12351/libs"
	"flag"
	"fmt"
	"os"
	"sync"
)

func takePart(iniFile string, callBack func()) {
	if !libs.CheckFile(iniFile) {
		fmt.Println("缺少配置文件,`-help`查看参数")
		os.Exit(1)
	}

	iniSection := "user1"
	conf := libs.LoadConf(iniFile, iniSection)
	cookie := "JSESSIONID=" + conf.JSESSIONID

	url := "http://apphd.bjzgh12351.org/mobile/api/service/reserve.jsp"
	res := libs.HttpGet(url, conf, cookie)
	fmt.Println(res)

	callBack()
}

var t = flag.Int("t", 1, "并发线程数")
var i = flag.String("i", "./conf.ini", "配置文件ini位置")

func main() {
	flag.Parse()

	if *t < 1 {
		fmt.Println("线程数错误\n")
		return
	}
	if *i == "" {
		fmt.Println("配置文件路径错误\n")
		return
	}
	thread := *t
	iniFile := *i

	var wg sync.WaitGroup
	for i := 0; i < thread; i++ {
		wg.Add(1)
		go takePart(iniFile, func() {
			defer wg.Done()
		})
	}
	wg.Wait()
}
