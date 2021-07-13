package main

import (
	"12351/libs"
	"flag"
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
)

const (
	snapUpUrl = "http://apphd.bjzgh12351.org/mobile/api/service/reserve.jsp"
	signInUrl = "http://apphd.bjzgh12351.org/mobile/api/check/take-service.jsp"
)

func signIn(iniFile string, callBack func()) {
	if libs.CheckFile(iniFile) == false {
		libs.PrintParamsError("缺少配置文件")
	}
	configs := libs.LoadConf(iniFile)
	for section, conf := range configs {
		cookie := "JSESSIONID=" + conf.JSESSIONID
		ch := make(chan string)
		go libs.HttpGet(signInUrl, conf, cookie, ch)
		fmt.Printf("[Section.%s执行结果]:%s\n", section, <-ch)
	}

	callBack()
}

func snapUp(iniFile string, callBack func()) {
	if libs.CheckFile(iniFile) == false {
		libs.PrintParamsError("缺少配置文件")
	}

	configs := libs.LoadConf(iniFile)
	for section, conf := range configs {
		cookie := "JSESSIONID=" + conf.JSESSIONID
		ch := make(chan string)
		go libs.HttpGet(snapUpUrl, conf, cookie, ch)
		fmt.Printf("[Section.%s执行结果]:%s\n", section, <-ch)
	}
	callBack()
}

var t = flag.Int("t", 1, "并发线程数")
var i = flag.String("i", "./conf.ini", "配置文件ini位置")
var s = flag.String("s", "snapUp", "- snapUp, 抢购\n- signIn, 抢购成功后的扫码确认")
var p = flag.Int("p", 100, "线程池大小")

func main() {
	flag.Parse()

	thread := *t
	iniFile := *i
	funcName := *s
	poolSize := *p

	funcs := map[string]interface{}{
		"signIn": signIn,
		"snapUp": snapUp,
	}

	if thread < 1 {
		libs.PrintParamsError("线程数错误")
	}
	if iniFile == "" {
		libs.PrintParamsError("配置文件路径错误")
	}
	if poolSize < 1 {
		libs.PrintParamsError("线程池大小错误")
	}

	if funcs[funcName] == nil {
		libs.PrintParamsError(fmt.Sprintf("不支持的参数`%s`", funcName))
	}

	//fmt.Println(funcName)
	var wg sync.WaitGroup

	p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
		libs.DynamicCall(funcs, funcName, iniFile, func() {
			defer wg.Done()
		})
	})
	defer p.Release()

	for i := 0; i < thread; i++ {
		wg.Add(1)
		_ = p.Invoke(i)
	}
	wg.Wait()
}
