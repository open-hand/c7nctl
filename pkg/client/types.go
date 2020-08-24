package client

import (
	"github.com/choerodon/c7nctl/pkg/utils"
)

type ChartValue struct {
	Name  string
	Value string
	Input utils.Input
	Case  string
	Check string
}

type PreValue struct {
	Name  string
	Value string
	Check string
	Input utils.Input
}

type PreValueList []*PreValue

type BackendTask struct {
	Name    string
	Success bool
}
