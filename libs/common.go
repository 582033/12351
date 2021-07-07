package libs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"gopkg.in/ini.v1"
)

type UserInfo struct {
	JSESSIONID    string
	service_id    string
	token         string
	receive_id    string
	phone         string
	firm_id       string
	reserve_count string
}

func CheckFile(file string) bool {
	isExist := true
	_, err := os.Stat(file)
	if err != nil {
		isExist = true
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

func LoadConf(file string, section string) UserInfo {
	cfg, err := ini.Load(file)
	if err != nil {
		panic(err)
	}

	conf := UserInfo{
		JSESSIONID:    cfg.Section(section).Key("JSESSIONID").String(),
		service_id:    cfg.Section(section).Key("service_id").String(),
		token:         cfg.Section(section).Key("token").String(),
		receive_id:    cfg.Section(section).Key("receive_id").String(),
		phone:         cfg.Section(section).Key("phone").String(),
		firm_id:       cfg.Section(section).Key("firm_id").String(),
		reserve_count: "1",
	}
	return conf
}

func HttpGet(url string, data UserInfo, cookie string) string {
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
	return strings.ReplaceAll(string(content), "\n", "")
}
