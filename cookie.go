package main

import (
	"bufio"
	"github.com/chucklqsun/glog"
	"io"
	"os"
	"strings"
)

type aci_cookie map[string]string

func splitCookie(cookieStr string) aci_cookie {
	var ret aci_cookie = make(map[string]string)
	cookie := strings.Split(cookieStr, ";")
	for _, v := range cookie {
		v = strings.Trim(v, " ")
		key := strings.Split(v, "=")
		if len(key) == 2 {
			ret[key[0]] = key[1]
		} else {
			ret[key[0]] = ""
		}
	}
	return ret
}

func readCookie(infile string) (str string, err error) {

	file, err := os.Open(infile)
	if err != nil {
		glog.Errorln("Failed to open the input file ", infile)
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
			glog.Errorln("A too long line, seems unexpected.")
			return
		}
		str = string(line)
	}
	return
}
