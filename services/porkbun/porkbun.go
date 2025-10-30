package porkbun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lum8rjack/redcompass/services/types"
	"golang.org/x/time/rate"
)

type Settings struct {
	ApiKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
}

type Client struct {
	client           *http.Client
	apiKey           string
	secretKey        string
	perMinuteLimiter *rate.Limiter
}

func NewClient(settings string) (*Client, error) {
	var porkbunSettings Settings
	err := json.Unmarshal([]byte(settings), &porkbunSettings)
	if err != nil {
		return nil, err
	}

	// Check if the settings are valid
	if porkbunSettings.ApiKey == "" || porkbunSettings.SecretKey == "" {
		return nil, errors.New("invalid settings")
	}

	if !strings.HasPrefix(porkbunSettings.ApiKey, "pk") || !strings.HasPrefix(porkbunSettings.SecretKey, "sk") {
		return nil, errors.New("invalid settings")
	}

	// Create the web client
	webclient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Porkbun doesn't document any specific rate limits but some other libraries use
	// a 1.5 second delay between requests so we'll use that as a starting point.
	return &Client{
		client:           webclient,
		apiKey:           porkbunSettings.ApiKey,
		secretKey:        porkbunSettings.SecretKey,
		perMinuteLimiter: rate.NewLimiter(rate.Every(time.Second*3/2), 1),
	}, nil
}

// GetName returns the name of the service
func (c *Client) GetName() string {
	return "Porkbun"
}

type PorkbunListAllRequest struct {
	SecretApiKey  string `json:"secretapikey"`
	ApiKey        string `json:"apikey"`
	Start         string `json:"start"`
	IncludeLabels string `json:"includeLabels"`
}

type DomainListResponse struct {
	Status  string `json:"status"`
	Domains []struct {
		Domain       string `json:"domain"`
		Status       string `json:"status"`
		Tld          string `json:"tld"`
		CreateDate   string `json:"createDate"`
		ExpireDate   string `json:"expireDate"`
		SecurityLock string `json:"securityLock"`
		WhoisPrivacy string `json:"whoisPrivacy"`
		AutoRenew    int    `json:"autoRenew"`
		NotLocal     int    `json:"notLocal"`
		Labels       []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Color string `json:"color"`
		} `json:"labels"`
	} `json:"domains"`
}

// GetDomains returns a list of all domains for the user
func (c *Client) GetDomains() ([]types.Domain, error) {
	var domains []types.Domain

	// Wait for rate limiter
	if err := c.waitForRateLimit(); err != nil {
		return nil, err
	}

	requestBody := PorkbunListAllRequest{
		SecretApiKey:  c.secretKey,
		ApiKey:        c.apiKey,
		Start:         "0",
		IncludeLabels: "yes",
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return domains, err
	}

	reader := bytes.NewReader(requestBodyBytes)

	resp, err := c.client.Post(
		"https://api.porkbun.com/api/json/v3/domain/listAll",
		"application/json",
		reader,
	)
	if err != nil {
		return domains, err
	}
	defer resp.Body.Close()

	var response DomainListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return domains, err
	}

	if response.Status != "SUCCESS" {
		return domains, errors.New(response.Status)
	}

	layout := "2006-01-02 15:04:05"
	for _, x := range response.Domains {
		newDomain := types.Domain{}

		// Check the values
		if x.Domain != "" {
			newDomain.Name = x.Domain
		}

		if x.CreateDate != "" {
			newDomain.Created, err = time.Parse(layout, x.CreateDate)
			if err != nil {
				return domains, err
			}
		}

		if x.ExpireDate != "" {
			newDomain.Expires, err = time.Parse(layout, x.ExpireDate)
			if err != nil {
				return domains, err
			}
		}

		newDomain.IsExpired = false
		if newDomain.Expires.Before(time.Now()) {
			newDomain.IsExpired = true
		}

		newDomain.IsLocked = false

		newDomain.AutoRenew = false
		if x.AutoRenew == 1 {
			newDomain.AutoRenew = true
		}

		newDomain.WhoIsGuard = false
		if x.WhoisPrivacy == "1" {
			newDomain.WhoIsGuard = true
		}

		newDomain.IsOurDNS = true
		usingPorkbunNameservers, err := c.isUsingPorkbunNameservers(x.Domain)
		if err != nil {
			return domains, err
		}
		newDomain.IsOurDNS = usingPorkbunNameservers

		domains = append(domains, newDomain)
	}

	return domains, nil
}

type PorkbunRetrieveRecordsRequest struct {
	SecretApiKey string `json:"secretapikey"`
	ApiKey       string `json:"apikey"`
}

type PorkbunRetrieveRecordsResponse struct {
	Status  string `json:"status"`
	Records []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
		TTL     string `json:"ttl"`
		Prio    string `json:"prio"`
		Notes   string `json:"notes"`
	} `json:"records"`
}

// GetDomainRecords returns a list of all records for a domain
func (c *Client) GetDomainRecords(domain string) ([]types.Record, error) {
	var records []types.Record

	// Wait for rate limiter
	if err := c.waitForRateLimit(); err != nil {
		return records, err
	}

	requestBody := PorkbunRetrieveRecordsRequest{
		SecretApiKey: c.secretKey,
		ApiKey:       c.apiKey,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return records, err
	}

	reader := bytes.NewReader(requestBodyBytes)

	url := fmt.Sprintf("https://api.porkbun.com/api/json/v3/dns/retrieve/%s", domain)

	resp, err := c.client.Post(
		url,
		"application/json",
		reader,
	)
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()

	var response PorkbunRetrieveRecordsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return records, err
	}

	if response.Status != "SUCCESS" {
		return records, errors.New(response.Status)
	}

	for _, x := range response.Records {
		r := types.Record{
			Domain: domain,
		}

		// Check the values
		if x.Name != "" {
			r.Name = x.Name
		}

		if x.Type != "" {
			r.Type = x.Type
		}

		if x.Content != "" {
			r.Address = x.Content
		}

		if x.TTL != "" {
			r.TTL = x.TTL
		}

		if x.ID != "" {
			r.HostId = x.ID
		}

		records = append(records, r)
	}
	return records, nil
}

type PorkbunGetNsResponse struct {
	Status      string   `json:"status"`
	Nameservers []string `json:"ns"`
}

// Check nameservers for a domain
func (c *Client) isUsingPorkbunNameservers(domain string) (bool, error) {
	// Wait for rate limiter
	if err := c.waitForRateLimit(); err != nil {
		return false, err
	}

	requestBody := PorkbunRetrieveRecordsRequest{
		SecretApiKey: c.secretKey,
		ApiKey:       c.apiKey,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return false, err
	}

	reader := bytes.NewReader(requestBodyBytes)

	url := fmt.Sprintf("https://api.porkbun.com/api/json/v3/domain/getNs/%s", domain)

	resp, err := c.client.Post(
		url,
		"application/json",
		reader,
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var response PorkbunGetNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, err
	}

	if response.Status != "SUCCESS" {
		return false, errors.New(response.Status)
	}

	if len(response.Nameservers) == 0 {
		return false, nil
	}

	isPorkbunNameserver := strings.HasSuffix(response.Nameservers[0], ".porkbun.com")

	return isPorkbunNameserver, nil
}

// Wait for the rate limiters
func (c *Client) waitForRateLimit() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := c.perMinuteLimiter.Wait(ctx); err != nil {
		return errors.New("failed to wait for per minute rate limit")
	}

	return nil
}
