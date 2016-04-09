package main

import (
	"flag"
	"github.com/chucklqsun/glog"
	"sync"
)

var w sync.WaitGroup

func login(vendorName string, account string) bool {
	aa, err := readAccount("account")
	if err != nil {
		glog.Errorln("Account err:", err)
		return false
	}
	username := aa[vendorName+"."+account]["username"]
	password := aa[vendorName+"."+account]["password"]

	login_req := aci_request{
		account: account,
		url:     myVendor.config[vendorName]["url_login"].(string),
		proxy:   myVendor.config[vendorName]["proxy"].(string),
		method:  myVendor.config[vendorName]["method"].(string),
		head:    myVendor.config[vendorName]["head"].(aci_head).data,
		body:    myVendor.config[vendorName]["body_login"].(func(string, string) string)(username, password),
		cookie:  "",
		result:  myVendor.config[vendorName]["result_login"],
	}
	return login_req.sendRequest()
	return true
}

func checkin(vendorName string, account string) bool {
	//todo: load cookie file in run dir now, auto loading in the further
	var (
		cookie_err error
		cookie     string
	)

	cookie, cookie_err = readCookie(vendorName + "." + account + ".cookie")
	if cookie_err != nil {
		glog.Errorln("Cookie err:", cookie_err)
		return false
	}

	checkin_req := aci_request{
		account: account,
		url:     myVendor.config[vendorName]["url_checkin"].(string),
		proxy:   myVendor.config[vendorName]["proxy"].(string),
		method:  myVendor.config[vendorName]["method"].(string),
		head:    myVendor.config[vendorName]["head"].(aci_head).data,
		body:    myVendor.config[vendorName]["body_checkin"].(func(string) string)(cookie),
		cookie:  cookie,
		result:  myVendor.config[vendorName]["result_checkin"],
	}
	return checkin_req.sendRequest()
}

func controller(vendorName string, account string) {
	mode := myVendor.config[vendorName]["mode"].(int)
	switch mode {
	case 1:
		login(vendorName, account)
	case 2:
		checkin(vendorName, account)
	case 3:
		if !checkin(vendorName, account) {
			login(vendorName, account) //current not support duokan dueto capt
			checkin(vendorName, account)
		}
	}
	w.Done()
}

func main() {
	//initial command for log
	flag.Parse()
	flag.Lookup("log_dir").Value.Set("./")
	flag.Lookup("log_name").Value.Set("log")
	flag.Lookup("alsologtostderr").Value.Set("true")

	glog.Info("Starting ...")

	for key, _ := range myVendor.config {
		if myVendor.config[key]["status"].(int) == 0 {
			continue
		}
		for _, account := range myVendor.config[key]["account"].([]string) {
			w.Add(1)
			go controller(key, account)
			glog.Infoln("Job -> ", key, ":", account)
		}
	}
	w.Wait()
	glog.Info("End")

	//make sure log msg flush into file
	glog.Flush()
}
