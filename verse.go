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
	defaultTimeout = 5 * time.Second
)

type verse struct {
	Bookname string `json:"bookname"`
	Chapter  string `json:"chapter"`
	Verse    string `json:"verse"`
	Text     string `json:"text"`
}

func (v verse) getHeader() string {
	return v.Bookname + " " + v.Chapter + ":" + v.Verse
}

func (v verse) toString(includeHeader bool) string {
	str := ""

	if includeHeader {
		str += v.getHeader() + "\n"
	}

	str += "\t(" + v.Verse + ") " + v.Text

	return str
}

type votd []verse

func (v votd) toString() string {
	numVerses := len(v)

	if numVerses > 1 {
		header := v[0].getHeader() + "-" + v[numVerses-1].Verse + "\n"
		text := ""

		for t := range v {
			text += v[t].toString(false) + "\n"
		}

		return header + text
	}

	return v[0].toString(true)
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

	var votd votd
	err = json.Unmarshal(data, &votd)
	if err != nil {
		panic(err)
	}

	fmt.Print(votd.toString())
}
