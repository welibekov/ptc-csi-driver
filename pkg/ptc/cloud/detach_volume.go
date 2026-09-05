package cloud

import (
	"context"

	"github.com/welibekov/ptc-csi-driver-ptc/pkg/ptc/util"
	"github.com/welibekov/ptc-csi-driver/pkg/ptc/client/operations"
)

func (c *Client) DetachVolume(ctx context.Context, instanceID, volumeID string) error {
	authorizer, err := c.GetAuthorizer(ctx)
	if err != nil {
		return err
	}

	// 1. Map Machine Spec to PTC API CreateVolumeParams
	params := operations.NewDetachVolumeFromVMParams().WithContext(ctx)
	params.InstanceID = instanceID
	params.VolumeID = volumeID
	params.Force = util.Ptr(true)

	// 2. Trigger Volume creation call
	resp, err := c.apiClient.Operations.DetachVolumeFromVM(params, authorizer)
	if err != nil {
		return err
	}

	if _, err := c.waitForTaskComplete(ctx, resp.GetPayload()); err != nil {
		return err
	}

	return nil
}
