package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

func json_decode(body []byte) (data map[string]interface{}, err error) {
	err = nil
	buf := bytes.NewBuffer(body)
	reader, err := gzip.NewReader(buf)
	if err == nil { //gzip result
		dec := json.NewDecoder(reader)
		err = dec.Decode(&data)
	} else { //plain text
		err = json.Unmarshal(body, &data)
	}
	return
}
