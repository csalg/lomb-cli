package assert

import "fmt"

func NotNil(t interface{}, msg string) {
	if t == nil {
		panic(fmt.Sprintf("Expected not nil: %s", msg))
	}
}
