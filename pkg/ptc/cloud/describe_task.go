package cloud

import (
	"context"
	"fmt"

	"github.com/welibekov/ptc-csi-driver/pkg/ptc/client/operations"
)

func (c *Client) DescribeTask(ctx context.Context, taskID string) (*Task, error) {
	authorizer, err := c.GetAuthorizer(ctx)
	if err != nil {
		return nil, err
	}

	params := operations.NewDescribeTaskParams()
	params.TaskID = taskID

	resp, err := c.apiClient.Operations.DescribeTaskContext(ctx, params, authorizer)
	if err != nil {
		return nil, fmt.Errorf("error calling API with auth: %w", err)
	}

	task, err := toTask(resp.GetPayload())
	if err != nil {
		return nil, err
	}

	return task, nil
}
