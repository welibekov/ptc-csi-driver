package cloud

import (
	"context"

	"github.com/welibekov/ptc-csi-driver/pkg/ptc/client/operations"
)

func (c *Client) DeleteVolume(ctx context.Context, volumeID string) error {
	authorizer, err := c.GetAuthorizer(ctx)
	if err != nil {
		return err
	}

	// 1. Map Machine Spec to PTC API CreateVolumeParams
	params := operations.NewDeleteVolumeParams().WithContext(ctx)
	params.VolumeID = volumeID

	// 2. Trigger Volume creation call
	resp, err := c.apiClient.Operations.DeleteVolume(params, authorizer)
	if err != nil {
		return err
	}

	if _, err := c.waitForTaskComplete(ctx, resp.GetPayload()); err != nil {
		return err
	}

	return nil
}
