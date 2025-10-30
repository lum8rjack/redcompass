package namecheap

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/lum8rjack/redcompass/services/types"
	nc "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"golang.org/x/time/rate"
)

type Settings struct {
	ApiKey   string `json:"apiKey"`
	Username string `json:"username"`
	IP       string `json:"ipAddress"`
}

type Client struct {
	client           *nc.Client
	perMinuteLimiter *rate.Limiter
	perHourLimiter   *rate.Limiter
	perDayLimiter    *rate.Limiter
}

func NewClient(settings string) (*Client, error) {
	var namecheapSettings Settings
	err := json.Unmarshal([]byte(settings), &namecheapSettings)
	if err != nil {
		return nil, err
	}

	// Check if the settings are valid
	if namecheapSettings.ApiKey == "" || namecheapSettings.IP == "" || namecheapSettings.Username == "" {
		return nil, errors.New("invalid settings")
	}

	// Check if the IP is a valid IP address
	if net.ParseIP(namecheapSettings.IP) == nil {
		return nil, errors.New("invalid IP address")
	}

	c := nc.NewClient(&nc.ClientOptions{
		UserName:   namecheapSettings.Username,
		ApiUser:    namecheapSettings.Username,
		ApiKey:     namecheapSettings.ApiKey,
		ClientIp:   namecheapSettings.IP,
		UseSandbox: false,
	})

	// Namecheap limit: 50/min, 700/hour, and 8000/day across the whole key
	return &Client{
		client:           c,
		perMinuteLimiter: rate.NewLimiter(rate.Every(time.Minute/50), 1),
		perHourLimiter:   rate.NewLimiter(rate.Every(time.Hour/700), 650),
		perDayLimiter:    rate.NewLimiter(rate.Every(24*time.Hour/8000), 7900),
	}, nil
}

// GetName returns the name of the service
func (c *Client) GetName() string {
	return "Namecheap"
}

// GetDomains returns a list of all domains for the user
func (c *Client) GetDomains() ([]types.Domain, error) {
	var domains []types.Domain

	// Wait for rate limiter
	if err := c.waitForRateLimit(); err != nil {
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
		newDomain := types.Domain{}

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

// GetDomainRecords returns a list of all records for a domain
func (c *Client) GetDomainRecords(domain string) ([]types.Record, error) {
	var records []types.Record

	// Wait for rate limiter
	if err := c.waitForRateLimit(); err != nil {
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
		r := types.Record{
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
func (c *Client) waitForRateLimit() error {
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
