package harvestclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// DefaultClient implements HTTP client for Harvest API
type DefaultClient struct {
	client        *http.Client
	accessToken   string
	accountID     string
	userAgentInfo string
}

// NewClient creates a new DefaultClient
func NewClient(accessToken, accountID, userAgentInfo string) *DefaultClient {
	return &DefaultClient{
		client:        &http.Client{},
		accessToken:   accessToken,
		accountID:     accountID,
		userAgentInfo: userAgentInfo,
	}
}

// GetTimeEntries retrieves time entries for a specific date
func (c *DefaultClient) GetTimeEntries(ctx context.Context, date string) ([]byte, error) {
	url := fmt.Sprintf("https://api.harvestapp.com/v2/time_entries?from=%s&to=%s", date, date)
	return c.doRequest(ctx, "GET", url)
}

// GetWhoAmI retrieves the current user's information
func (c *DefaultClient) GetWhoAmI(ctx context.Context) ([]byte, error) {
	return c.doRequest(ctx, "GET", "https://api.harvestapp.com/api/v2/users/me.json")
}

func (c *DefaultClient) doRequest(ctx context.Context, method, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.accessToken)
	req.Header.Add("Harvest-Account-ID", c.accountID)
	req.Header.Add("User-Agent", c.userAgentInfo)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
