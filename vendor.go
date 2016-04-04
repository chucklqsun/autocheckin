package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type vendor_config map[string]interface{}

type vendor struct {
	config map[string]vendor_config
}

var myVendor vendor

func init() {
	myVendor.config = make(map[string]vendor_config)

	//duokan app
	myVendor.config["duokan"] = map[string]interface{}{
		"mode":        2, //checkin only
		"account":     []string{"a1", "a2"},
		"status":      1,
		"url_login":   "xxxxxx",
		"url_checkin": "https://www.duokan.com/checkin/v0/checkin",
		//"url_checkin": "https://www.duokan.com/checkin/v0/status",
		"proxy": "",
		//"proxy":  "http://192.168.159.1:8888",
		"method": "POST",
		"head": aci_head{
			data: map[string]string{
				"Origin":           "https://www.duokan.com",
				"X-Requested-With": "com.duokan.reader",
				"Referer":          "https://www.duokan.com/hs/user/task",
			},
		},
		"body_login": func() string {
			return ""
		},
		"body_checkin": func(cookieStr string) string {
			csrf := func(e int, t byte) int {
				return (131*e + int(t)) % 65536
			}
			timestamp := time.Now().Unix()
			cookie := strings.Split(cookieStr, ";")
			var device_id string
			for _, v := range cookie {
				v = strings.Trim(v, " ")
				if key := strings.Split(v, "="); key[0] == "device_id" {
					device_id = key[1]
				}
			}

			tk := []byte(fmt.Sprintf("%s&%s", device_id, strconv.FormatInt(timestamp, 10)))
			var e int = 0
			for _, v := range tk {
				e = csrf(e, v)
			}
			body := fmt.Sprintf("_t=%s&_c=%d", strconv.FormatInt(timestamp, 10), e)
			//fmt.Println(body)
			return body
		},
		"result_login": func(feedback map[string]interface{}) bool {
			if value, ok := feedback["result"]; ok {
				retCode := value.(float64)
				if retCode == 0.0 || retCode == 500002.0 {
					return true
				} else {
					return false
				}
			}
			return false
		},
		"result_checkin": func(feedback map[string]interface{}) bool {
			if value, ok := feedback["result"]; ok {
				retCode := value.(float64)
				if retCode == 0.0 || retCode == 500002.0 {
					return true
				} else {
					return false
				}
			}
			return false
		},
	}

	//zimuzu.tv
	myVendor.config["zimuzu"] = map[string]interface{}{
		"mode":        1, //login only
		"account":     []string{"a1"},
		"status":      1,
		"url_login":   "http://www.zimuzu.tv/User/Login/ajaxLogin",
		"url_checkin": "",
		"proxy":       "",
		"method":      "POST",
		"head": aci_head{
			data: map[string]string{
				"Origin":           "http://www.zimuzu.tv",
				"X-Requested-With": "XMLHttpRequest",
				"Referer":          "http://www.zimuzu.tv/user/login",
			},
		},
		"body_login": func(account string, password string) string {
			body := fmt.Sprintf("account=%s&password=%s&remember=1&url_back=http://www.zimuzu.tv/", account, password)
			return body
		},
		"body_checkin": func(cookieStr string) string {
			return ""
		},
		"result_login": func(feedback map[string]interface{}) bool {
			if value, ok := feedback["status"]; ok {
				retCode := value.(float64)
				if retCode == 1.0 {
					return true
				} else {
					return false
				}
			}
			return false
		},
		"result_checkin": func(feedback map[string]interface{}) bool {
			if value, ok := feedback["status"]; ok {
				retCode := value.(float64)
				if retCode == 1.0 {
					return true
				} else {
					return false
				}
			}
			return false
		},
	}

	//other vendor...
}
