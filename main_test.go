package main

import (
	"net"
	"os"
	"regexp"
	"testing"
)

func TestGetIpAddr(t *testing.T) {
	result := getIpAddr("10.0.0.1")
	resultIp := net.ParseIP("10.0.0.1")
	if resultIp.String() != result.String() {
		t.Errorf("Function did not return valid IP")
	}
}

func TestPollTrue(t *testing.T) {
	host := mockAPI{}
	result := poll(host)
	if result != true {
		t.Errorf("Poll did not return true")
	}
}

func TestPollFalse(t *testing.T) {
	host := mockAPIFail{}
	result := poll(host)
	if result != false {
		t.Errorf("Poll did not return false")
	}
}

func TestGetImage(t *testing.T) {
	host := mockAPI{}
	testFilePath := "test/0_01-02-2002 1:2:3.jpeg"
	_ = os.Remove(testFilePath)
	result, err := getImage("IMG_0.JPG", host, "test")
	if err != nil {
		t.Errorf("Function returned error %s", err)
	}

	re := regexp.MustCompile(`test/0_[0-9-]+_[0-9:]+.jpeg`)
	if !re.MatchString(result) {
		t.Errorf("Function did not return correct file path %s", result)
	}

	_, err = os.Stat(result)
	if err != nil {
		t.Errorf("Could not stat test file")
	}
	_ = os.Remove(result)
}

func TestDeleteFile(t *testing.T) {
	host := mockAPI{}
	err := deleteFile("foobar", host)
	if err != nil {
		t.Errorf("Function errored: %s", err)
	}
}
