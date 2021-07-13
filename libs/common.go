package libs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type (
	Config   map[string]UserInfo
	UserInfo struct {
		JSESSIONID    string
		service_id    string
		token         string
		receive_id    string
		phone         string
		firm_id       string
		reserve_count string
	}
)

func CheckFile(file string) bool {
	isExist := true
	_, err := os.Stat(file)
	if err != nil {
		isExist = false
	}
	return isExist
}

func StructToUrlParams(foo UserInfo) string {
	k := reflect.TypeOf(foo)
	v := reflect.ValueOf(foo)

	var paramList []string
	for i := 0; i < v.NumField(); i++ {
		key := k.Field(i).Name
		value := v.Field(i)
		param := fmt.Sprintf("%s=%s", key, value)
		paramList = append(paramList, param)
	}

	params := strings.Join(paramList, "&")
	return params
}

func HttpGet(url string, data UserInfo, cookie string, ch chan string) {
	params := StructToUrlParams(data)
	//fmt.Println(params)

	//设置headers
	var header = make(map[string]string)
	header["Host"] = "apphd.bjzgh12351.org"
	header["Content-Type"] = "application/x-www-form-urlencoded;charset=utf-8"
	header["Origin"] = "file://"
	header["Cookie"] = cookie
	header["Accept"] = "application/json, text/plain, */*"
	header["User-Agent"] = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 (4385191040)"
	header["Accept-Language"] = "zh-cn"

	url = url + "?" + params

	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	for k, v := range header {
		request.Header.Set(k, v)
	}
	res, _ := client.Do(request)

	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(content))
	ch <- strings.ReplaceAll(string(content), "\n", "")
}

func PrintParamsError(msg string) {
	fmt.Printf("%s; `-h`查看参数\n", msg)
	os.Exit(1)
}

func DynamicCall(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	go f.Call(in)
	return
}
