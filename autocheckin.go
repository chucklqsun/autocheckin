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

func checkin(vendorName string) {
	//todo: load cookie file in run dir now, auto loading in the further
	var cookie_err error
	myVendor.config[vendorName]["head"].(aci_head).data["cookie"], cookie_err = readCookie("cookie")
	if cookie_err != nil {
		fmt.Println(cookie_err)
		return
	}
	checkin_req := aci_request{
		url:    myVendor.config[vendorName]["url_checkin"].(string),
		proxy:  myVendor.config[vendorName]["proxy"].(string),
		method: myVendor.config[vendorName]["method"].(string),
		head:   myVendor.config[vendorName]["head"].(aci_head).data,
		body:   myVendor.config[vendorName]["body"],
	}
	checkin_req.sendRequest()

	w.Done()
}

func main() {
	fmt.Println("Begin")
	for key, _ := range myVendor.config {
		if myVendor.config[key]["status"].(int) == 0 {
			continue
		}

		w.Add(1)
		go checkin(key)
		fmt.Println("Job -> ", key)
	}
	w.Wait()
	fmt.Println("End")
}
