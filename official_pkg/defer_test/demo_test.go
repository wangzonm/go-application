package defer_test

import (
	"fmt"
	"testing"
)

func TestDemo(t *testing.T) {
}

func demo() {
	defer func() {
		fmt.Println("defer 1")
	}()
	func() {
		defer func() {
		}()
	}()
}
