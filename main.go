package main

import (
	"12351/libs"
	"flag"
	"fmt"
	"os"
	"sync"
)

func signIn(iniFile string, callBack func()) {
	if libs.CheckFile(iniFile) == false {
		fmt.Println("缺少配置文件,`-help`查看参数")
		os.Exit(1)
	}
	configs := libs.LoadConf(iniFile)
	//url := "http://apphd.bjzgh12351.org/mobile/api/check/take-service.jsp"
	for conf := range configs {
		/*
			cookie := "JSESSIONID=" + conf.JSESSIONID
			res := libs.HttpGet(url, conf, cookie)
			fmt.Println(res)
		*/
		fmt.Println(conf)
	}

	callBack()
}

func snapUp(iniFile string, callBack func()) {
	if libs.CheckFile(iniFile) == false {
		fmt.Println("缺少配置文件,`-help`查看参数")
		os.Exit(1)
	}

	configs := libs.LoadConf(iniFile)
	url := "http://apphd.bjzgh12351.org/mobile/api/check/take-service.jsp"
	for section, conf := range configs {
		cookie := "JSESSIONID=" + conf.JSESSIONID
		ch := make(chan string)
		go libs.HttpGet(url, conf, cookie, ch)
		fmt.Printf("[Section.%s执行结果]:%s\n", section, <-ch)
	}
	callBack()
}

var t = flag.Int("t", 1, "并发线程数")
var i = flag.String("i", "./conf.ini", "配置文件ini位置")
var s = flag.String("s", "snapUp", "snapUp, 抢购\n - signIn, 抢购成功后的扫码确认")

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
		go snapUp(iniFile, func() {
			defer wg.Done()
		})
	}
	wg.Wait()
}
