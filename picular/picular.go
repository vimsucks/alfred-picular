package picular

import (
	"encoding/json"
	"io"
	"net/http"
)

func SearchColor(query string) (*Response, error) {
	url := "https://backend.picular.co/api/search?query=" + query
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type Response struct {
	Colors    []Color
	Primary   string
	Secondary string
}

type Color struct {
	Color string
	Img   string
	Light bool
}
