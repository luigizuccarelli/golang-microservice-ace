package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// MiddlewareLogin a http response and request wrapper for database insert
// It takes a both response and request objects and returns void
func MiddlewareLogin(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Could not read body data (MiddlewareLogin) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	}

	apitoken, err := LoginData(body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Login error (MiddlewareLogin) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{MetaInfo: apitoken}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW data successfully retrieved", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareData a http response and request wrapper for database update
// It takes a both response and request objects and returns void
func MiddlewareData(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Could not read body data (MiddlewareData) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	}

	data, err := AllData(body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Subscriptions data read (MiddlewareSubscriptions) " + err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		e := json.Unmarshal(data, &payload)
		if e != nil {
			response = Response{StatusCode: "500", Status: "ERROR", Message: "Subscriptions unmarshal error (MiddlewareSubscriptions) " + e.Error()}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			response = Response{StatusCode: "200", Status: "OK", Message: "MW data read successfully", Payload: payload}
		}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	logger.Trace(fmt.Sprintf("used to mask cc %v", r))
	logger.Trace(fmt.Sprintf("config data  %v", config))
	fmt.Fprintf(w, "ok version 1.0")
}
