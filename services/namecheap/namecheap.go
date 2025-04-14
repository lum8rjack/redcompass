package namecheap

import (
	"strconv"
	"time"

	nc "github.com/lum8rjack/go-namecheap-sdk/v2/namecheap"
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
	client *nc.Client
}

func NewClient(username string, apikey string, ip string) Client {
	c := nc.NewClient(&nc.ClientOptions{
		UserName:   username,
		ApiUser:    username,
		ApiKey:     apikey,
		ClientIp:   ip,
		UseSandbox: false,
	})

	return Client{
		client: c,
	}
}

// NamecheapGetDomains returns a list of all domains for the user
func (c *Client) NamecheapGetDomains() ([]Domain, error) {
	var domains []Domain

	ncresp, err := c.client.Domains.GetList(&nc.DomainsGetListArgs{
		ListType: nc.String("ALL"),
		PageSize: nc.Int(100),
	})
	if err != nil {
		return domains, err
	}

	// If the Domains is nil, no domains were found
	if ncresp.Domains == nil {
		return domains, nil
	}

	for _, x := range *ncresp.Domains {
		newDomain := Domain{}

		// Check the pointers for nil values
		if x.Name != nil {
			newDomain.Name = *x.Name
		}

		if x.Created != nil {
			newDomain.Created = x.Created.Time
		}

		if x.Expires != nil {
			newDomain.Expires = x.Expires.Time
		}

		if x.IsExpired != nil {
			newDomain.IsExpired = *x.IsExpired
		}

		if x.IsLocked != nil {
			newDomain.IsLocked = *x.IsLocked
		}

		if x.AutoRenew != nil {
			newDomain.AutoRenew = *x.AutoRenew
		}

		if x.WhoisGuard != nil {
			newDomain.WhoIsGuard = *x.WhoisGuard == "ENABLED"
		}

		if x.IsOurDNS != nil {
			newDomain.IsOurDNS = *x.IsOurDNS
		}

		domains = append(domains, newDomain)
	}

	return domains, nil
}

// NamecheapGetDomainRecords returns a list of all records for a domain
func (c *Client) NamecheapGetDomainRecords(domain string) ([]Record, error) {
	var records []Record

	ncresp, err := c.client.DomainsDNS.GetHosts(domain)
	if err != nil {
		return records, err
	}

	// If the Hosts is nil, no records were found
	if ncresp.DomainDNSGetHostsResult.Hosts == nil {
		return records, nil
	}

	for _, x := range *ncresp.DomainDNSGetHostsResult.Hosts {
		r := Record{
			Domain: domain,
		}

		// Check the pointers for nil values
		if x.Name != nil {
			r.Name = *x.Name
		}

		if x.Type != nil {
			r.Type = *x.Type
		}

		if x.Address != nil {
			r.Address = *x.Address
		}

		if x.HostId != nil {
			r.HostId = strconv.Itoa(*x.HostId)
		}

		records = append(records, r)
	}

	return records, nil
}
