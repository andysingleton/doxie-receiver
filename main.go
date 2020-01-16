package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

type Image struct {
	Name     string `json:"Name"`
	Size     int    `json:"Size"`
	Modified string `json:"Modified"`
}

func GetImage(host string, path string) (string, error) {
	var file image.Image
	var localFile *os.File
	fmt.Println("Downloading image", path)
	commandString := fmt.Sprintf("http://%s/scans/%s", host, path)
	r, err := myClient.Get(commandString)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	file, err = jpeg.Decode(r.Body)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("images/%s", filepath.Base(path))
	fmt.Println("Creating", fileName)
	localFile, err = os.Create(fileName)
	if err != nil {
		fmt.Println("Could not create file", err)
		return "", err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, file, nil)

	bytesWritten, err := localFile.Write(buf.Bytes())
	fmt.Println(bytesWritten, "Bytes written")

	return fileName, nil
}

func GetList(host string, command string) ([]Image, error) {
	var result []Image
	commandString := fmt.Sprintf("http://%s/%s", host, command)
	r, err := myClient.Get(commandString)
	if err != nil {
		return result, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func GetItem(host string, command string) (map[string]interface{}, error) {
	var result map[string]interface{}
	commandString := fmt.Sprintf("http://%s/%s", host, command)
	r, err := myClient.Get(commandString)
	if err != nil {
		return result, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("Body is", err)
		return result, err
	}

	return result, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func poll() (bool, error) {
	result, err := GetItem("192.168.1.131", "hello.json")
	if err != nil {
		return false, err
	}
	fmt.Println("Connected to ", result["MAC"])
	return true, nil
}

func fetchFileList() ([]Image, error) {
	result, err := GetList("192.168.1.131", "scans.json")
	if err != nil {
		return result, err
	}
	return result, nil
}

func deleteFile(host string, scanName string) error {
	var req *http.Request
	var err error

	// todo: Replace all raw prints with Log
	log.Println("Deleting downloaded file from Doxie", scanName)
	commandString := fmt.Sprintf("http://%s/scans/%s", host, scanName)

	// http client doesnt appear to implement DELETE requests directly
	req, err = http.NewRequest("DELETE", commandString, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer req.Body.Close()
	return nil
}

func main() {
	daemonize := flag.Bool("daemon", false, "Run continuously in the foreground")
	hostIp := flag.String("ip", "", "Specify the IP to poll")
	var poll_result bool
	var scans []Image
	var err error
	var fileName string

	for true {
		poll_result, _ = poll()
		if poll_result == true {
			scans, err = fetchFileList()
			if err != nil {
				fmt.Println("Failed to fetch file list", err)
				continue
			}
			for scan := range scans {
				scanName := scans[scan].Name

				// Dont download the webhook log
				if scanName == "/DOXIE/WEBHOOK/LOG.TXT" {
					continue
				}

				fileName, err = GetImage(*hostIp, scanName)
				if err != nil {
					fmt.Println("Failed to download file", err)
					continue
				}

				err = deleteFile(*hostIp, scanName)
				if err != nil {
					fmt.Println("Failed to delete downloaded file", err)
					continue
				}
				go processFile(fileName)
			}
		}
		if *daemonize != true {
			break
		}
		time.Sleep(5 * time.Minute)
	}
}
