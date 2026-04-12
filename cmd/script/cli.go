package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type MjAgentBody struct {
	Base64Array []string `json:"base64Array"`
	NotifyHook  string   `json:"notifyHook"`
	Prompt      string   `json:"prompt"`
	State       string   `json:"state"`
	BotType     string   `json:"botType"`
}

func main() {

	var token = flag.String("token", "", "AI Generate Image Token")

	url := "https://aiyiwei.vip/mj/submit/upload-discord-images"
	method := "POST"

	params := MjAgentBody{
		Base64Array: nil,
		NotifyHook:  "",
		Prompt:      "",
		State:       "",
		BotType:     "MID_JOURNEY",
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	payload := bytes.NewReader(jsonData)

	client := &http.Client{}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(body)
}
