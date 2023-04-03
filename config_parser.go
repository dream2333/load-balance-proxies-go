package main

import (
	"encoding/json"
	"log"
	"os"
)

type Proxy struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Weight int    `json:"weight"`
}

type Proxies []Proxy

func LoadJson(filename string) Proxies {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("文件打开出错: ", err)
	}

	// Now let's unmarshall the data into `payload`
	var proxies Proxies
	err = json.Unmarshal(content, &proxies)
	if err != nil {
		log.Fatal("JSON反序列化失败: ", err)
	}

	// Let's print the unmarshalled data!
	return proxies
}

func WeightedRandomProxy(proxies Proxies) Proxy {
	var weights []float32
	for _, v := range proxies {
		weights = append(weights, float32(v.Weight))
	}
	index := WeightedRandomIndex(weights)
	return proxies[index]
}
