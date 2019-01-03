package utils

import (
	"testing"
	"fmt"
)

func TestRandomString(t *testing.T) {
	fmt.Println(RandomString(), RandomString())
	if RandomString() == RandomString() {
		t.Error("Func RandowString not work")
	}
}
