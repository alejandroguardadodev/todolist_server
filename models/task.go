package models

import (
	"time"

	"todolistserver.com/test/models/customtypes"
	"todolistserver.com/test/validation"
)

type Task struct {
	ID          uint                              `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string                            `json:"title" gorm:"not null;type:varchar(100)" validate:"required"`
	Description string                            `json:"description" gorm:"type:text"`
	Status      customtypes.CustomeTaskStatusType `json:"status" gorm:"not null;type:task_status_type" validate:"customRequiredEnum,customCheckEnumValue"`
	ProjectID   uint                              `json:"project_id" gorm:"not null" validate:"required"`
	Project     Project                           `gorm:"foreignKey:ProjectID"`
	DueDate     string                            `json:"due_date" gorm:"type:text" validate:"customValidDate,customValiDateAfterOrEqualThanToday"`
	DueTime     string                            `json:"due_time"`
	Starred     bool                              `json:"starred" gorm:"default:false"`
	CreatedAt   time.Time
	UpdateAt    time.Time
}

func (t Task) GetDictionary(isDefaultProject bool) *Dictionary {
	dic := Dictionary{
		"id":          t.ID,
		"title":       t.Title,
		"description": t.Description,
		"status":      t.Status.String(),
		"project":     t.Project.GetDictionary(),
		"due_date":    t.DueDate,
		"due_time":    t.DueTime,
		"starred":     t.Starred,
	}

	if isDefaultProject {
		dic["project"] = ""
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
