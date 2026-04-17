package virustotal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	client           *http.Client
	apiKey           string
	username         string
	perMinuteLimiter *rate.Limiter
}

type vt_settings struct {
	APIKey   string `json:"apiKey"`
	Username string `json:"username"`
}

func NewClient(settings string) (*Client, error) {
	if settings == "" {
		return nil, errors.New("settings are empty")
	}

	var vtSettings vt_settings
	err := json.Unmarshal([]byte(settings), &vtSettings)
	if err != nil {
		return nil, err
	}
	if len(vtSettings.APIKey) != 64 {
		return nil, errors.New("invalid API key")
	}

	if vtSettings.Username == "" {
		return nil, errors.New("username is empty")
	}

	webclient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// VirusTotal limit: 4/min and 500/day
	return &Client{
		client:           webclient,
		apiKey:           vtSettings.APIKey,
		username:         vtSettings.Username,
		perMinuteLimiter: rate.NewLimiter(rate.Every(time.Minute/4), 1),
	}, nil
}

// GetName returns the name of the scanner
func (c *Client) GetName() string {
	return "VirusTotal"
}

// IsKeyValid checks if the API key is valid
func (c *Client) IsKeyValid() bool {
	req, _ := http.NewRequest("GET", "https://www.virustotal.com/api/v3/metadata", nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", c.apiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return false
	}

	return true
}

// GetDailyQuotaRemaining returns the daily quota remaining
func (c *Client) GetDailyQuotaRemaining() (int, error) {
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/users/%s/api_usage", c.username)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", c.apiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, errors.New("failed to get daily quota remaining")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var response VT_QuotaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, err
	}

	// Get the current date in YYYY-MM-DD format based on UTC time
	now := time.Now().UTC()
	date := now.Format("2006-01-02")

	// Get the daily quota usage for the current date
	dailyUsage := response.Data.Daily[date]
	if dailyUsage == nil {
		return 0, errors.New("failed to get daily quota usage")
	}

	remaining := 500 - dailyUsage["/api/v3/(domains)"]

	return remaining, nil
}

// GetResults returns the domain reputation results for a domain
func (c *Client) GetResults(domain string) ([]byte, error) {
	// Wait for rate limiter
	if err := c.waitForRateLimit(); err != nil {
		return nil, err
	}

	// Setup request
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s", domain)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", c.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		e := fmt.Sprintf("request returned a status code of %d", res.StatusCode)
		return nil, errors.New(e)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Confirm we get a valid JSON response
	var response VT_Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	lastAnalysis, err := convertLastAnalysisResultsToMap(response.Data.Attributes.LastAnalysisResults)
	if err != nil {
		return nil, err
	}

	vtresults := VT_Database{
		Domain:              response.Data.ID,
		VotesHarmless:       response.Data.Attributes.TotalVotes.Harmless,
		VotesMalicious:      response.Data.Attributes.TotalVotes.Malicious,
		Malicious:           response.Data.Attributes.LastAnalysisStats.Malicious,
		Suspicious:          response.Data.Attributes.LastAnalysisStats.Suspicious,
		Undetected:          response.Data.Attributes.LastAnalysisStats.Undetected,
		Harmless:            response.Data.Attributes.LastAnalysisStats.Harmless,
		Timeout:             response.Data.Attributes.LastAnalysisStats.Timeout,
		LastAnalysisResults: lastAnalysis,
	}

	return json.Marshal(vtresults)
}

// convertLastAnalysisResultsToMap turns VT_Response's fixed-field last_analysis_results
// into a map keyed by engine id, using JSON so we do not list every engine by hand.
func convertLastAnalysisResultsToMap(v any) (map[string]VT_Analysis_Engine_Result, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal last_analysis_results: %w", err)
	}
	var m map[string]VT_Analysis_Engine_Result
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("unmarshal last_analysis_results map: %w", err)
	}
	if m == nil {
		m = map[string]VT_Analysis_Engine_Result{}
	}
	return m, nil
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
