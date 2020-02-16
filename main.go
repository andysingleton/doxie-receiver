package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

type Image struct {
	Name     string `json:"Name"`
	Size     int    `json:"Size"`
	Modified string `json:"Modified"`
}

func getIpAddr(hostIp string) net.IP {
	ipAddr := net.ParseIP(hostIp)
	if ipAddr == nil {
		log.Fatal("You must provide a valid address")
	} else {
		log.Println("Using address", ipAddr.String())
	}
	return ipAddr
}

func poll(host hostAPIInterface) bool {
	log.Println("Polling")
	result, err := host.GetStatus()
	if err != nil {
		return false
	}
	log.Println("Connected to", result["MAC"])
	return true
}

func getImage(path string, host hostAPIInterface, imagePath string) (string, error) {
	var localFile *os.File
	log.Println("Downloading image", path)
	image, err := host.GetJpeg(path)
	if err != nil {
		return "", err
	}

	// todo: Check and create the images directory. make "images" a config variable
	t := time.Now()
	fileDate := fmt.Sprintf(t.Format("02-01-2006_15:04:05"))

	// get the number from the filename
	re := regexp.MustCompile(`IMG_([0-9]+).JPG`)
	fileId := re.FindStringSubmatch(path)
	if len(fileId) == 0 {
		log.Println("Could not parse an Id from the filename: ", path)
		return "", errors.New("Could not parse an Id from the filename")
	}

	filePath := fmt.Sprintf("%s/%s_%s.jpeg", imagePath, fileId[1], fileDate)

	check, _ := os.Stat(filePath)
	if check != nil {
		log.Println("File already exists:", imagePath)
		return "", errors.New("Refusing to overwrite pre-existing target file")
	}

	localFile, err = os.Create(filePath)
	if err != nil {
		log.Println("Could not create file", err)
		return "", err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, image, nil)
	bytesWritten, err := localFile.Write(buf.Bytes())
	log.Println(bytesWritten, "Bytes written")

	return filePath, nil
}

func deleteFile(scanName string, host hostAPIInterface) error {
	log.Println("Deleting downloaded file from Doxie", scanName)
	err := host.Delete(scanName)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var scans []Image
	var err error
	var fileName string

	daemonize := flag.Bool("daemon", false, "Run continuously in the foreground")
	hostIp := flag.String("ip", "", "Specify the Doxie IP")
	hostPort := flag.Int("port", 80, "Specify the Doxie port")
	imagePath := flag.String("output", "images", "Output path for stored images")
	noDelete := flag.Bool("no-delete", false, "Don't remove images from Doxie")
	flag.Parse()

	if *hostIp == "" {
		fmt.Println("You must provide an IP address for your Doxie")
		os.Exit(1)
	}

	host := hostAPI{
		Ip:   getIpAddr(*hostIp),
		Port: *hostPort,
	}

	for true {
		if poll(host) {
			scans, err = host.GetFileList()
			if err != nil {
				log.Println("Failed to fetch file list", err)
				continue
			}
			for scan := range scans {
				scanName := scans[scan].Name

				// Dont download the webhook log
				if scanName == "/DOXIE/WEBHOOK/LOG.TXT" {
					continue
				}

				fileName, err = getImage(scanName, host, *imagePath)
				if err != nil {
					log.Println("Failed to download file", err)
					continue
				}

				if *noDelete == false {
					err = deleteFile(scanName, host)
					if err != nil {
						log.Println("Failed to delete downloaded file", err)
					}
				}

				go processFile(fileName)
			}
			log.Println("Done")
		}
		if *daemonize != true {
			break
		}
		time.Sleep(30 * time.Second)
	}
}

// todo: Replace all raw prints with Log
// todo: Handle partial json returns
