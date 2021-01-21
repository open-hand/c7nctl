package utils

import (
	"fmt"
	"testing"
)

func TestRandomString(t *testing.T) {
	str1 := RandomString(10)
	str2 := RandomString(10)
	fmt.Println(str1, str2)
	if str1 == str2 {
		t.Error("Func RandowString not work")
	}
}
