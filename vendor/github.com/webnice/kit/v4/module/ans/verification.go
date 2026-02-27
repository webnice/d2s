package ans

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/webnice/dic"
	kitModuleTrace "github.com/webnice/kit/v4/module/trace"

	"github.com/go-playground/validator/v10"
)

// Verify Верификация данных с использованием библиотеки github.com/go-playground/validator.
// Результатом будет nil - если ошибок нет, либо заполненный ошибками объект интерфейса RestErrorInterface.
func (ans *impl) Verify(variable any) (ret RestErrorInterface) {
	const errPanicCatchingException = "верификация данных прервана паникой:\n%s\n%s"
	var (
		err         error
		rt          reflect.Type
		rv          reflect.Value
		item        any
		field       []RestErrorField
		fieldErrors []RestErrorField
		n           int
	)

	// При вызове reflect возможна паника.
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(errPanicCatchingException, e, kitModuleTrace.StackShort())
			if ans.logErrorf(err.Error()); ret == nil {
				ret = ans.NewRestError(dic.Status().InternalServerError, err)
			}
		}
	}()
	// Разные способы верификации в зависимости от типа объекта.
	switch rt = indirectType(reflect.TypeOf(variable)); rt.Kind() {
	case reflect.Slice:
		rv = indirectValue(reflect.ValueOf(variable))
		for n = 0; n < rv.Len(); n++ {
			item = rv.Index(n).Interface()
			if fieldErrors, err = ans.verifyValidatorV10(item); len(fieldErrors) > 0 {
				field = append(field, fieldErrors...)
			}
			if err != nil {
				break
			}
		}
	default:
		field, err = ans.verifyValidatorV10(variable)
	}
	// Если ошибок не найдено.
	if len(field) <= 0 {
		return
	}
	// Если ошибки есть, сбор всех ошибок в объект интерфейса RestErrorInterface.
	ret = ans.NewRestError(dic.Status().BadRequest, err)
	for n = range field {
		ret.AddWithKey(
			field[n].Field,
			field[n].FieldValue,
			field[n].Message,
			field[n].I18nKey,
		)
	}

	return
}

// Проверка объекта через внешний верификатор. Описания проверок определяются тегами структуры данных объекта.
func (ans *impl) verifyValidatorV10(variable any) (ret []RestErrorField, err error) {
	var (
		verificatorError validator.ValidationErrors
		verificator      *validator.Validate
		ok               bool
		n                int
	)

	verificator = validator.New(validator.WithRequiredStructEnabled())
	if err = verificator.Struct(variable); err == nil {
		return
	}
	if ok = errors.As(err, &verificatorError); !ok {
		return
	}
	for n = range verificatorError {
		ret = append(ret, RestErrorField{
			Field:      verificatorError[n].StructField(),
			FieldValue: fmt.Sprintf("%v", verificatorError[n].Value()),
			Message:    verificatorError[n].ActualTag(),
		})
	}

	return
}
