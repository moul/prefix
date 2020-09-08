package prefix

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	sprig "github.com/Masterminds/sprig/v3"
	"moul.io/u"
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

	t         *template.Template
	startedAt time.Time // used by {{.Uptime}}
	lastTime  time.Time // used by {{.Duration}}
}

func New(format string) LinePrefixer {
	funcMap := template.FuncMap{}
	for k, v := range sprig.FuncMap() {
		funcMap[k] = v
	}

	return &linePrefixer{
		Format:     DefaultFormat,
		LineNumber: 0,

		t:         template.Must(template.New("").Funcs(funcMap).Parse(format)),
		startedAt: time.Now(),
		lastTime:  time.Now(),
	}
}

func (p *linePrefixer) PrefixLine(line string) string {
	p.LineNumber++
	var prefix bytes.Buffer
	_ = p.t.Execute(&prefix, p)
	p.lastTime = time.Now()
	return prefix.String() + line
}

/// Template helpers

func (p *linePrefixer) LineNumber3() string { return fmt.Sprintf("%-3d", p.LineNumber) }
func (p *linePrefixer) LineNumber4() string { return fmt.Sprintf("%-4d", p.LineNumber) }
func (p *linePrefixer) LineNumber5() string { return fmt.Sprintf("%-5d", p.LineNumber) }
func (p *linePrefixer) Uptime() string {
	return fmt.Sprintf("%-7s", u.ShortDuration(time.Since(p.startedAt)))
}
func (p *linePrefixer) Duration() string {
	return fmt.Sprintf("%-7s", u.ShortDuration(time.Since(p.lastTime)))
}
