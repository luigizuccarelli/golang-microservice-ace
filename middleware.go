package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"io/ioutil"
	"net/http"
	"time"
)

func LoginData(body []byte) (string, error) {

	logger.Trace("In function LoginData")

	var j map[string]interface{}
	var schema SchemaInterface
	var apitoken = ""

	e := json.Unmarshal(body, &j)
	if e != nil {
		logger.Error(fmt.Sprintf("LoginData %v\n", e.Error()))
		return apitoken, e
	}
	logger.Debug(fmt.Sprintf("json data %v\n", j))

	// use this for cors
	//res.setHeader("Access-Control-Allow-Origin", "*")
	//res.setHeader("Access-Control-Allow-Methods", "POST")
	//res.setHeader("Access-Control-Allow-Headers", "accept, content-type")

	// lets first check in the in-memory cache
	key := j["username"].(string) + ":" + j["password"].(string)
	val, err := redisdb.Get("hash").Result()

	if val == "" || err != nil {
		logger.Info(fmt.Sprintf("Key not found in cache %s", key))
		// make the middleware call
		// set up http object
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest("GET", config.Url+"/username/"+j["username"].(string)+"/password/"+j["password"].(string), nil)
		req.Header.Set("token", config.Token)
		resp, err := client.Do(req)
		if err != nil {
			logger.Error(fmt.Sprintf("%s ", err.Error()))
			return apitoken, err
		}
		logger.Info(fmt.Sprintf("Connected to host %s", config.Url))
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(fmt.Sprintf("%s ", err.Error()))
			return apitoken, err
		}
		logger.Trace(fmt.Sprintf("Response from MW %s", string(body)))

		errs := json.Unmarshal(body, &schema)
		if err != nil {
			logger.Error(fmt.Sprintf("%s ", errs.Error()))
			return apitoken, errs
		}
		logger.Debug(fmt.Sprintf("Schema Data %v ", schema))

		all, _ := json.MarshalIndent(schema, "", "	")
		err = redisdb.Set("all", all, time.Hour).Err()
		err = redisdb.Set("hash", key, time.Hour).Err()

		if err != nil {
			return apitoken, err
		}

		apitoken = xid.New().String()
		err = redisdb.Set("apitoken", apitoken, time.Hour).Err()
	} else if val != key {
		logger.Error(fmt.Sprintf("Hash token's don't match %s != %s", val, key))
		return apitoken, errors.New("hash token does not match")
	} else if val == key {
		logger.Info(fmt.Sprintf("Key found in cache %s", key))
		apitoken = xid.New().String()
		err = redisdb.Set("apitoken", apitoken, time.Hour).Err()
	}
	return apitoken, nil
}

func AllData(b []byte) ([]byte, error) {

	logger.Trace("In function AllData")

	var subs []byte
	var data string
	var j map[string]interface{}

	e := json.Unmarshal(b, &j)
	if e != nil {
		logger.Error(fmt.Sprintf("AllData %v\n", e.Error()))
		return subs, e
	}
	logger.Debug(fmt.Sprintf("json data %v\n", j))

	// lets first check in the in-memory cache
	apitoken := j["apitoken"].(string)
	val, err := redisdb.Get("apitoken").Result()
	if apitoken != val || err != nil {
		return subs, err
	} else {
		data, e = redisdb.Get("all").Result()
		if e != nil {
			return subs, e
		}
	}
	subs = []byte(data)
	return subs, nil
}
