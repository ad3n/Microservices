package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const pluginName = "authenticator"

var HandlerRegisterer registrable = registrable(pluginName)

type registrable string

func (r registrable) RegisterHandlers(f func(
	name string,
	handler func(
		context.Context,
		map[string]interface{},
		http.Handler,
	) (http.Handler, error),
)) {
	f(pluginName, func(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
		config := parse(extra)
		cache := NewCache()
		request := NewRequest(&http.Client{Timeout: time.Duration(config.timeout) * time.Second}, cache)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sha := sha1.New()
			service := r.Header.Get(config.header.Service)
			_, err := sha.Write([]byte(fmt.Sprintf("%s:%s", service, r.Header.Get(config.header.Identity))))
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			response := do(handler, r, w, request, cache, config, service, config.whitelist, fmt.Sprintf("%x", sha.Sum(nil)))
			if response == nil {
				return
			}

			handler.ServeHTTP(w, request.copy(r, config, response))
		}), nil
	})
}

func do(
	handler http.Handler,
	r *http.Request,
	w http.ResponseWriter,
	executor request,
	cache *cache,
	config config,
	service string,
	whitelist string,
	cacheKey string,
) *data {
	match, _ := regexp.MatchString(whitelist, r.URL.Path)
	if match {
		h := strings.ToLower(r.Header.Get(config.header.Service))
		s, ok := config.services[h]
		if h == "" || !ok {
			for _, v := range config.services {
				if v.LogoutPath == r.URL.Path || s.LoginPath == r.URL.Path {
					go cache.invalidate(cacheKey)

					break
				}
			}
		} else if s.LogoutPath == r.URL.Path || s.LoginPath == r.URL.Path {
			go cache.invalidate(cacheKey)
		}

		handler.ServeHTTP(w, r)

		return nil
	}

	response, err := executor.do(r, config, service, cacheKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)

		return nil
	}

	return response
}

func parse(configs map[string]interface{}) config {
	config := config{}

	services, err := json.Marshal(configs["services"])
	if err == nil {
		json.Unmarshal(services, &config.services)
	}

	method, ok := configs["method"].(string)
	if !ok {
		method = "POST"
	}
	config.method = method

	header, err := json.Marshal(configs["header"])
	if err == nil {
		json.Unmarshal(header, &config.header)
	}

	timeout, ok := configs["timeout"].(int)
	if !ok {
		timeout = 3
	}
	config.timeout = timeout

	cacheTTL, ok := configs["cache_ttl"].(int)
	if !ok {
		cacheTTL = 490
	}
	config.cacheTTL = cacheTTL

	whitelist, ok := configs["whitelist"].(string)
	if !ok {
		whitelist = ""
	}
	config.whitelist = whitelist

	responseMap, err := json.Marshal(configs["response_to_header"])
	if err == nil {
		json.Unmarshal(responseMap, &config.responseMap)
	}

	return config
}

func main() {}
