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
		"url": "https://www.duokan.com/checkin/v0/status",
		//"proxy":  "http://192.168.159.1:8888",
		"method": "POST",
		"head": aci_head{
			data: map[string]string{
				"Origin":           "https://www.duokan.com",
				"X-Requested-With": "com.duokan.reader",
				"Referer":          "https://www.duokan.com/hs/user/task",
				"cookie":           "",
			},
		},
		"body": func() string {
			csrf := func(e int, t byte) int {
				return (131*e + int(t)) % 65536
			}
			timestamp := time.Now().Unix()
			cookie := strings.Split(myVendor.config["duokan"]["head"].(aci_head).data["cookie"], ";")
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
			return body
		},
	}

	//other vendor...
}
