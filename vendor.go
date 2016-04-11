package main

import (
	"fmt"
	"github.com/chucklqsun/ptlogin"
	"strconv"
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
		"body_login": func(account string, password string, idx string) string {
			return ""
		},
		"body_checkin": func(cookieStr string) string {
			csrf := func(e int, t byte) int {
				return (131*e + int(t)) % 65536
			}
			timestamp := time.Now().Unix()
			var device_id string
			cookieHash := splitCookie(cookieStr)
			if value, ok := cookieHash["device_id"]; ok {
				device_id = value
			} else {
				device_id = ""
			}

			tk := []byte(fmt.Sprintf("%s&%s", device_id, strconv.FormatInt(timestamp, 10)))
			var e int = 0
			for _, v := range tk {
				e = csrf(e, v)
			}
			body := fmt.Sprintf("_t=%s&_c=%d", strconv.FormatInt(timestamp, 10), e)
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
		"body_login": func(account string, password string, idx string) string {
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

	//daoju.qq.com get JD
	myVendor.config["daoju_jd"] = map[string]interface{}{
		"mode":        3, //checkin only
		"account":     []string{"a1"},
		"status":      1,
		"url_login":   "xxxxxx",
		"url_checkin": "http://apps.game.qq.com/ams/ame/ame.php?ameVersion=0.3&sServiceType=dj&iActivityId=11117&sServiceDepartment=djc&set_info=djc",
		"proxy":       "",
		//"proxy":       "http://192.168.159.1:8888",
		"method": "POST",
		"head": aci_head{
			data: map[string]string{
				"Origin":           "http://daoju.qq.com/index.shtml",
				"X-Requested-With": "XMLHttpRequest",
				"Referer":          "http://daoju.qq.com/index.shtml",
			},
		},
		"body_login": func(account string, password string, idx string) string {
			var pt ptlogin.Ptlogin
			//username,password
			pt.SetInput(account, password)
			pt.SetCookieName("daoju_jd." + idx + ".cookie")
			pt.Ptui_checkVC()
			body := ""
			return body
		},
		"body_checkin": func(cookieStr string) string {
			getACSRFToken := func(str string) int {
				if str != "" {
					var hash int = 5381
					strBit := []byte(str)
					for _, v := range strBit {
						hash += (hash << 5) + int(v)
					}
					return hash & 0x7fffffff
				}
				return 0
			}
			var g_tk string
			cookieHash := splitCookie(cookieStr)
			if skey, ok := cookieHash["skey"]; ok {
				g_tk = fmt.Sprintf("%d", getACSRFToken(skey))
			} else {
				g_tk = ""
			}

			params := map[string]string{
				"iFlowId":            "95581",
				"sServiceType":       "dj",
				"sServiceDepartment": "djc",
				"ch":                 "10000",
				"sDeviceID":          "00000000-60e9-6e8d-222b-9e7c5de6ffde",
				"appVersion":         "39",
				"g_tk":               g_tk,
				"appSource":          "android",
				"iActivityId":        "11117",
			}
			var body string
			for key, value := range params {
				body += key + "=" + value + "&"
			}
			return body
		},
		"result_login": func(feedback map[string]interface{}) bool {
			if value, ok := feedback["ret"]; ok {
				retCode := value.(string)
				if retCode == "0" {
					return true
				} else {
					return false
				}
			}
			return false
		},
		"result_checkin": func(feedback map[string]interface{}) bool {
			if value, ok := feedback["ret"]; ok {
				retCode := value.(string)
				if retCode == "0" || retCode == "600" {
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
