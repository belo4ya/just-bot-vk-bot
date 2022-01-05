package fuapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	host             string = "https://ruz.fa.ru/api"
	searchEndpoint   string = host + "/search"
	scheduleEndpoint string = host + "/schedule"
)

type (
	Group struct {
		Id          string `json:"id"`
		Label       string `json:"label"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}
)

func GetGroup(name string) {
	p := fmt.Sprintf("?term=%s&type=group", name)
	url := searchEndpoint + p

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var group []Group
	if err := json.Unmarshal(body, &group); err != nil {
		log.Fatalln(err)
	}

	log.Println(url)
	log.Println(group)
	log.Println(string(body))
}
