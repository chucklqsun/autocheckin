package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type aci_account map[string]map[string]string

func readAccount(infile string) (aa aci_account, err error) {
	aa = make(map[string]map[string]string)
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
		ln := strings.Split(string(line), "=")
		aa[ln[0]] = make(map[string]string)
		up := strings.Split(string(ln[1]), "|")
		aa[ln[0]]["username"] = up[0]
		aa[ln[0]]["password"] = up[1]
	}
	return
}
