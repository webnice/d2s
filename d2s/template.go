package d2s // import "gopkg.in/webnice/d2s.v1/d2s"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"text/template"
	"time"
)

var tpl = template.Must( // nolint: gochecknoglobals
	template.New("").
		Parse(`// Code generated by running "go generate"; DO NOT EDIT.
// This file was generated at {{ .Timestamp.UTC.Format "02.01.2006 15:04:05 UTC" }}
// The structure is based on the database table structure.
// Database: "{{ .Database }}"
// Table: "{{ .Table }}"

package {{ .Package }}

{{- if or .IsImportNil .IsImportTime }}

import (
{{- if .IsImportTime }}
	"time"
{{- end }}
{{- if .IsImportNil }}
{{- if .IsImportTime }}
{{ end }}
	nul "gopkg.in/webnice/lin.v1/nl"{{ end }}
){{ end -}}
{{- $cf := .ColumnsFormat }}

// {{ .Struct }} {{ .Comment }}
type {{ .Struct }} struct {
{{- range .Columns }}
	{{ printf $cf .Name .Type .Tags .Comment }}
{{- end }}
}
`))

type tplVar struct {
	Timestamp        time.Time
	Database         string
	Table            string
	Package          string
	Struct           string
	Comment          string
	Columns          []*tplVarCol
	ColumnsFormat    string
	ColumnNameLength int
	ColumnTypeLength int
	ColumnTagsLength int
	IsImportNil      bool
	IsImportTime     bool
}

type tplVarCol struct {
	Name    string
	Type    string
	Tags    string
	Comment string
}
