package task

import (
	"context"
	"errors"
	"time"
)

type Status string

var (
	ErrInvalidStatus error = errors.New("invalid status")
)

func (s Status) String() string {
	return string(s)
}

func ParseStatus(s string) (*Status, error) {
	switch s {
	case "high", "moderate", "low":
		status := Status(s)
		return &status, nil
	default:
		return nil, ErrInvalidStatus
	}
}

const (
	High     Status = "high"
	Moderate Status = "moderate"
	Low      Status = "low"
)

type Task struct {
	ID          int64
	Name        string
	Description string
	Status      Status
	Label       string
	CreatedAt   time.Time
}

type TaskRepository interface {
	CreateTask(ctx context.Context, t Task) error
	GetTaskById(ctx context.Context, tID int64) (*Task, error)
	UpdateTaskById(ctx context.Context, tID int64) error
	DeleteTaskById(ctx context.Context, tID int64) error
}
