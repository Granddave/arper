package arp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetVendorName(macAddress string) string {
	url := fmt.Sprintf("https://api.macvendors.com/%s", macAddress)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to execute request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	// Don't care if the API doesn't find any vendor
	if resp.StatusCode == http.StatusNotFound {
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %v", resp.StatusCode)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body %v", err)
		return ""
	}

	return string(body)
}
