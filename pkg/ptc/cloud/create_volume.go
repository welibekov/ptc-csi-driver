package cloud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/welibekov/ptc-csi-driver/pkg/ptc/client/operations"
	"github.com/welibekov/ptc-csi-driver/pkg/ptc/util"
)

func (c *Client) CreateVolume(ctx context.Context, params *operations.CreateVolumeParams) (*string, error) {
	authorizer, err := c.GetAuthorizer(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.apiClient.Operations.CreateVolumeContext(ctx, params, authorizer)
	if err != nil {
		return nil, err
	}

	task, err := c.waitForTaskComplete(ctx, resp.GetPayload())
	if err != nil {
		return nil, err
	}

	var taskOutput struct {
		VolumeID string `json:"volume-id"`
	}

	// Extract generated Volume ID from Task Output JSON payload
	if len(task.Output) > 0 {
		// Step 1: Decode the outer JSON string
		var rawJSONString string
		if err := json.Unmarshal(task.Output, &rawJSONString); err == nil {
			// Step 2: Decode the inner JSON object
			if err := json.Unmarshal([]byte(rawJSONString), &taskOutput); err != nil {
				return nil, fmt.Errorf("failed to unmarshal inner task output payload: %w", err)
			}
		} else {
			// Fallback: Attempt single-step unmarshal if task output is direct JSON
			if err := json.Unmarshal(task.Output, &taskOutput); err != nil {
				return nil, fmt.Errorf("failed to unmarshal single-step task if direct JSON: %w", err)
			}
		}
	}

	return util.Ptr(taskOutput.VolumeID), nil
}
