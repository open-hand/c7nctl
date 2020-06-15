package utils

import (
	"github.com/choerodon/c7nctl/pkg/context"
	"testing"
)

//func TestAskAgreeTerms(t *testing.T){
//	AskAgreeTerms()
//}

func TestCheckMatch(t *testing.T) {
	input := context.Input{
		Regex:   "^\\S{8}$",
		Include: []context.KV{context.KV{Name: "必须包含数字", Value: "\\d+"}, context.KV{Name: "必须包含大写", Value: "[A-Z]+"}},
	}

	if CheckMatch("abcdef", input) {
		t.Error("regex check error")
	}

	if !CheckMatch("ce123Qwe", input) {
		t.Error("regex check error")
	}

	if CheckMatch("ce123456", input) {
		t.Error("regex check error")
	}

	if CheckMatch("cdsfaQes", input) {
		t.Error("regex check error")
	}
}
