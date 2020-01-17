package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net"
	"net/http"
)

type hostAPIInterface interface {
	getIp() string
	Get(string) ([]byte, error)
	GetFileList() ([]Image, error)
	GetStatus() (map[string]interface{}, error)
	GetJpeg(string) (image.Image, error)
	Delete(string) error
}

type hostAPI struct {
	Ip   net.IP
	Port int
}

func (host hostAPI) getIp() string {
	return host.Ip.String()
}

func (host hostAPI) GetStatus() (map[string]interface{}, error) {
	var status map[string]interface{}
	result, err := host.Get("hello.json")
	if err != nil {
		return status, err
	}
	err = json.Unmarshal(result, &status)
	if err != nil {
		fmt.Println("Body is", err)
		return status, err
	}
	return status, nil
}

func (host hostAPI) Get(path string) ([]byte, error) {
	commandString := fmt.Sprintf("http://%s:%d/%s", host.Ip, host.Port, path)
	r, err := myClient.Get(commandString)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	return body, err
}

func (host hostAPI) GetJpeg(path string) (image.Image, error) {
	var file image.Image
	commandString := fmt.Sprintf("http://%s:%d/scans/%s", host.Ip, host.Port, path)
	r, err := myClient.Get(commandString)
	if err != nil {
		return file, err
	}
	defer r.Body.Close()

	file, err = jpeg.Decode(r.Body)
	return file, err
}

func (host hostAPI) GetFileList() ([]Image, error) {
	var scans []Image
	result, err := host.Get("scans.json")
	if err != nil {
		return scans, err
	}
	err = json.Unmarshal(result, &scans)
	return scans, err
}

func (host hostAPI) Delete(scanName string) error {
	var req *http.Request

	// http client doesnt appear to implement DELETE requests directly
	commandString := fmt.Sprintf("http://%s:%d/scans/%s", host.Ip, host.Port, scanName)
	fmt.Println("1")
	req, err := http.NewRequest("DELETE", commandString, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
