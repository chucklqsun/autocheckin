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
	return true
}

func checkin(vendorName string, account string) bool {
	//todo: load cookie file in run dir now, auto loading in the further
	var cookie_err error
	myVendor.config[vendorName]["head"].(aci_head).data["cookie"], cookie_err = readCookie(vendorName + "." + account + ".cookie")
	if cookie_err != nil {
		fmt.Println(cookie_err)
		return false
	}
	checkin_req := aci_request{
		url:    myVendor.config[vendorName]["url_checkin"].(string),
		proxy:  myVendor.config[vendorName]["proxy"].(string),
		method: myVendor.config[vendorName]["method"].(string),
		head:   myVendor.config[vendorName]["head"].(aci_head).data,
		body:   myVendor.config[vendorName]["body"],
		result: myVendor.config[vendorName]["result"],
	}
	return checkin_req.sendRequest()
}

func controller(vendorName string, account string) {
	if !checkin(vendorName, account) {
		login(vendorName, account)
		checkin(vendorName, account)
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
