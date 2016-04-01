package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var _, _ = url.Parse("")

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

func sendRequest(vendorName string) {
	//todo: load cookie file in run dir now, auto loading in the further
	var cookie_err error
	myVendor.config[vendorName]["head"].(aci_head).data["cookie"], cookie_err = readCookie("cookie")
	if cookie_err != nil {
		fmt.Println(cookie_err)
		return
	}

	//get vendor api url
	targetUrl := myVendor.config[vendorName]["url"].(string)

	tr := &http.Transport{}

	//ignore https verify
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//proxy setup,not neccessary
	if proxy_url, ok := myVendor.config[vendorName]["proxy"]; ok {
		proxy, _ := url.Parse(proxy_url.(string))
		tr.Proxy = http.ProxyURL(proxy)
	}

	client := &http.Client{Transport: tr}

	//get vendor method
	method := myVendor.config[vendorName]["method"].(string)
	var (
		req  *http.Request
		err  error
		body io.Reader //used for POST only
	)

	//setup request body
	switch {
	case method == "POST":
		data_func := myVendor.config[vendorName]["body"].(func() string)
		body = strings.NewReader(data_func())
	case method == "GET":
		body = nil
	default:
		body = nil
	}
	req, err = http.NewRequest(method, targetUrl, body)

	//setup common header
	for key, value := range common_head.data {
		req.Header.Set(key, value)
	}

	//setup vendor header
	//"head" convert interface to aci_head type
	vender_head := myVendor.config[vendorName]["head"].(aci_head)
	for key, value := range vender_head.data {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			if feedback, err := json_decode(body); err == nil {
				fmt.Println(feedback)
			} else {
				fmt.Println(err)
				return
			}
		}
	}

}

func main() {
	for key, _ := range myVendor.config {
		sendRequest(key)
	}
}
