package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var _, _ = url.Parse("")

type aci_request struct {
	url    string
	proxy  string
	method string
	head   map[string]string
	body   interface{}
	result interface{}
}

func (ar *aci_request) sendRequest() bool {

	tr := &http.Transport{}

	//ignore https verify
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//proxy setup,left empty if no use
	if !strings.EqualFold(ar.proxy, "") {
		proxy, _ := url.Parse(ar.proxy)
		tr.Proxy = http.ProxyURL(proxy)
	}

	client := &http.Client{Transport: tr}

	//get vendor method
	method := ar.method
	var (
		req  *http.Request
		err  error
		body io.Reader //used for POST only
	)

	//setup request body
	switch {
	case method == "POST":
		data_func := ar.body.(func() string)
		body = strings.NewReader(data_func())
	case method == "GET":
		body = nil
	default:
		body = nil
	}
	req, err = http.NewRequest(method, ar.url, body)

	//setup common header
	for key, value := range common_head.data {
		req.Header.Set(key, value)
	}

	//setup vendor header
	//"head" convert interface to aci_head type
	for key, value := range ar.head {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return false
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return false
		} else {
			if feedback, err := json_decode(body); err == nil {
				fmt.Println(feedback)
				if ret := ar.result.(func(map[string]interface{}) bool)(feedback); !ret {
					fmt.Println("exec ", ar.url, " fail")
					return false
				} else {
					fmt.Println("exec ", ar.url, " success")
					return true
				}
			} else {
				fmt.Println(err)
				return false
			}
		}
	}

	return false

}
