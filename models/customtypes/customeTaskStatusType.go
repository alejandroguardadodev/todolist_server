package customtypes

import (
	"database/sql/driver"
	"strings"
)

type CustomeTaskStatusType string

const (
	CUSTOM_TASK_STATUS_TYPE_TODO  CustomeTaskStatusType = "TODO"
	CUSTOM_TASK_STATUS_TYPE_DOING CustomeTaskStatusType = "DOING"
	CUSTOM_TASK_STATUS_TYPE_DONE  CustomeTaskStatusType = "DONE"
)

func (cm *CustomeTaskStatusType) Scan(value interface{}) error {
	if val, ok := value.(string); ok {
		*cm = CustomeTaskStatusType(val)
	} else {
		*cm = CustomeTaskStatusType(value.([]byte))
	}

	return nil
}

func (cm CustomeTaskStatusType) Value() (driver.Value, error) {
	return string(cm), nil
}

func (cmt CustomeTaskStatusType) String() string {
	return string(cmt)
}

func (cm CustomeTaskStatusType) Empty() bool {
	value, _ := cm.Value()

	return value.(string) == ""
}

func (cm CustomeTaskStatusType) CheckValue() bool {
	str := strings.ToUpper(string(cm))

	return str == string(CUSTOM_TASK_STATUS_TYPE_TODO) || str == string(CUSTOM_TASK_STATUS_TYPE_DOING) || str == string(CUSTOM_TASK_STATUS_TYPE_DONE)
}
