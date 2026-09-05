package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/welibekov/ptc-csi-driver/pkg/ptc/util"
)

type TaskID string

// String converts TaskID to a string.
func (t TaskID) String() string {
	return string(t)
}

// Task represents a unit of work with associated metadata.
type Task struct {
	ID         TaskID          `json:"task-id"`               // Unique identifier for the task
	Title      string          `json:"title,omitempty"`       // Title or description of the task
	ErrMessage *string         `json:"err_message,omitempty"` // Optional error message if the task fails
	Status     TaskStatus      `json:"status"`                // Current status of the task
	CreatedAt  time.Time       `json:"created_at,omitempty"`  // Timestamp when the task was created
	UpdatedAt  time.Time       `json:"updated_at,omitempty"`  // Timestamp when the task was last updated
	Output     json.RawMessage `json:"output,omitempty"`      // Output result of handler execution
}

// Done returns true if the task status is completed.
func (t *Task) Done() bool {
	return t.Status == TaskCompleted
}

// Failed returns true if the task status is failed.
func (t *Task) Failed() bool {
	return t.Status == TaskFailed
}

// TaskStatus defines the possible states of a task.
type TaskStatus string

const (
	TaskAccepted  TaskStatus = "accepted"  // Task has been accepted for processing
	TaskCompleted TaskStatus = "completed" // Task has been successfully completed
	TaskFailed    TaskStatus = "failed"    // Task has failed during processing
	TaskRunning   TaskStatus = "running"   // Task is currently being processed
)

// Helper functions
func toTask(input any) (*Task, error) {
	task := Task{}
	if err := rest2Task(input, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// rest2Task converts a raw payload into a types.Task, stripping metadata fields starting with "_".
func rest2Task(src interface{}, dest *Task) error {
	if src == nil || dest == nil {
		return errors.New("source payload and destination task pointer must not be nil")
	}

	// 1. Marshal source to JSON bytes
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("failed to marshal source payload: %w", err)
	}

	// 2. Unmarshal into map to filter ignored keys (starting with "_")
	var rawMap map[string]interface{}
	if err := json.Unmarshal(srcBytes, &rawMap); err != nil {
		return fmt.Errorf("failed to unmarshal payload to map: %w", err)
	}

	filteredMap := make(map[string]interface{}, len(rawMap))
	for key, value := range rawMap {
		if !strings.HasPrefix(key, "_") {
			filteredMap[key] = value
		}
	}

	// 3. Remarshal filtered map
	filteredBytes, err := json.Marshal(filteredMap)
	if err != nil {
		return fmt.Errorf("failed to marshal filtered payload: %w", err)
	}

	// 4. Decode into destination struct
	decoder := json.NewDecoder(bytes.NewReader(filteredBytes))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("failed to decode payload into Task: %w", err)
	}

	return nil
}

// Wait blocks until the asynchronous task completes or fails, polling the task state periodically.
//
// It uses the provided describeTaskFn to fetch the current state of the task
// at regular intervals (currently set to every second). The function returns
// nil if the task completes successfully, or an error if the task fails, the
// context is canceled, or an error occurs while fetching the task state.
//
// Parameters:
//   - ctx: Context for cancellation and deadlines.
//   - describeTaskFn: A function that retrieves the current state of the task.
//
// Returns:
//   - error: Returns nil on successful task completion, or an error if the task fails,
//     the context is canceled, or an error occurs during task state retrieval.
func (t *Task) Wait(ctx context.Context, describTaskFn func(*Task) error) error {
	// TODO: The ticker interval can be taken from a config parameter.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := describTaskFn(t); err != nil {
				return err
			}

			if t.Done() {
				return nil
			}

			if t.Failed() {
				return fmt.Errorf("failed: %v", util.Deref(t.ErrMessage))
			}
		}
	}
}

func (t Task) String() string {
	taskBytes, _ := json.MarshalIndent(t, "", "  ")

	return string(taskBytes)
}
