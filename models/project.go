package models

import (
	"time"

	"todolistserver.com/test/validation"
)

type Project struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string `json:"title" gorm:"not null;type:varchar(300)" validate:"required"`
	User      string `json:"user" gorm:"not null;type:varchar(100)" validate:"required"`
	CreatedAt time.Time
	UpdateAt  time.Time
}

type CompressedProject struct {
	ID    uint
	Title string
}

func (p Project) GetCompressedProjectVersion() *CompressedProject {
	compress := CompressedProject{
		ID:    p.ID,
		Title: p.Title,
	}

	return &compress
}

func (p Project) GetDictionary() *Dictionary {
	dic := Dictionary{
		"id":    p.ID,
		"title": p.Title,
	}

	return &dic
}

func (p Project) Validate() (*[]validation.ErrField, error) {

	if err := validation.Validate.StructExcept(p, "Project"); err != nil {
		errFields := validation.GetValidateInformation(err, p)

		return errFields, err
	}

	return nil, nil
}
