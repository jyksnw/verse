package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	apiURL         = "http://labs.bible.org/api/?passage=votd&type=json"
	defaultTimeout = 30 * time.Second
)

type VOTD struct {
	Bookname string `json:"bookname"`
	Chapter  string `json:"chapter"`
	Verse    string `json:"verse"`
	Text     string `json:"text"`
}

func (v VOTD) GetHeader() string {
	return v.Bookname + " " + v.Chapter + ":" + v.Verse
}

func (v VOTD) ToString(includeHeader bool) string {
	str := ""

	if includeHeader {
		str += v.GetHeader() + "\n"
	}

	str += "\t(" + v.Verse + ") " + v.Text

	return str
}

type Verses []VOTD

func (v Verses) ToString() string {
	numVerses := len(v)

	if numVerses > 1 {
		header := v[0].GetHeader() + "-" + v[numVerses-1].Verse + "\n"
		text := ""

		for t := range v {
			text += v[t].ToString(false) + "\n"
		}

		return header + text
	}

	return v[0].ToString(true)
}

func main() {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: defaultTimeout,
		}).Dial,
		TLSHandshakeTimeout: defaultTimeout,
	}

	var netClient = &http.Client{
		Timeout:   defaultTimeout,
		Transport: netTransport,
	}

	response, err := netClient.Get(apiURL)
	if err != nil {
		panic(err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 400 {
		panic(response.Status)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var votd Verses
	json.Unmarshal(data, &votd)
	fmt.Print(votd.ToString())
}
