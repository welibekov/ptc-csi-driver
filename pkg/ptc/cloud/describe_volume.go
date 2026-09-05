package cloud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/welibekov/ptc-csi-driver/pkg/ptc/client/operations"
)

type Volume struct {
	// The instance ID of the VM to which the volume is attached, if any
	AttachedTo string `json:"attached_to,omitempty"`

	// The timestamp when the volume was created
	CreatedAt string `json:"created_at,omitempty"`

	// The name of the volume
	Name string `json:"name,omitempty"`

	// The size of the volume
	Size string `json:"size,omitempty"`

	// The current status of the volume
	Status string `json:"status,omitempty"`

	// Tags associated with the volume
	Tags []string `json:"tags"`

	// The type of the volume
	Type string `json:"type,omitempty"`

	// The unique identifier of the volume
	VolumeID string `json:"volume_id,omitempty"`
}

func (c *Client) DescribeVolume(ctx context.Context, volumeID string) (*Volume, error) {
	authorizer, err := c.GetAuthorizer(ctx)
	if err != nil {
		return nil, err
	}

	params := operations.NewDescribeVolumeParams()
	resp, err := c.apiClient.Operations.DescribeVolumeContext(ctx, params, authorizer)
	if err != nil {
		return nil, fmt.Errorf("error calling API with auth: %w", err)
	}

	payloadBytes, err := json.Marshal(resp.GetPayload())
	if err != nil {
		return nil, fmt.Errorf("failed to get volume list: %w", err)
	}

	var volume Volume
	if err := json.Unmarshal(payloadBytes, &volumes); err != nil {
		return nil, fmt.Errorf("failed to get volume list: %w", err)
	}

	return &volume, nil
}
