package porkbun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Domain struct {
	Name       string
	Created    time.Time
	Expires    time.Time
	IsExpired  bool
	AutoRenew  bool
	IsLocked   bool
	WhoIsGuard bool
	IsOurDNS   bool
}

type Record struct {
	Domain  string
	HostId  string
	Name    string
	Type    string
	Address string
}

type Settings struct {
	Username string `json:"username"`
	ApiKey   string `json:"apiKey"`
	IP       string `json:"ipAddress"`
}

type Client struct {
	client           *http.Client
	apiKey           string
	secretKey        string
	perMinuteLimiter *rate.Limiter
}

func NewClient(apikey string, secretkey string) Client {
	webclient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Porkbun doesn't document any specific rate limits but some other libraries use
	// a 1.5 second delay between requests so we'll use that as a starting point.
	return Client{
		client:           webclient,
		apiKey:           apikey,
		secretKey:        secretkey,
		perMinuteLimiter: rate.NewLimiter(rate.Every(time.Second*3/2), 1),
	}
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

// Endpoint: https://api.porkbun.com/api/json/v3/domain/listAll
// Get all domain names in account. Domains are returned in chunks of 1000.

// Method: POST
// Body: {
// 	"secretapikey": "YOUR_SECRET_API_KEY",
// 	"apikey": "YOUR_API_KEY",
// 	"start": "0",
// 	"includeLabels": "yes"
// }

// PorkbunGetDomains returns a list of all domains for the user
func (c *Client) PorkbunGetDomains() ([]Domain, error) {
	var domains []Domain

	// Wait for rate limiter
	if err := c.WaitForRateLimit(); err != nil {
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

	for _, x := range response.Domains {
		newDomain := Domain{}

		// Check the pointers for nil values
		if x.Domain != "" {
			newDomain.Name = x.Domain
		}

		if x.CreateDate != "" {
			newDomain.Created, err = time.Parse(time.RFC3339, x.CreateDate)
			if err != nil {
				return domains, err
			}
		}

		if x.ExpireDate != "" {
			newDomain.Expires, err = time.Parse(time.RFC3339, x.ExpireDate)
			if err != nil {
				return domains, err
			}
		}

		newDomain.IsExpired = false
		if newDomain.Expires.Before(time.Now()) {
			newDomain.IsExpired = true
		}

		newDomain.IsLocked = false
		if x.SecurityLock == "1" {
			newDomain.IsLocked = true
		}

		newDomain.AutoRenew = false
		if x.AutoRenew == 1 {
			newDomain.AutoRenew = true
		}

		newDomain.WhoIsGuard = false
		if x.WhoisPrivacy == "1" {
			newDomain.WhoIsGuard = true
		}

		newDomain.IsOurDNS = true
		if x.NotLocal == 1 {
			newDomain.IsOurDNS = false
		}

		domains = append(domains, newDomain)
	}

	return domains, nil
}

// PorkbunGetDomainRecords returns a list of all records for a domain
// func (c *Client) PorkbunGetDomainRecords(domain string) ([]Record, error) {
// 	var records []Record

// 	// Wait for rate limiter
// 	if err := c.WaitForRateLimit(); err != nil {
// 		return nil, err
// 	}

// 	ncresp, err := c.client.DomainsDNS.GetHosts(domain)
// 	if err != nil {
// 		return records, err
// 	}

// 	// If the Hosts is nil, no records were found
// 	if ncresp.DomainDNSGetHostsResult.Hosts == nil {
// 		return records, nil
// 	}

// 	for _, x := range *ncresp.DomainDNSGetHostsResult.Hosts {
// 		r := Record{
// 			Domain: domain,
// 		}

// 		// Check the pointers for nil values
// 		if x.Name != nil {
// 			r.Name = *x.Name
// 		}

// 		if x.Type != nil {
// 			r.Type = *x.Type
// 		}

// 		if x.Address != nil {
// 			r.Address = *x.Address
// 		}

// 		if x.HostId != nil {
// 			r.HostId = strconv.Itoa(*x.HostId)
// 		}

// 		records = append(records, r)
// 	}

// 	return records, nil
// }

// Wait for the rate limiters
func (c *Client) WaitForRateLimit() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := c.perMinuteLimiter.Wait(ctx); err != nil {
		return errors.New("failed to wait for per minute rate limit")
	}

	return nil
}
