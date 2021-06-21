package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type request struct {
	client *http.Client
	cache  *cache
}

func NewRequest(client *http.Client, cache *cache) request {
	return request{
		client: client,
		cache:  cache,
	}
}

func (r request) do(request *http.Request, config config, service string, cacheKey string) (*data, error) {
	var response *data
	token := request.Header.Get(config.header.Identity)
	if token == "" {
		return nil, errors.New("authentication token is not provided")
	}

	res, found := r.cache.get(cacheKey)
	if found {
		response = res
	} else {
		s, ok := config.services[service]
		res := make(chan *data)
		if service != "" && ok {
			go r.exe(config, s.ValidateUrl, token, res)
			response = <-res
			if response != nil {
				go r.cache.set(cacheKey, response, config.cacheTTL)
			}
		} else {
			ok := false
			for _, v := range config.services {
				go r.exe(config, v.ValidateUrl, token, res)
			}

			for range config.services {
				response = <-res
				if response != nil {
					go r.cache.set(cacheKey, response, config.cacheTTL)
					ok = true
					break
				}
			}

			if !ok {
				return nil, errors.New("authentication token not provided or not valid")
			}
		}
	}

	return response, nil
}

func (r request) copy(request *http.Request, config config, data *data) *http.Request {
	request.Header.Del(config.header.Identity)
	for k, v := range config.responseMap {
		header, ok := data.Payload[k]
		if ok {
			request.Header.Set(v, header)
		}
	}

	return request
}

func (r request) exe(config config, endpoint string, token string, response chan *data) {
	rq, err := http.NewRequest(config.method, endpoint, nil)
	if err != nil {
		fmt.Println("Error Response ", err)
		response <- nil
	}

	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set(config.header.Identity, token)

	rs, err := r.client.Do(rq)
	if err != nil {
		fmt.Println("Error Response ", err)
		response <- nil
	}

	defer rs.Body.Close()

	rsBodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		fmt.Println("Error Response ", err)
		response <- nil
	}

	data := &data{}
	err = json.Unmarshal(rsBodyBytes, data)
	if err != nil {
		fmt.Println("Error Decode ", string(rsBodyBytes))
		response <- nil
	}

	if rs.StatusCode != http.StatusOK {
		fmt.Println("Error Upstream ", data)
		response <- nil
	}

	response <- data
}
