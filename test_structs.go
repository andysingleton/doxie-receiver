package main

import (
	"errors"
	"image"
	"image/jpeg"
	"net"
	"os"
	"testing"
)

type mockAPI struct {
	Ip   net.IP
	Port int
}

func (host mockAPI) getIp() string {
	return net.ParseIP("10.0.0.1").String()
}

func (host mockAPI) Get(path string) ([]byte, error) {
	return []byte("wef"), nil
}

func (host mockAPI) GetJpeg(path string) (image.Image, error) {
	var file image.Image
	var t testing.T
	testImage, err := os.Open("test-image.jpg")
	if err != nil {
		t.Errorf("Couldn't open test image file")
	}
	file, err = jpeg.Decode(testImage)
	if err != nil {
		t.Errorf("Couldn't decode test image file")
	}
	return file, nil
}

func (host mockAPI) GetFileList() ([]Image, error) {
	var images []Image
	images = append(images, Image{
		Name:     "foo",
		Size:     1,
		Modified: "bar",
	})
	return images, nil
}

func (host mockAPI) GetStatus() (map[string]interface{}, error) {
	var result map[string]interface{}
	return result, nil
}

func (host mockAPI) Delete(scanName string) error {
	return nil
}

type mockAPIFail struct {
	Ip   net.IP
	Port int
}

func (host mockAPIFail) getIp() string {
	return net.ParseIP("10.0.0.1").String()
}

func (host mockAPIFail) Get(path string) ([]byte, error) {
	return []byte("wef"), errors.New("failure")
}

func (host mockAPIFail) GetJpeg(path string) (image.Image, error) {
	var file image.Image
	return file, errors.New("failure")
}

func (host mockAPIFail) GetFileList() ([]Image, error) {
	var images []Image
	images = append(images, Image{
		Name:     "foo",
		Size:     1,
		Modified: "bar",
	})
	return images, errors.New("failure")
}

func (host mockAPIFail) GetStatus() (map[string]interface{}, error) {
	var result map[string]interface{}
	return result, errors.New("failure")
}

func (host mockAPIFail) Delete(scanName string) error {
	return errors.New("failure")
}
