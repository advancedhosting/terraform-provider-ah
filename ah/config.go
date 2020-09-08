package ah

import (
	"github.com/advancedhosting/advancedhosting-api-go/ah"
)

// Config represents provider's configuration
type Config struct {
	Token       string
	APIEndpoint string
}

// Client returns a new client to communicate with AH Cloud
func (c *Config) Client() (*ah.APIClient, error) {
	clientOptions := &ah.ClientOptions{
		Token:   c.Token,
		BaseURL: c.APIEndpoint,
	}

	return ah.NewAPIClient(clientOptions)

}
