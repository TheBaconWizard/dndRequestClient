package dndRequestClient

import (
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"time"

	"bytes"
	"compress/gzip"
	"strings"

	"github.com/Carcraftz/cclient"
	"github.com/andybalholm/brotli"

	http "github.com/Carcraftz/fhttp"

	tls "github.com/Carcraftz/utls"
)

// Get makes a simple get request to the specified url.
// url is the website url you want to access formatted ass http://...
// headers are the headers you want to include.
// returns a http response object as well as a decoded string response
func Get(url string, headers map[string]string) (responseRequest *http.Response, respBody string) {
	return HandleReq("GET", url, "", headers)
}

// Post makes a simple post request to the specified url.
// url is the website url you want to access formatted ass http://...
// headers are the headers you want to include.
// body is the body you want to include if none us ""
// returns a http response object as well as a decoded string response
func Post(url string, headers map[string]string, body string) (responseRequest *http.Response, respBody string) {
	return HandleReq("POST", url, body, headers)
}

// Patch makes a simple patch request to the specified url.
// url is the website url you want to access formatted ass http://...
// headers are the headers you want to include.
// body is the body you want to include if none us ""
// returns a http response object as well as a decoded string response
func patch(url string, headers map[string]string, body string) (responseRequest *http.Response, respBody string) {
	return HandleReq("PATCH", url, body, headers)
}

// HandleReq can be used to handle and types of requests
// method is the method of the request ex. "GET","POST","PATCH"
// myUrl is the website url
// input is the body of the request if none put ""
// headers are the headers you want to include
// returns a http response object as well as a decoded string response
func HandleReq(method string, myUrl string, input string, headers map[string]string) (responseRequest *http.Response, respBody string) {

	client, err := cclient.NewClient(tls.HelloChrome_Auto, "", true, time.Duration(6))
	if err != nil {
		log.Fatal(err)
	}

	var req *http.Request
	if input == "" {
		req, err = http.NewRequest(method, myUrl, nil)
	} else {
		req, err = http.NewRequest(method, myUrl, strings.NewReader(input))
	}

	if err != nil {
		panic(err)
	}

	// convert header map[string]string to map[string][]string

	scrappyHeader := http.Header{}

	for k, v := range headers {
		scrappyHeader.Add(k, v)

	}

	//master header order, all your headers will be ordered based on this list and anything extra will be appended to the end
	//if your site has any custom headers, see the header order chrome uses and then add those headers to this list
	masterheaderorder := []string{
		"host",
		"connection",
		"cache-control",
		"device-memory",
		"viewport-width",
		"rtt",
		"downlink",
		"ect",
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-full-version",
		"sec-ch-ua-arch",
		"sec-ch-ua-platform",
		"sec-ch-ua-platform-version",
		"sec-ch-ua-model",
		"upgrade-insecure-requests",
		"user-agent",
		"accept",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-user",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"cookie",
	}
	headermap := make(map[string]string)
	//TODO: REDUCE TIME COMPLEXITY (This code is very bad)
	headerorderkey := []string{}
	for _, key := range masterheaderorder {
		for k, v := range scrappyHeader {
			lowercasekey := strings.ToLower(k)
			if key == lowercasekey {
				headermap[k] = v[0]
				headerorderkey = append(headerorderkey, lowercasekey)
			}
		}

	}
	for k, v := range req.Header {
		if _, ok := headermap[k]; !ok {
			headermap[k] = v[0]
			headerorderkey = append(headerorderkey, strings.ToLower(k))
		}
	}

	//ordering the pseudo headers and our normal headers
	req.Header = http.Header{
		http.HeaderOrderKey:  headerorderkey,
		http.PHeaderOrderKey: {":method", ":authority", ":scheme", ":path"},
	}
	//set our Host header
	u, err := url.Parse(myUrl)
	if err != nil {
		panic(err)
	}
	for k := range scrappyHeader {
		if k != "Content-Length" && !strings.Contains(k, "Poptls") {
			v := scrappyHeader.Get(k)
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Host", u.Host)
	resp, err := client.Do(req)
	//forward decoded response body
	encoding := resp.Header["Content-Encoding"]
	responseBody, err := ioutil.ReadAll(resp.Body)
	finalres := ""
	finalres = string(responseBody)
	if len(encoding) > 0 {
		if encoding[0] == "gzip" {
			unz, err := gUnzipData(responseBody)
			if err != nil {
				panic(err)
			}
			finalres = string(unz)
		} else if encoding[0] == "deflate" {
			unz, err := enflateData(responseBody)
			if err != nil {
				panic(err)
			}
			finalres = string(unz)
		} else if encoding[0] == "br" {
			unz, err := unBrotliData(responseBody)
			if err != nil {
				panic(err)
			}
			finalres = string(unz)
		} else {
			fmt.Println("UNKNOWN ENCODING: " + encoding[0])
			finalres = string(responseBody)
		}
	} else {
		finalres = string(responseBody)
	}

	return resp, finalres
}

func gUnzipData(data []byte) (resData []byte, err error) {
	gz, _ := gzip.NewReader(bytes.NewReader(data))
	defer gz.Close()
	respBody, err := ioutil.ReadAll(gz)
	return respBody, err
}
func enflateData(data []byte) (resData []byte, err error) {
	zr, _ := zlib.NewReader(bytes.NewReader(data))
	defer zr.Close()
	enflated, err := ioutil.ReadAll(zr)
	return enflated, err
}
func unBrotliData(data []byte) (resData []byte, err error) {
	br := brotli.NewReader(bytes.NewReader(data))
	respBody, err := ioutil.ReadAll(br)
	return respBody, err
}
