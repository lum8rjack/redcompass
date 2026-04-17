package types

type Scanner interface {
	GetName() string
	IsKeyValid() bool
	GetResults(domain string) ([]byte, error)
	GetDailyQuotaRemaining() (int, error)
}
