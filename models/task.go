package models

import (
	"time"

	"todolistserver.com/test/models/customtypes"
	"todolistserver.com/test/validation"
)

type Task struct {
	ID          uint                              `json:"id" gorm:"primaryKey;autoIncrement"`
	Task        string                            `json:"task" gorm:"not null;type:varchar(100)" validate:"required"`
	Description string                            `json:"description" gorm:"type:text"`
	Status      customtypes.CustomeTaskStatusType `json:"status" gorm:"not null;type:task_status_type" validate:"customRequiredEnum,customCheckEnumValue"`
	ProjectID   int                               `json:"project_id"`
	Project     Project                           `gorm:"foreignKey:ProjectID"`
	Due         string                            `json:"due_date" gorm:"type:Date" validate:"required,customValidDate,customValiDateAfterOrEqualThanToday"`
	DueTime     string                            `json:"due_time"`
	Starred     bool                              `json:"starred" gorm:"default:false"`
	CreatedAt   time.Time
	UpdateAt    time.Time
}

func (t Task) GetDictionary() *Dictionary {
	dic := Dictionary{
		"id":          t.ID,
		"task":        t.Task,
		"description": t.Description,
		"status":      t.Status.String(),
		"project":     t.Project.GetDictionary(),
		"due":         t.Due,
		"due_time":    t.DueTime,
		"starred":     t.Starred,
	}

	return &dic
}

func (p Task) Validate() (*[]validation.ErrField, error) {

	if err := validation.Validate.StructExcept(p, "Project"); err != nil {
		errFields := validation.GetValidateInformation(err, p)

		return errFields, err
	}

	return nil, nil
}
