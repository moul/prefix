package prefix

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	sprig "github.com/Masterminds/sprig/v3"
	"moul.io/u"
)

/// Public API

// AvailablePatterns is the list of available patterns to be used by the user.
//
// This variable is only used to generate usage.
var AvailablePatterns = map[string]string{
	`{{.LineNumber}}`:                `display line number`,
	`{{.LineNumber3}}`:               `alias for {{printf "%-3d" .LineNumber}}`,
	`{{.LineNumber4}}`:               `alias for {{printf "%-4d" .LineNumber}}`,
	`{{.LineNumber5}}`:               `alias for {{printf "%-5d" .LineNumber}}`,
	`{{.Format}}`:                    `the value you set with -format`,
	`{{.Uptime}}`:                    `time since the the prefixer was initialized`,
	`{{.Duration}}`:                  `time since previous line was started`,
	`{{.Uptime | short_duration}}`:   `{{.Uptime}} displayed in a pretty & short format (len<=7)`,
	`{{.Duration | short_duration}}`: `{{.Duration}} displayed in a pretty & short format (len<=7)`,
	`{{now}}`:                        `current date (format: 2006-01-02 15:04:05.999999999 -0700 MST)`,
	`{{now | unixEpoch}}`:            `current timestamp`,
	`{{uuidv4}}`:                     `UUID of the v4 (randomly generated) type`,
	`{{env "USER"}}`:                 `replace with content of the $USER env var`,
	`{{.ShortDuration}}`:             `alias for {{.Duration | short_duration}}`,
	`{{.ShortUptime}}`:               `alias for {{.Uptime | short_duration}}`,
}

// AvailablePresets is the list of available presets.
//
// Those presets are automatically replaced during the initialization of the prefixer.
var AvailablePresets = map[string]string{
	"{{DEFAULT}}":    `{{.LineNumber3}} up={{.ShortUptime}} d={{.ShortDuration}} |`,
	"{{SLOW_LINES}}": `{{if (gt .Duration 1000000000)}}SLOW{{else}}    {{end}} {{.Duration | short_duration}} `,
	"{{SHORT_DATE}}": `{{now | date "06/02/01 15:04:05"}}`,
}

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

func (p *linePrefixer) String() string {
	return fmt.Sprintf("LinePrefixer{%q}", p.Format)
}

// New returns an initialized LinePrefixer.
func New(format string) LinePrefixer {
	if format == "" {
		format = AvailablePresets["{{DEFAULT}} "]
	}

	// apply presets
	for {
		before := format
		for k, v := range AvailablePresets {
			format = strings.ReplaceAll(format, k, v)
		}
		if before == format {
			break
		}
	}

	// build funcmap
	funcMap := template.FuncMap{}
	for k, v := range sprig.FuncMap() {
		funcMap[k] = v
	}
	funcMap["short_duration"] = shortDuration

	return &linePrefixer{
		Format:     format,
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

func (p *linePrefixer) LineNumber3() string     { return fmt.Sprintf("%-3d", p.LineNumber) }
func (p *linePrefixer) LineNumber4() string     { return fmt.Sprintf("%-4d", p.LineNumber) }
func (p *linePrefixer) LineNumber5() string     { return fmt.Sprintf("%-5d", p.LineNumber) }
func (p *linePrefixer) Uptime() time.Duration   { return time.Since(p.startedAt) }
func (p *linePrefixer) Duration() time.Duration { return time.Since(p.lastTime) }
func (p *linePrefixer) ShortDuration() string   { return shortDuration(p.Duration()) }
func (p *linePrefixer) ShortUptime() string     { return shortDuration(p.Uptime()) }

func shortDuration(d time.Duration) string { return fmt.Sprintf("%-7s", u.ShortDuration(d)) }
