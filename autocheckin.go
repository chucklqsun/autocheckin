package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

var w sync.WaitGroup

func readCookie(infile string) (str string, err error) {

	file, err := os.Open(infile)
	if err != nil {
		fmt.Println("Failed to open the input file ", infile)
		return
	}

	defer file.Close()

	br := bufio.NewReader(file)
	for {
		line, isPrefix, err1 := br.ReadLine()
		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}
		if isPrefix {
			fmt.Println("A too long line, seems unexpected.")
			return
		}
		str = string(line)
	}
	return
}

func login(vendorName string, account string) bool {
	aa, err := readAccount("account")
	if err != nil {
		fmt.Println("Account err:", err)
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
		fmt.Println("Cookie err:", cookie_err)
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
		//		if !checkin(vendorName, account) {
		//			//login(vendorName, account)    //current not support duokan dueto capt
		//			checkin(vendorName, account)
		//		}
	}
	w.Done()
}

func main() {
	fmt.Println("Begin")
	for key, _ := range myVendor.config {
		if myVendor.config[key]["status"].(int) == 0 {
			continue
		}
		for _, account := range myVendor.config[key]["account"].([]string) {
			w.Add(1)
			go controller(key, account)
			fmt.Println("Job -> ", key, ":", account)
		}
	}
	w.Wait()
	fmt.Println("End")
}
