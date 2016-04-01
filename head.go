package main

type aci_head struct {
	data map[string]string
}

func (h *aci_head) Set(key string, value string) {
	h.data[key] = value
}

var common_head aci_head

func init() {
	//initialize
	common_head.data = make(map[string]string)
	//set common key-value
	common_head.data["User-Agent"] = "Mozilla/5.0 (Linux; Android 4.4; Nexus 5 Build/_BuildID_) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36 XiaoMi/MiuiBrowser/2.1.1"
	common_head.data["Connection"] = "keep-alive"
	common_head.data["x-wap-profile"] = "http://wap1.huawei.com/uaprof/HW_HUAWEI_MT7-CL00_2_20140903.xml"
	common_head.data["Accept"] = "*/*"
	common_head.data["Accept-Language"] = "zh-CN,en-US;q=0.8"
	common_head.data["Accept-Encoding"] = "gzip,deflate"
	common_head.data["Content-type"] = "application/x-www-form-urlencoded"
}
