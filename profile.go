package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/batmac/nextdns_exporter/nextdns"
)

type Profile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Fingerprint string `json:"fingerprint"`
}

type Profiles struct {
	Profiles []Profile `json:"data"`
}

func getProfiles(c *nextdns.Client) ([]Profile, error) {
	req, err := c.Get(c.BaseURL + "profiles")
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	if req.StatusCode != 200 {
		log.Printf("Error: %v", req.Status)
		return nil, errors.New(req.Status)
	}
	defer req.Body.Close()

	var content bytes.Buffer
	_, err = content.ReadFrom(req.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	log.Printf("Content: %v", content.String())

	var profiles Profiles
	err = json.Unmarshal(content.Bytes(), &profiles)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	if len(profiles.Profiles) == 0 {
		log.Printf("Error: no profile found")
		return nil, errors.New("no profile found")
	}

	return profiles.Profiles, nil
}

func getProfileIDs(c *nextdns.Client) ([]string, error) {
	p := make([]string, 0)
	profiles, err := getProfiles(c)
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		p = append(p, profile.ID)
	}

	return p, nil
}

func getProfile(id string, c *http.Client) (string, error) {
	req, err := c.Get("https://api.nextdns.io/profiles/" + id)
	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	if req.StatusCode != 200 {
		log.Printf("Error: %v", req.Status)
		return "", errors.New(req.Status)
	}
	defer req.Body.Close()

	var content bytes.Buffer
	_, err = content.ReadFrom(req.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	return content.String(), nil
}
