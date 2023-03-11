package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)


func makeGetData(parameters map[string]interface{}) string {
	if parameters == nil {
		return ""
	}
	stringList := make([]string, len(parameters))
	for k := range parameters {
		stringList = append(stringList, fmt.Sprintf("%s=%s", k, parameters[k]))
	}
	return strings.Join(stringList, "&")
}

func makePostData(parameters map[string]interface{}) *bytes.Reader {
	if parameters == nil {
		return nil
	}
	jsonData, err := json.Marshal(parameters)
	if err != nil {
		fmt.Println("makePostData func")
		log.Fatal(jsonData)
	}
	return bytes.NewReader(jsonData)
}


func Request(url string, method string, parameters map[string]interface{}) map[string]interface{} {
	
	var data *bytes.Reader
	
	switch method {

	case "GET":
		getParameters := makeGetData(parameters)
		if getParameters != "" {
			url = url + "?" + getParameters
		}
		fmt.Println("trying to request")
		request, err := http.NewRequest(method, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Println("sending request")
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		responseByte, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		ResponseData := make(map[string]interface{})
		json.Unmarshal(responseByte, &ResponseData)

		fmt.Println(res.StatusCode)
		return ResponseData
	
	case "POST":
		postParameters := makePostData(parameters)
		data = postParameters
		request, err := http.NewRequest(method, url, data)
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Println("sending request")
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Println("sending request err: ", err.Error())
		}
		responseByte, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		ResponseData := make(map[string]interface{})
		json.Unmarshal(responseByte, &ResponseData)

		fmt.Println(res.StatusCode)
		return ResponseData
	}

	return make(map[string]interface{})
}


func TestRequest() {

	Request(
		"http://154.91.170.231:4000/get-all-orders?broker=binance&email=mohamad@gmail.com",
		"GET",
		nil)

	fmt.Println("request function done!")
}