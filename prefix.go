package prefix

import (
	"bytes"

	"text/template"
)

const DefaultFormat = `{{printf "%-3d" .LineNumber}} `

type LinePrefixer interface {
	PrefixLine(string) string
}

type linePrefixer struct {
	Format     string
	LineNumber int

	t *template.Template
}

func New(format string) LinePrefixer {
	return &linePrefixer{
		Format:     DefaultFormat,
		LineNumber: 0,

		t: template.Must(template.New("").Parse(format)),
	}
}

func (p *linePrefixer) PrefixLine(line string) string {
	p.LineNumber++
	var prefix bytes.Buffer
	_ = p.t.Execute(&prefix, p)
	return prefix.String() + line
}
