package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	apiURL         = "http://labs.bible.org/api/?passage=votd&type=json"
	cacheDir       = ".verse"
	defaultTimeout = 5 * time.Second
)

var (
	fileName = time.Now().Local().Format("20060102")
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func checkError(err error) {
	if err != nil {
		// Fail without causing an error
		fmt.Print("🚧 Could not load verse 🚧")
		os.Exit(0)
	}
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		checkError(err)
	}

	exPath := filepath.Dir(ex)
	cache := exPath + string(os.PathSeparator) + cacheDir
	exist, err := exists(cache)
	if err != nil {
		checkError(err)
	}

	if !exist {
		err = os.Mkdir(cache, os.ModePerm)
		if err != nil {
			checkError(err)
		}
	}

	cacheFile := cache + string(os.PathSeparator) + fileName
	exist, err = exists(cacheFile)
	if err != nil {
		checkError(err)
	}

	var votd votd

	if exist {
		data, err := ioutil.ReadFile(cacheFile)
		if err != nil {
			checkError(err)
		}

		err = json.Unmarshal(data, &votd)
		if err != nil {
			checkError(err)
		}
	} else {
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
			checkError(err)
		}

		if response.StatusCode < 200 || response.StatusCode >= 400 {
			panic(response.Status)
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			checkError(err)
		}

		err = json.Unmarshal(data, &votd)
		if err != nil {
			checkError(err)
		}

		_ = ioutil.WriteFile(cacheFile, data, os.ModePerm)
	}

	fmt.Print(votd.toString())
}
