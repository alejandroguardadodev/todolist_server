package validation

type CustomTypeBehavior interface {
	CheckValue() bool
	Empty() bool
	String() string
}
