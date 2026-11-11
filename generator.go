package d2s

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	d2sTypes "github.com/webnice/d2s/types"
)

// Generator Создание контента golang файла на основе загруженных данных таблицы.
func (d2s *impl) Generator(inf *d2sTypes.TableInfo) (ret *bytes.Buffer, err error) {
	const (
		prefixTime    = "time."
		prefixNix     = "nix."
		typeJSON      = "json.RawMessage"
		columnsFormat = "%%-%ds %%-%ds `%%-%ds` // %%s"
	)
	var (
		vars             *tplVar
		col              *tplVarCol
		i, maxLengthName int
		importTime       bool
		importNix        bool
		importJSON       bool
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
		if col, err = d2s.generatorColumn(inf.Columns[i], maxLengthName); err != nil {
			return
		}
		if vars.ColumnNameLength < len(col.Name) {
			vars.ColumnNameLength = len(col.Name)
		}
		if vars.ColumnTypeLength < len(col.Type) {
			vars.ColumnTypeLength = len(col.Type)
		}
		if vars.ColumnTagsLength < len(col.Tags) {
			vars.ColumnTagsLength = len(col.Tags)
		}
		if strings.HasPrefix(col.Type, prefixTime) {
			importTime = true
		}
		if strings.Contains(col.Type, prefixNix) {
			importNix = true
		}
		if strings.Contains(col.Type, typeJSON) {
			importJSON = true
		}
		vars.Columns = append(vars.Columns, col)
	}
	vars.Imports = d2s.importBlock(importJSON, importTime, importNix)
	vars.ColumnsFormat = fmt.Sprintf(columnsFormat, vars.ColumnNameLength, vars.ColumnTypeLength, vars.ColumnTagsLength)
	err = tpl.Execute(ret, vars)

	return
}

// importBlock Формирование блока импортов сгенерированного golang файла.
func (d2s *impl) importBlock(importJSON, importTime, importNix bool) (ret string) {
	var std, ext []string

	if importJSON {
		std = append(std, `"encoding/json"`)
	}
	if importTime {
		std = append(std, `"time"`)
	}
	if importNix {
		ext = append(ext, `"github.com/webnice/nix/v2"`)
	}
	switch {
	case len(std)+len(ext) == 0:
		return
	case len(std)+len(ext) == 1:
		return "import " + append(std, ext...)[0]
	}
	ret = "import (\n"
	for _, p := range std {
		ret += "\t" + p + "\n"
	}
	if len(std) > 0 && len(ext) > 0 {
		ret += "\n"
	}
	for _, p := range ext {
		ret += "\t" + p + "\n"
	}
	ret += ")"

	return
}

// generatorColumn Формирование описания колонки с учётом переопределения типа.
// Приоритет переопределений: командная строка, комментарий колонки базы данных,
// базовый маппинг типов.
func (d2s *impl) generatorColumn(ci *d2sTypes.ColumnInfo, maxLengthName int) (ret *tplVarCol, err error) {
	const keyNewLine, keyReturn = "\n", "\r"

	var (
		canonicalCLI, canonicalComment, cleanedComment string
		hasCLI, hasComment                             bool
	)

	ret = &tplVarCol{
		Name: ci.Name.CamelCase(),
		Tags: d2s.GeneratorTags(ci, maxLengthName),
	}
	canonicalCLI, hasCLI = d2s.override[strings.ToLower(string(ci.Name))]
	canonicalComment, cleanedComment, hasComment = d2sTypes.OverrideComment(ci.Comment)
	switch {
	case hasCLI:
		ret.Type = d2sTypes.OverrideGoType(canonicalCLI, ci.IsNullable)
	case hasComment:
		ret.Type = d2sTypes.OverrideGoType(canonicalComment, ci.IsNullable)
	case ci.TypeGolang != nil && ci.TypeGolang.Simple == d2sTypes.TUnknown:
		// Неизвестный тип базы данных и нет переопределения типа поля.
		err = fmt.Errorf("не реализовано сопоставление типа %q, укажите тип поля через переопределение", ci.TypeSimple)
		return
	default:
		ret.Type = ci.TypeGolang.String()
	}
	// Маркер {...} удаляется из комментария независимо от победившего переопределения.
	ret.Comment = strings.Replace(cleanedComment, keyReturn, "", -1)
	ret.Comment = strings.Replace(ret.Comment, keyNewLine, "\\n", -1)

	return
}

// GeneratorTags Создание тэгов описания колонки.
func (d2s *impl) GeneratorTags(col *d2sTypes.ColumnInfo, max int) (ret string) {
	const (
		stdDbTag        = `db:"%s"`
		stdDbTagLength  = len(stdDbTag) - 2
		grmTagStart     = `gorm:"`
		grmTagEnd       = `"`
		grmColumn       = `column:%s`
		grmPrimaryKey   = `primaryKey`
		grmAutoInc      = `autoIncrement`
		grmNotNull      = `not null`
		grmDefault      = `default:%s`
		grmDefaultNull  = `default:null`
		grmSize         = `size:%d`
		grmPrecision    = `precision:%d`
		grmScale        = `scale:%d`
		grmType         = `type:%s`
		keyUnsigned     = ` UNSIGNED`
		keySemicolon    = `;`
		escapeSemicolon = `\;`
	)
	var (
		settings []string
		std      string
		stdMax   int
	)

	// Теги для database/sql и github.com/jmoiron/sqlx.
	stdMax = stdDbTagLength + max
	std = fmt.Sprintf(stdDbTag, col.Name)
	// Теги для gorm.io/gorm.
	settings = append(settings, fmt.Sprintf(grmColumn, col.Name))
	if col.IsPrimary {
		settings = append(settings, grmPrimaryKey)
	}
	if col.IsAutoIncrement {
		settings = append(settings, grmAutoInc)
	}
	if !col.IsNullable {
		settings = append(settings, grmNotNull)
	}
	switch {
	case col.Default != nil:
		settings = append(settings, fmt.Sprintf(grmDefault, strings.ReplaceAll(*col.Default, keySemicolon, escapeSemicolon)))
	case col.IsNullable:
		settings = append(settings, grmDefaultNull)
	}
	if (col.TypeGolang.Simple == d2sTypes.TString || col.TypeGolang.Simple == d2sTypes.TBytes) &&
		col.Size != nil && *col.Size > 0 {
		settings = append(settings, fmt.Sprintf(grmSize, *col.Size))
	}
	switch col.TypeGolang.Simple {
	case d2sTypes.TFloat32, d2sTypes.TFloat64,
		d2sTypes.TInt8, d2sTypes.TInt16, d2sTypes.TInt32, d2sTypes.TInt64,
		d2sTypes.TUint8, d2sTypes.TUint16, d2sTypes.TUint32, d2sTypes.TUint64:
		if col.Precision != nil && *col.Precision > 0 {
			settings = append(settings, fmt.Sprintf(grmPrecision, *col.Precision))
		}
		if col.Scale != nil && *col.Scale > 0 {
			settings = append(settings, fmt.Sprintf(grmScale, *col.Scale))
		}
	}
	if col.IsUnsigned {
		settings = append(settings, fmt.Sprintf(grmType, col.TypeFull+keyUnsigned))
	} else {
		settings = append(settings, fmt.Sprintf(grmType, col.TypeFull))
	}
	ret = fmt.Sprintf("%-*s %s%s%s", stdMax, std, grmTagStart, strings.Join(settings, keySemicolon), grmTagEnd)

	return
}
