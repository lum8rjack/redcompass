package namecheap

import (
	"strconv"
	"time"

	nc "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
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

	for _, x := range *ncresp.Domains {
		newDomain := Domain{
			Name:       *x.Name,
			Created:    x.Created.Time,
			Expires:    x.Expires.Time,
			IsExpired:  *x.IsExpired,
			IsLocked:   *x.IsLocked,
			AutoRenew:  *x.AutoRenew,
			WhoIsGuard: *x.WhoisGuard == "ENABLED",
			IsOurDNS:   *x.IsOurDNS,
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

	for _, x := range *ncresp.DomainDNSGetHostsResult.Hosts {
		records = append(records, Record{
			Domain:  domain,
			HostId:  strconv.Itoa(*x.HostId),
			Name:    *x.Name,
			Type:    *x.Type,
			Address: *x.Address,
		})
	}

	return records, nil
}
