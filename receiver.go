package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
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

func main() {
	/*
		# For each one, process the file
		# Download, Delete, Process (thread)
		GET /scans/DOXIE/JPEG/IMG_XXXX.JPG
		# 404 Not Found

		DELETE /scans/DOXIE/JPEG/IMG_XXXX.JPG
		# 404 Not Found
	*/

	daemonize := flag.Bool("daemon", false, "Run continuously in the foreground")
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
				if scanName == "/DOXIE/WEBHOOK/LOG.TXT" {
					continue
				}
				fileName, err = GetImage("192.168.1.131", scanName)
				if err != nil {
					fmt.Println("Failed to download file", err)
					continue
				}
				fmt.Println("filename is", fileName)

				// generate searcheable pdf
				// name by date

				//err = processFile(file)
				//if err != nil {
				//	fmt.Println("Could not process file", err)
				//	continue
				//}
				//err = deleteFile(scan)
				//if err != nil {
				//	fmt.Println("Failed to delete downloaded file", err)
				//	continue
				//}
			}
		}
		if *daemonize != true {
			break
		}
		//time.Sleep(1 * time.Minute)
		time.Sleep(5 * time.Second)
	}
}
