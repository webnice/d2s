package uuid

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"time"
)

func init() { singleton = newUUID() }

// Get Получение объекта сущности пакета, возвращается интерфейс пакета.
//
//goland:noinspection GoUnusedExportedFunction
func Get() Interface { return singleton }

func newUUID() *impl {
	return &impl{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       rand.Reader,
	}
}

// Bytes Возвращает срез байт.
func (u *uuid) Bytes() []byte { return u.data[:] }

// Version Возвращает версию алгоритма используемого при генерации UUID.
func (u *uuid) Version() (ret VersionType) {
	switch u.data[6] >> 4 {
	case Version.V1.value:
		ret = Version.V1
	case Version.V2.value:
		ret = Version.V2
	case Version.V3.value:
		ret = Version.V3
	case Version.V4.value:
		ret = Version.V4
	case Version.V5.value:
		ret = Version.V5
	default:
		ret = Version.Unknown
	}

	return
}

// Variant Возвращает версию варианта макета.
func (u *uuid) Variant() (ret VariantType) {
	switch {
	case (u.data[8] >> 7) == 0x00:
		ret = Variant.NCS
	case (u.data[8] >> 6) == 0x02:
		ret = Variant.RFC4122
	case (u.data[8] >> 5) == 0x06:
		ret = Variant.Microsoft
	case (u.data[8] >> 5) == 0x07:
		fallthrough
	default:
		ret = Variant.Future
	}

	return
}

// String Возвращает каноническое строковое представление UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u *uuid) String() (ret string) {
	var (
		buf       []byte
		separator = byte('-')
	)

	buf = make([]byte, 36)
	hex.Encode(buf[0:8], u.data[0:4])
	buf[8] = separator
	hex.Encode(buf[9:13], u.data[4:6])
	buf[13] = separator
	hex.Encode(buf[14:18], u.data[6:8])
	buf[18] = separator
	hex.Encode(buf[19:23], u.data[8:10])
	buf[23] = separator
	hex.Encode(buf[24:], u.data[10:])
	ret = string(buf)

	return
}

// Equal Сравнение UUID. Возвращает истину, если переданный UUID эквивалентен UUID объекта. Иначе возвращается ложь.
func (u *uuid) Equal(uu UUID) bool { return bytes.Equal(u.data[:], uu.(*uuid).data[:]) }

// SetVersion Устанавливает бит версии UUID.
func (u *uuid) SetVersion(v VersionType) { u.data[6] = (u.data[6] & 0x0f) | (v.value << 4) }

// SetVariant Устанавливает бит варианта макета UUID.
func (u *uuid) SetVariant(v VariantType) {
	switch v {
	case Variant.NCS:
		u.data[8] = u.data[8]&(0xff>>1) | (0x00 << 7)
	case Variant.RFC4122:
		u.data[8] = u.data[8]&(0xff>>2) | (0x02 << 6)
	case Variant.Microsoft:
		u.data[8] = u.data[8]&(0xff>>3) | (0x06 << 5)
	case Variant.Future:
		fallthrough
	default:
		u.data[8] = u.data[8]&(0xff>>3) | (0x07 << 5)
	}
}
