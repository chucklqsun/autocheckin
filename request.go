package main

import (
	"crypto/tls"
	"github.com/chucklqsun/glog"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var _, _ = url.Parse("")

type aci_request struct {
	account string
	url     string
	proxy   string
	method  string
	head    map[string]string
	body    string
	cookie  string
	result  interface{}
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
		body = strings.NewReader(ar.body)
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
	req.Header.Set("cookie", ar.cookie)

	//setup vendor header
	//"head" convert interface to aci_head type
	for key, value := range ar.head {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		glog.Errorln("Resp err:", err)
		return false
	}

	defer resp.Body.Close()

	respBody, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		glog.Errorln("ReadAll err:", err)
		return false
	} else {
		if feedback, err := json_decode(respBody); err == nil {
			glog.Infoln(feedback)
			if ret := ar.result.(func(map[string]interface{}) bool)(feedback); !ret {
				glog.Infoln("exec ", ar.account, ar.url, " fail")
				return false
			} else {
				glog.Infoln("exec ", ar.account, ar.url, " success")
				return true
			}
		} else {
			glog.Errorln("Json decode error", err)
			glog.Errorln("body:", respBody)
			return false
		}
	}

	return false

}
