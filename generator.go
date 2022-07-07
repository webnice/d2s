// Package d2s
package d2s

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	d2sTypes "github.com/webnice/d2s/types"
)

// Generator Создание контента golang файла на основе загруженных данных таблицы
func (d2s *impl) Generator(inf *d2sTypes.TableInfo) (ret *bytes.Buffer, err error) {
	const (
		libTime       = "time.Time"
		libNil        = "nul."
		keyNewLine    = "\n"
		keyReturn     = "\r"
		columnsFormat = "%%-%ds %%-%ds `%%-%ds` // %%s"
	)
	var (
		vars             *tplVar
		col              *tplVarCol
		i, maxLengthName int
	)

	ret = &bytes.Buffer{}
	vars = &tplVar{
		Timestamp: time.Now(),
		Database:  inf.Database,
		Table:     inf.Table,
		Package:   inf.Package,
		Struct:    inf.Struct,
		Comment:   inf.Comment,
		Columns:   make([]*tplVarCol, 0, len(inf.Columns)),
	}
	for i = range inf.Columns {
		if maxLengthName < len(inf.Columns[i].Name) {
			maxLengthName = len(inf.Columns[i].Name)
		}
	}
	for i = range inf.Columns {
		col = &tplVarCol{
			Name:    inf.Columns[i].Name.CamelCase(),
			Type:    inf.Columns[i].TypeGolang.String(),
			Tags:    d2s.GeneratorTags(inf.Columns[i], maxLengthName),
			Comment: strings.Replace(inf.Columns[i].Comment, keyReturn, "", -1),
		}
		col.Comment = strings.Replace(col.Comment, keyNewLine, "\\n", -1)
		if vars.ColumnNameLength < len(inf.Columns[i].Name.CamelCase()) {
			vars.ColumnNameLength = len(inf.Columns[i].Name.CamelCase())
		}
		if vars.ColumnTypeLength < len(inf.Columns[i].TypeGolang.String()) {
			vars.ColumnTypeLength = len(inf.Columns[i].TypeGolang.String())
		}
		if vars.ColumnTagsLength < len(col.Tags) {
			vars.ColumnTagsLength = len(col.Tags)
		}
		if strings.EqualFold(col.Type, libTime) {
			vars.IsImportTime = true
		}
		if strings.Contains(col.Type, libNil) {
			vars.IsImportNil = true
		}
		vars.Columns = append(vars.Columns, col)
	}
	vars.ColumnsFormat = fmt.Sprintf(columnsFormat, vars.ColumnNameLength, vars.ColumnTypeLength, vars.ColumnTagsLength)
	err = tpl.Execute(ret, vars)

	return
}

// GeneratorTags Создание тэгов описания колонки
func (d2s *impl) GeneratorTags(col *d2sTypes.ColumnInfo, max int) (ret string) {
	const (
		stdDbTag      = `db:"%s"`
		grmDbtag      = `gorm:"column:%s;%%-%ds`
		grmPrimaryKey = `primary_key;`
		semicolon     = `;`
	)
	var (
		buf          *bytes.Buffer
		std, grm     string
		stdLengthMax int
	)

	buf = &bytes.Buffer{}
	// Теги для database/sql и github.com/jmoiron/sqlx
	stdLengthMax = len(stdDbTag) - 2 + max
	std = fmt.Sprintf(stdDbTag, col.Name)
	_, _ = fmt.Fprintf(buf, fmt.Sprintf(`%%-%ds`, stdLengthMax), std)
	// Теги для github.com/jinzhu/gorm
	grm = fmt.Sprintf(grmDbtag, col.Name, max-len(col.Name)+len(grmPrimaryKey))
	if col.IsPrimary {
		grm = fmt.Sprintf(grm, grmPrimaryKey)
	} else {
		grm = fmt.Sprintf(grm, "")
	}
	_, _ = fmt.Fprintf(buf, " %s", grm)
	// AUTO_INCREMENT
	if col.IsAutoIncrement {
		_, _ = fmt.Fprintf(buf, "AUTO_INCREMENT;")
	}
	// NOT NULL
	if col.IsNullable {
		_, _ = fmt.Fprintf(buf, "NULL;")
	} else {
		_, _ = fmt.Fprintf(buf, "NOT NULL;")
	}
	// DEFAULT
	if col.Default == nil {
		_, _ = fmt.Fprintf(buf, "DEFAULT NULL;")
	} else {
		_, _ = fmt.Fprintf(buf, "DEFAULT '%s';", *col.Default)
	}
	// SIZE
	if col.TypeGolang.Simple == d2sTypes.TString && col.Size != nil && *col.Size > 0 {
		_, _ = fmt.Fprintf(buf, "size:%d;", *col.Size)
	}
	// PRECISION and SCALE
	switch col.TypeGolang.Simple {
	case d2sTypes.TFloat64, d2sTypes.TInt64, d2sTypes.TUint64:
		if col.Precision != nil && *col.Precision > 0 {
			_, _ = fmt.Fprintf(buf, "precision:%d;", *col.Precision)
		}
		if col.Scale != nil && *col.Scale > 0 {
			_, _ = fmt.Fprintf(buf, "scale:%d;", *col.Scale)
		}
	}
	// TYPE
	_, _ = fmt.Fprintf(buf, "type:%s", col.TypeFull)
	// UNSIGNED append for type
	if col.IsUnsigned {
		_, _ = fmt.Fprintf(buf, " UNSIGNED;")
	} else {
		_, _ = fmt.Fprintf(buf, ";")
	}
	// Удаление завершающей точки с запятой
	buf = bytes.NewBuffer(bytes.TrimRight(buf.Bytes(), semicolon))
	_, _ = fmt.Fprintf(buf, `"`)
	ret = buf.String()

	return
}
