package more_goconvey_assertions

import "fmt"

const (
	success                = ""
	needExactValues        = "This assertion requires exactly %d comparison values (you provided %d)."
	needNonEmptyCollection = "This assertion requires at least 1 comparison value (you provided 0)."
)

const (
	shouldBeString = "The argument to this assertion must be a string (you provided %v)."

	//TODO: These do not exist in goconvey yet
	shouldBeBool        = "The argument to this assertion must be a bool (you provided %v)."
	shouldHaveBeenTrue  = "Expected: true\nActual:   %v"
	shouldHaveBeenFalse = "Expected: false\nActual:   %v"
)

func need(needed int, expected []interface{}) string {
	if len(expected) != needed {
		return fmt.Sprintf(needExactValues, needed, len(expected))
	}
	return success
}
