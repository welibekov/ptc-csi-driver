package cloud

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"golang.org/x/oauth2"
)

// newAPIKeyAuthWriter creates a new ClientAuthInfoWriter for the given OAuth2 token.
func newAPIKeyAuthWriter(token *oauth2.Token) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(
		req runtime.ClientRequest,
		reg strfmt.Registry,
	) error {
		return req.SetHeaderParam("x-ptc-token", token.AccessToken) // Set the token in the request header.
	})
}
