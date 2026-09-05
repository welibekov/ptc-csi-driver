package cloud

import (
	"github.com/welibekov/ptc-csi-driver/pkg/ptc/client"
	swaggerClient "github.com/welibekov/ptc-csi-driver/pkg/ptc/client"
)

// PTCCredentials holds authenticated session details for PTC API.
type PTCCredentials struct {
	OriginalHostname string
	Scheme           string
	Hostname         string
	Username         string
	Password         string
	TenantID         string
	ClientID         string
}

type Client struct {
	creds        *PTCCredentials
	tokenManager *TokenManager
	baseURL      string
	apiClient    *swaggerClient.Ptc // Pre-configured SDK client
}

func NewClient(creds *PTCCredentials, tm *TokenManager) *Client {
	cfg := client.DefaultTransportConfig().
		WithHost(creds.Hostname).
		WithSchemes([]string{creds.Scheme}).
		WithBasePath("")

	c := &Client{
		creds:        creds,
		apiClient:    client.NewHTTPClientWithConfig(nil, cfg),
		baseURL:      creds.Hostname,
		tokenManager: tm,
	}

	return c
}
