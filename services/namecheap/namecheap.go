package namecheap

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const namecheapEndpoint string = "https://api.namecheap.com/xml.response"

type ApiResponseGetDomains struct {
	XMLName          xml.Name `xml:"ApiResponse"`
	Text             string   `xml:",chardata"`
	Status           string   `xml:"Status,attr"`
	Xmlns            string   `xml:"xmlns,attr"`
	Errors           string   `xml:"Errors"`
	Warnings         string   `xml:"Warnings"`
	RequestedCommand string   `xml:"RequestedCommand"`
	CommandResponse  struct {
		Text                string `xml:",chardata"`
		Type                string `xml:"Type,attr"`
		DomainGetListResult struct {
			Text   string `xml:",chardata"`
			Domain []struct {
				Text       string `xml:",chardata"`
				ID         string `xml:"ID,attr"`
				Name       string `xml:"Name,attr"`
				User       string `xml:"User,attr"`
				Created    string `xml:"Created,attr"`
				Expires    string `xml:"Expires,attr"`
				IsExpired  string `xml:"IsExpired,attr"`
				IsLocked   string `xml:"IsLocked,attr"`
				AutoRenew  string `xml:"AutoRenew,attr"`
				WhoisGuard string `xml:"WhoisGuard,attr"`
				IsPremium  string `xml:"IsPremium,attr"`
				IsOurDNS   string `xml:"IsOurDNS,attr"`
			} `xml:"Domain"`
		} `xml:"DomainGetListResult"`
		Paging struct {
			Text        string `xml:",chardata"`
			TotalItems  string `xml:"TotalItems"`
			CurrentPage string `xml:"CurrentPage"`
			PageSize    string `xml:"PageSize"`
		} `xml:"Paging"`
	} `xml:"CommandResponse"`
	Server            string `xml:"Server"`
	GMTTimeDifference string `xml:"GMTTimeDifference"`
	ExecutionTime     string `xml:"ExecutionTime"`
}

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

type ApiResponseGetRecords struct {
	XMLName          xml.Name `xml:"ApiResponse"`
	Text             string   `xml:",chardata"`
	Status           string   `xml:"Status,attr"`
	Xmlns            string   `xml:"xmlns,attr"`
	Errors           string   `xml:"Errors"`
	Warnings         string   `xml:"Warnings"`
	RequestedCommand string   `xml:"RequestedCommand"`
	CommandResponse  struct {
		Text                    string `xml:",chardata"`
		Type                    string `xml:"Type,attr"`
		DomainDNSGetHostsResult struct {
			Text          string `xml:",chardata"`
			Domain        string `xml:"Domain,attr"`
			EmailType     string `xml:"EmailType,attr"`
			IsUsingOurDNS string `xml:"IsUsingOurDNS,attr"`
			Host          []struct {
				Text               string `xml:",chardata"`
				HostId             string `xml:"HostId,attr"`
				Name               string `xml:"Name,attr"`
				Type               string `xml:"Type,attr"`
				Address            string `xml:"Address,attr"`
				MXPref             string `xml:"MXPref,attr"`
				TTL                string `xml:"TTL,attr"`
				AssociatedAppTitle string `xml:"AssociatedAppTitle,attr"`
				FriendlyName       string `xml:"FriendlyName,attr"`
				IsActive           string `xml:"IsActive,attr"`
				IsDDNSEnabled      string `xml:"IsDDNSEnabled,attr"`
			} `xml:"host"`
		} `xml:"DomainDNSGetHostsResult"`
	} `xml:"CommandResponse"`
	Server            string `xml:"Server"`
	GMTTimeDifference string `xml:"GMTTimeDifference"`
	ExecutionTime     string `xml:"ExecutionTime"`
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
	username string
	apikey   string
	ip       string
}

func NewClient(username string, apikey string, ip string) Client {
	c := Client{
		username: username,
		apikey:   apikey,
		ip:       ip,
	}

	return c
}

func (c *Client) NamecheapGetDomains() ([]Domain, error) {
	var apiresponse ApiResponseGetDomains
	var domains []Domain

	url := fmt.Sprintf("%s?Command=namecheap.domains.getList&PageSize=100&UserName=%s&ApiUser=%s&ApiKey=%s&ClientIp=%s", namecheapEndpoint, c.username, c.username, c.apikey, c.ip)

	// Namecheap GET request to get the list of all domains
	resp, err := http.Get(url)
	if err != nil {
		return domains, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e := fmt.Sprintf("status code %d", resp.StatusCode)
		return domains, errors.New(e)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domains, err
	}

	err = xml.Unmarshal(body, &apiresponse)
	if err != nil {
		return domains, err
	}

	pageSize, err := strconv.Atoi(apiresponse.CommandResponse.Paging.PageSize)
	if err != nil {
		return domains, err
	}

	totalItems, err := strconv.Atoi(apiresponse.CommandResponse.Paging.TotalItems)
	if err != nil {
		return domains, err
	}

	if pageSize < totalItems {
		// FIXME: do more requests
		fmt.Println("FIXME: Do more Namecheap requests")
	}

	for _, x := range apiresponse.CommandResponse.DomainGetListResult.Domain {
		newDomain := Domain{
			Name: x.Name,
		}

		c, err := time.Parse("01/02/2006", x.Created)
		if err != nil {
			continue
		}
		newDomain.Created = c

		e, err := time.Parse("01/02/2006", x.Expires)
		if err != nil {
			continue
		}
		newDomain.Expires = e

		if x.IsExpired == "false" {
			newDomain.IsExpired = false
		} else {
			newDomain.IsExpired = true
		}

		if x.IsLocked == "false" {
			newDomain.IsLocked = false
		} else {
			newDomain.IsLocked = true
		}

		if x.AutoRenew == "false" {
			newDomain.AutoRenew = false
		} else {
			newDomain.AutoRenew = true
		}

		if x.WhoisGuard == "ENABLED" {
			newDomain.WhoIsGuard = true
		} else {
			newDomain.WhoIsGuard = false
		}

		if x.IsOurDNS == "false" {
			newDomain.IsOurDNS = false
		} else {
			newDomain.IsOurDNS = true
		}

		domains = append(domains, newDomain)
	}

	return domains, nil
}

func (c *Client) NamecheapGetDomainRecords(domain string) ([]Record, error) {
	var apiresponse ApiResponseGetRecords
	var records []Record

	parts := strings.Split(domain, ".")
	sld := parts[0]
	tld := strings.Join(parts[1:], ".")

	url := fmt.Sprintf("%s?Command=namecheap.domains.dns.getHosts&UserName=%s&ApiUser=%s&ApiKey=%s&ClientIp=%s&SLD=%s&TLD=%s", namecheapEndpoint, c.username, c.username, c.apikey, c.ip, sld, tld)

	// Namecheap GET request to get the list of all domains
	resp, err := http.Get(url)
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e := fmt.Sprintf("status code %d", resp.StatusCode)
		return records, errors.New(e)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}

	err = xml.Unmarshal(body, &apiresponse)
	if err != nil {
		return records, err
	}

	for _, x := range apiresponse.CommandResponse.DomainDNSGetHostsResult.Host {
		records = append(records, Record{
			Domain:  domain,
			HostId:  x.HostId,
			Name:    x.Name,
			Type:    x.Type,
			Address: x.Address,
		})
	}

	return records, nil
}
