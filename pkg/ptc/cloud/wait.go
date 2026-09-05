package cloud

import (
	"context"
	"fmt"
)

func (c *Client) waitForTaskComplete(ctx context.Context, payload any) (*Task, error) {
	task, err := toTask(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task response: %w", err)
	}

	err = task.Wait(ctx, func(target *Task) error {
		freshTask, err := c.DescribeTask(ctx, task.ID.String())
		if err != nil {
			return err
		}
		*target = *freshTask
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error waiting for VM creation task: %w", err)
	}

	return task, nil
}
