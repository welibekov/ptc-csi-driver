package cloud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/welibekov/ptc-csi-driver/pkg/ptc/client/operations"
)

func (c *Client) ListVolumes(ctx context.Context) (*[]Volume, error) {
	authorizer, err := c.GetAuthorizer(ctx)
	if err != nil {
		return nil, err
	}

	params := operations.NewListVolumesParams()
	resp, err := c.apiClient.Operations.ListVolumesContext(ctx, params, authorizer)
	if err != nil {
		return nil, fmt.Errorf("error calling API with auth: %w", err)
	}

	payloadBytes, err := json.Marshal(resp.GetPayload())
	if err != nil {
		return nil, fmt.Errorf("failed to get volume list: %w", err)
	}

	var volumes []Volume
	if err := json.Unmarshal(payloadBytes, &volumes); err != nil {
		return nil, fmt.Errorf("failed to get volume list: %w", err)
	}

	return &volumes, nil
}
