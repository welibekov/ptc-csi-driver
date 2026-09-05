package cloud

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"

	"github.com/go-openapi/runtime"
	"github.com/welibekov/ptc-csi-driver/pkg/ptc/client/operations"
	"github.com/welibekov/ptc-csi-driver/pkg/ptc/util"
)

// Wrapper around TokenManager GetToken() method.
func (c *Client) GetToken(ctx context.Context) (*oauth2.Token, error) {
	return c.tokenManager.GetToken(ctx, c)
}

func (c *Client) GetAuthorizer(ctx context.Context) (runtime.ClientAuthInfoWriter, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	return newAPIKeyAuthWriter(token), nil
}

// login uses the stored credentials to request a new token
func (c *Client) login(ctx context.Context) (*oauth2.Token, error) {
	params := &operations.GetTokenParams{
		TenantID:  c.creds.TenantID,
		ClientID:  c.creds.ClientID,
		GrantType: "password",
		Username:  util.Ptr(c.creds.Username),
		Password:  util.Ptr(c.creds.Password),
	}

	resp, err := c.apiClient.Operations.GetTokenContext(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("login token request failed: %w", err)
	}

	return parseOAuth2Token(resp.Payload)
}

// refreshToken uses an existing refresh token string to obtain a new token
func (c *Client) refreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	params := &operations.GetTokenParams{
		TenantID:     c.creds.TenantID,
		ClientID:     c.creds.ClientID,
		GrantType:    "refresh_token",
		RefreshToken: &refreshToken,
	}

	resp, err := c.apiClient.Operations.GetTokenContext(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("refresh token request failed: %w", err)
	}

	return parseOAuth2Token(resp.Payload)
}

func parseOAuth2Token(payload interface{}) (*oauth2.Token, error) {
	tokenJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("couldn't serialize token payload: %w", err)
	}

	var newToken oauth2.Token
	if err := json.Unmarshal(tokenJSON, &newToken); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal into oauth2.Token: %w", err)
	}

	return &newToken, nil
}
