package patreonproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/patreon-db-sync/internal/config"
	"github.com/TicketsBot/patreon-db-sync/internal/utils"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	config config.Config
	client *http.Client
}

func NewClient(config config.Config) *Client {
	return NewClientWithHttpClient(config, &http.Client{
		Timeout: time.Second * 5,
	})
}

func NewClientWithHttpClient(config config.Config, client *http.Client) *Client {
	return &Client{
		config: config,
		client: client,
	}
}

func (c *Client) ListEntitlements(ctx context.Context, legacyOnly bool) (*ListEntitlementsResponse, error) {
	url := c.buildUrl(legacyOnly)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.PatreonProxy.AuthToken))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var response ListEntitlementsResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) buildUrl(legacyOnly bool) *url.URL {
	url := utils.Ptr(*c.config.PatreonProxy.RootUrl) // Clone URL
	url.Path = "/all"
	url.Query().Set("legacyOnly", fmt.Sprintf("%t", legacyOnly))
	return url
}
