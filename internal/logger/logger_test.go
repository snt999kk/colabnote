package logger

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	Log(fmt.Errorf("MyError"))
}
