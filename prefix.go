package prefix

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
)

/// Public API

const DefaultFormat = `{{.LineNumber3}} `

type LinePrefixer interface {
	PrefixLine(string) string
}

/// Main implementation

type linePrefixer struct {
	Format     string
	LineNumber int

	t *template.Template
}

func New(format string) LinePrefixer {
	funcMap := template.FuncMap{}
	for k, v := range sprig.FuncMap() {
		funcMap[k] = v
	}

	return &linePrefixer{
		Format:     DefaultFormat,
		LineNumber: 0,

		t: template.Must(template.New("").Funcs(funcMap).Parse(format)),
	}
}

func (p *linePrefixer) PrefixLine(line string) string {
	p.LineNumber++
	var prefix bytes.Buffer
	_ = p.t.Execute(&prefix, p)
	return prefix.String() + line
}

/// Template helpers

func (p *linePrefixer) LineNumber3() string { return fmt.Sprintf("%-3d", p.LineNumber) }
func (p *linePrefixer) LineNumber4() string { return fmt.Sprintf("%-4d", p.LineNumber) }
func (p *linePrefixer) LineNumber5() string { return fmt.Sprintf("%-5d", p.LineNumber) }
