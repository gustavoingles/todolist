package database

import (
	"context"
	"fmt"
	"todo-list/task"

	"gorm.io/gorm"
)

type TaskStatusDB struct {
	gorm.Model
	StatusName string `gorm:"column:status_name;not null; unique"`
}

type TaskLabelDB struct {
	gorm.Model
	LabelName string `gorm:"column: label_name;size(15); not null; unique; check: label_name >= 1"`
}

type TaskDB struct {
	gorm.Model
	Name        string `gorm:"not null; uniqueIndex; size(30); check: name >= 1"`
	Description string `gorm:"not null"`
	Status      TaskStatusDB
	Label       []TaskLabelDB
}

type TaskStore struct {
	db *gorm.DB
}

func (ts TaskStore) MigrateTaskStatus(ctx context.Context, statusToCreate *[]TaskStatusDB) error {
	return gorm.G[TaskStatusDB](ts.db).CreateInBatches(ctx, statusToCreate, 1)
}

func (ts TaskStore) CreateTask(ctx context.Context, t task.Task) error {
	parsedTaskLabels := make([]TaskLabelDB, 0, len(t.Label))
	for _, l := range t.Label {
		tldb := TaskLabelDB{
			LabelName: l,
		}
		parsedTaskLabels = append(parsedTaskLabels, tldb)
	}
	parsedTaskStatus := t.Status.String()

	return gorm.G[TaskDB](ts.db).Create(ctx, &TaskDB{
		Name:        t.Name,
		Description: t.Description,
		Status: TaskStatusDB{
			StatusName: parsedTaskStatus,
		},
		Label: parsedTaskLabels,
	})
}

func (ts TaskStore) GetTaskById(ctx context.Context, tID int64) (*task.Task, error) {
	taskDB, err := gorm.G[TaskDB](ts.db).Where("id = ?", tID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while fetchin task by its id: %w", err)
	}

	mappedTaskStatus, err := task.ParseStatus(taskDB.Status.StatusName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task status string into 'Status' type: %w", err)
	}
	mappedTaskLabel := make([]string, 0, len(taskDB.Label))
	for _, l := range taskDB.Label {
		mappedTaskLabel = append(mappedTaskLabel, l.LabelName)
	}

	return &task.Task{
		ID:          int64(taskDB.ID),
		Name:        taskDB.Name,
		Description: taskDB.Description,
		Label:       mappedTaskLabel,
		Status:      *mappedTaskStatus,
		CreatedAt:   taskDB.CreatedAt,
	}, nil
}

func (ts TaskStore) GetAllTasks(ctx context.Context) ([]task.Task, error) {
	tasksDB, err := gorm.G[TaskDB](ts.db).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while fetching all tasks: %w", err)
	}

	tasks := make([]task.Task, 0, len(tasksDB))
	for _, tdb := range tasksDB {
		mappedTaskStatus, err := task.ParseStatus(tdb.Status.StatusName)
		if err != nil {
			return nil, fmt.Errorf("failed to parse task status string into 'Status' type: %w", err)
		}

		mappedTaskLabel := make([]string, 0, len(tdb.Label))
		for _, l := range tdb.Label {
			mappedTaskLabel = append(mappedTaskLabel, l.LabelName)
		}

		tasks = append(tasks, task.Task{
			ID:          int64(tdb.ID),
			Name:        tdb.Name,
			Description: tdb.Description,
			Label:       mappedTaskLabel,
			Status:      *mappedTaskStatus,
			CreatedAt:   tdb.CreatedAt,
		})
	}

	return tasks, nil
}

func (ts TaskStore) DeleteTaskById(ctx context.Context, tID int64) error {
	_, err := gorm.G[TaskDB](ts.db).Where("id = ?", tID).Delete(ctx)
	return err
}

func (ts TaskStore) UpdateTaskById(ctx context.Context, tID int64, dataToChange task.TaskUpdateData) error {
	err := ts.db.Transaction(func(tx *gorm.DB) error {
		taskDB, err := gorm.G[TaskDB](tx).Where("id = ?", tID).First(ctx)
		if err != nil {
			return err
		}

		updateFn := func(tdb *TaskDB, dtc task.TaskUpdateData) TaskDB {
			tdb.Name = dtc.NewName
			tdb.Description = dtc.NewDescription
			tdb.Status = TaskStatusDB{StatusName: dtc.NewStatus.String()}

			newLabels := make([]TaskLabelDB, 0)
			for _, l := range dtc.NewLabel {
				nl := TaskLabelDB{LabelName: l}
				newLabels = append(newLabels, nl)
			}
			tdb.Label = newLabels
			return *tdb
		}

		_, err = gorm.G[TaskDB](tx).Where("id = ?", tID).Updates(ctx, updateFn(&taskDB, dataToChange))
		return err
	})
	return err
}
