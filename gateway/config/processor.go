package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	wd, _ := os.Getwd()
	main, _ := ioutil.ReadFile(fmt.Sprintf("%s/config/config.json", wd))
	config := make(map[string]interface{})
	json.Unmarshal(main, &config)

	var endpoints []map[string]interface{}
	files := []string{}
	filepath.Walk(fmt.Sprintf("%s/config/endpoints", wd), func(path string, f os.FileInfo, err error) error {
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	for _, f := range files {
		var endpoint []map[string]interface{}
		j, _ := ioutil.ReadFile(f)
		json.Unmarshal(j, &endpoint)
		endpoints = append(endpoints, endpoint...)
	}

	config["endpoints"] = endpoints
	krakend, _ := json.Marshal(config)

	whitelists, _ := ioutil.ReadFile(fmt.Sprintf("%s/config/whitelist.json", wd))
	whitelist := []string{}

	json.Unmarshal(whitelists, &whitelist)
	wl := strings.Join(whitelist, "|")
	wl = strings.ReplaceAll(string(krakend), "-!MSuryaIksanudin!-", wl)

	file := fmt.Sprintf("%s/krakend.json", wd)
	ioutil.WriteFile(file, []byte(wl), 0644)
}
