// Package gce is a Google Compute Environment utility package
package gce

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const timeout time.Duration = 1 * time.Second

// Get is a utility wrapper for handling http get calls
// and reading their response bodies to a string
func Get(url string) (string, error) {
	// Create a new http.Client with a timeout
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Metadata-Flavor", "Google")

	// Make the GET request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Invalid response: %v %v", resp.StatusCode, resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to parse response: %v", resp.Body)
	}
	if len(bodyBytes) < 1 {
		return "", fmt.Errorf("No response received")
	}
	return bytes.NewBuffer(bodyBytes).String(), nil
}

// Cluster returns the name of the GKE cluster
func Cluster() (string, error) {
	return Get("http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-name")
}

// Instance returns the instance id of the current instance
func Instance() (string, error) {
	h, err := Get("http://metadata.google.internal/computeMetadata/v1/instance/hostname")
	if err != nil {
		return "", err
	}

	// The instance/id path returns the numeric ID; we need the alphabetic ID, which
	// is the host portion of the full hostname.
	pieces := strings.Split(h, ".")
	return pieces[0], nil
}

// Project returns the project id of the current instance
func Project() (string, error) {
	return Get("http://metadata.google.internal/computeMetadata/v1/project/project-id")
}

// Zone returns the zone of the current instance
func Zone() (string, error) {
	z, err := Get("http://metadata.google.internal/computeMetadata/v1/instance/zone")
	if err != nil {
		return "", err
	}

	// The zone metadata call prefixes the project information; we are only interested in
	// the last segment: the zone itself.
	pieces := strings.Split(z, "/")
	return pieces[len(pieces)-1], nil
}

// PrivateIPv4 returns the internal IPv4 address of the node
func PrivateIPv4() (string, error) {
	return Get("http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip")
}

// PublicIPv4 returns the public IPv4 address of the node
func PublicIPv4() (string, error) {
	return Get("http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip")
}
