package types

import "time"

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
	TTL     string
}

type Service interface {
	GetName() string
	GetDomains() ([]Domain, error)
	GetDomainRecords(domain string) ([]Record, error)
}
