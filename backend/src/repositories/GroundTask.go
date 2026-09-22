package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type TaskRepository struct{ DB *gorm.DB }

func NewTaskRepository(db *gorm.DB) *TaskRepository { return &TaskRepository{DB: db} }

func (r *TaskRepository) CreateBatch(tasks []models.GroundTask) error {
	if len(tasks) == 0 {
		return nil
	}
	return r.DB.Create(&tasks).Error
}

func (r *TaskRepository) ListByTurnaround(turnaroundID uint) ([]models.GroundTask, error) {
	var rows []models.GroundTask
	err := r.DB.Where("turnaround_id = ?", turnaroundID).Order("planned_start ASC").Find(&rows).Error
	return rows, err
}

func (r *TaskRepository) List(page, size int, status string) ([]models.GroundTask, int64, error) {
	var rows []models.GroundTask
	var total int64
	q := r.DB.Model(&models.GroundTask{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&total)
	err := q.Order("planned_start ASC").Offset((page - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

func (r *TaskRepository) Get(id uint) (*models.GroundTask, error) {
	var t models.GroundTask
	if err := r.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) Save(t *models.GroundTask) error { return r.DB.Save(t).Error }

func (r *TaskRepository) CountByTurnaround(turnaroundID uint) (int64, error) {
	var n int64
	err := r.DB.Model(&models.GroundTask{}).Where("turnaround_id = ?", turnaroundID).Count(&n).Error
	return n, err
}
