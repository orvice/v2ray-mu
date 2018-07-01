package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"github.com/catpie/musdk-go"
)

func (u *UserManager) httpReq(uri string, method string, buffer string) (string, int, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, uri, strings.NewReader(buffer))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Token", cfg.WebApi.Token)
	req.Header.Set("ServiceType", strconv.Itoa(musdk.TypeV2ray))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return string(body), res.StatusCode, err
}

func (u *UserManager) httpGet(uri string) (string, int, error) {
	return u.httpReq(uri, http.MethodGet, "")
}

func (u *UserManager) httpPost(uri, data string) (string, int, error) {
	return u.httpReq(uri, http.MethodPost, data)
}