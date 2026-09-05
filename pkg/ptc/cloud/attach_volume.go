package cloud

import (
	"context"

	"github.com/welibekov/ptc-csi-driver/pkg/ptc/client/operations"
)

func (c *Client) AttachVolume(ctx context.Context, instanceID, volumeID string) error {
	authorizer, err := c.GetAuthorizer(ctx)
	if err != nil {
		return err
	}

	// 1. Map Machine Spec to PTC API CreateVolumeParams
	params := operations.NewAttachVolumeToVMParams().WithContext(ctx)
	params.InstanceID = instanceID
	params.VolumeID = volumeID

	// 2. Trigger Volume creation call
	resp, err := c.apiClient.Operations.AttachVolumeToVM(params, authorizer)
	if err != nil {
		return err
	}

	if _, err := c.waitForTaskComplete(ctx, resp.GetPayload()); err != nil {
		return err
	}

	return nil
}
