package connections

import(
	"gov/backend/common/helper"
	"gov/backend/common/config"
	"github.com/pkg/errors"
	"net/http/httputil"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"bytes"
	"log"
)

func ProxyRequest(httpType string, urlPath string, params map[string]string) []byte {
	var buffParams *bytes.Buffer
	var byteParams []byte

	//creating the proxyURL
	proxyStr 	  := config.Set.Proxy.Adress
	proxyURL, err := url.Parse(proxyStr)
	helper.CheckError(errors.Wrap(err, "Can't (url.Parse) passed val from [config.Set.Proxy.Adress]"))

	//creating the URL to be loaded through the proxy
	urlStr    	  := urlPath
	url, err  	  := url.Parse(urlStr)
	helper.CheckError(errors.Wrap(err, "Can't (url.Parse) passed val [urlPath string]"))

	//adding the proxy settings to the Transport object
	transport	  := &http.Transport{Proxy: http.ProxyURL(proxyURL),}

	//adding the Transport object to the http Client
	client 		  := &http.Client{Transport: transport,}

	// set params
	if params != nil {
		byteParams, err = json.Marshal(params)
		helper.CheckError(errors.Wrap(err, "Can't (json.Marshal) passed val [params map[string]string]"))
		buffParams      = bytes.NewBuffer(byteParams)
	} else {
		buffParams      = bytes.NewBuffer(nil)
	}

	//generating the HTTP request
	request, err  := http.NewRequest(httpType, url.String(), buffParams)
	helper.CheckError(errors.Wrap(err, "Can't generate the HTTP request (http.NewRequest)"))

	//adding proxy authentication
	auth          := config.Set.Proxy.Login + ":" + config.Set.Proxy.Pass
	basicAuth     := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	request.Header.Add("Proxy-Authorization", basicAuth)

	//printing the request to the console
	dump, err     := httputil.DumpRequest(request, false)
	helper.CheckError(errors.Wrap(err, "Can't print the request to the console (httputil.DumpRequest)"))
	log.Println(string(dump))

	//calling the URL
	response, err := client.Do(request)
	helper.CheckError(errors.Wrap(err, "Can't call the URL (client.Do)"))

	log.Println(response.StatusCode)
	log.Println(response.Status)
	//getting the response
	data, err 	  := ioutil.ReadAll(response.Body)
	helper.CheckError(errors.Wrap(err, "Can't get the response the URL (ioutil.ReadAll)"))

	return data
}


