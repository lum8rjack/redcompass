package namecheap

import (
	"context"
	"errors"
	"strconv"
	"time"

	nc "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
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
	client           *nc.Client
	perMinuteLimiter *rate.Limiter
	perHourLimiter   *rate.Limiter
	perDayLimiter    *rate.Limiter
}

func NewClient(username string, apikey string, ip string) Client {
	c := nc.NewClient(&nc.ClientOptions{
		UserName:   username,
		ApiUser:    username,
		ApiKey:     apikey,
		ClientIp:   ip,
		UseSandbox: false,
	})

	// Namecheap limit: 50/min, 700/hour, and 8000/day across the whole key
	return Client{
		client:           c,
		perMinuteLimiter: rate.NewLimiter(rate.Every(time.Minute/50), 1),
		perHourLimiter:   rate.NewLimiter(rate.Every(time.Hour/700), 650),
		perDayLimiter:    rate.NewLimiter(rate.Every(24*time.Hour/8000), 7900),
	}
}

// NamecheapGetDomains returns a list of all domains for the user
func (c *Client) NamecheapGetDomains() ([]Domain, error) {
	var domains []Domain

	// Wait for rate limiter
	if err := c.WaitForRateLimit(); err != nil {
		return nil, err
	}

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

	// Wait for rate limiter
	if err := c.WaitForRateLimit(); err != nil {
		return nil, err
	}

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

// Wait for the rate limiters
func (c *Client) WaitForRateLimit() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := c.perMinuteLimiter.Wait(ctx); err != nil {
		return errors.New("failed to wait for per minute rate limit")
	}

	if err := c.perHourLimiter.Wait(ctx); err != nil {
		return errors.New("failed to wait for per hour rate limit")
	}

	if err := c.perDayLimiter.Wait(ctx); err != nil {
		return errors.New("failed to wait for per day rate limit")
	}

	return nil
}
