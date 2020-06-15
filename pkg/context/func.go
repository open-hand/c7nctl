package context

import (
	"bytes"
	"github.com/vinkdong/gox/log"
	"text/template"
)

type Func struct {
	Name     string
	Template string
	Text     string
}

func (fn *Func) Execute(obj interface{}, opts ...string) string {
	if fn.Text != "" {
		return fn.Text
	}
	if fn.Template != "" {
		tpl, err := template.New("ntpl").Parse(fn.Template)
		var data bytes.Buffer
		if err != nil {
			log.Error(err)
			return ""
		}
		if err := tpl.Execute(&data, obj); err != nil {
			log.Error(err)
			return ""
		}
		return data.String()
	}
	return ""
}
