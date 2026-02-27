package uuid

import (
	"io"
	"net"
	"sync"
	"time"
)

// Interface Интерфейс пакета.
type Interface interface {
	// V1 Создание UUID версии 1 основанной на текущем времени и MAC адресе.
	V1() (ret UUID, err error)

	// V2 Создание UUID версии 2 основанной на POSIX UID/GID идентификаторах пользователя и группы соответственно.
	V2(domain DomainType) (ret UUID, err error)

	// V3 Создание UUIDold версии 3.
	V3(namespace NamespaceType, name string) (ret UUID)

	// V4 Создание UUID версии 4 основанной на генераторе случайных чисел.
	V4() (ret UUID)

	// V5 Создание UUID версии 5 основанной на SHA-1 хешировании от пространства имён и названия.
	V5(namespace NamespaceType, name string) (ret UUID)

	// V6 Создание UUID версии 6.
	V6() (ret UUID)

	// FromBytes Конвертация среза байт в UUID, в случае ошибки, возвращается константа NULL.
	FromBytes(input []byte) (ret UUID)

	// FromString Конвертация строки в UUID, в случае ошибки, возвращается константа NULL.
	FromString(input string) (ret UUID)
}

// UUID Интерфейс объекта UUID.
type UUID interface {
	// Bytes Возвращает срез байт.
	Bytes() []byte

	// Version Возвращает версию алгоритма используемого при генерации UUID.
	Version() (ret VersionType)

	// Variant Возвращает версию варианта макета.
	Variant() (ret VariantType)

	// String Возвращает каноническое строковое представление UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
	String() (ret string)

	// Equal Сравнение UUID. Возвращает истину, если переданный UUID эквивалентен UUID объекта. Иначе возвращается ложь.
	Equal(uu UUID) bool

	// SetVersion Устанавливает бит версии UUID.
	SetVersion(v VersionType)

	// SetVariant Устанавливает бит варианта макета UUID.
	SetVariant(v VariantType)

	// СТАНДАРТНЫЕ ИНТЕРФЕЙСЫ

	// MarshalBinary Реализация интерфейса encoding.BinaryMarshaler.
	MarshalBinary() (data []byte, err error)

	// UnmarshalBinary Реализация интерфейса encoding.BinaryUnmarshaler.
	UnmarshalBinary(data []byte) (err error)

	// MarshalText Реализация интерфейса encoding.TextMarshaler.
	MarshalText() (text []byte, err error)

	// UnmarshalText Реализация интерфейса encoding.TextUnmarshaler.
	UnmarshalText(text []byte) (err error)

	// UnmarshalJSON Реализация интерфейса json.Unmarshaler.
	UnmarshalJSON(b []byte) (err error)

	// MarshalJSON Реализация интерфейса json.Marshaler.
	MarshalJSON() (ret []byte, err error)
}

// Объект сущности пакета.
type impl struct {
	clockSequenceOnce sync.Once
	hardwareAddrOnce  sync.Once
	storageMutex      sync.Mutex
	rand              io.Reader
	epochFunc         epochFunc
	hwAddrFunc        hwAddrFunc
	lastTime          uint64
	clockSequence     uint16
	hardwareAddr      [6]byte
}

type uuid struct {
	data [size]byte // Соответствует спецификации, описанной в RFC 4122.
}

type (
	version struct {
		Unknown VersionType // Неизвестная.
		V1      VersionType // V1.
		V2      VersionType // V2.
		V3      VersionType // V3.
		V4      VersionType // V4.
		V5      VersionType // V5.
		V6      VersionType // V6.
	}
	variant struct {
		NCS       VariantType // NCS.
		RFC4122   VariantType // RFC4122.
		Microsoft VariantType // Microsoft.
		Future    VariantType // Future.
	}
	namespace struct {
		DNS  NamespaceType // DNS Пространства имён UUIDold - DNS.
		URL  NamespaceType // DNS Пространства имён UUIDold - URL.
		OID  NamespaceType // DNS Пространства имён UUIDold - OID.
		X500 NamespaceType // DNS Пространства имён UUIDold - X500.
	}
	domain struct {
		Person DomainType // Person персональный.
		Group  DomainType // Group группа.
		Org    DomainType // Org Организация.
	}
)

// VersionType Версия UUIDold.
type VersionType struct {
	value byte
	name  string
}

// VariantType Тип данных для вариантов макета.
type VariantType byte

// NamespaceType Тип данных для пространства имён UUIDold.
type NamespaceType [size]byte

// DomainType Тип данных для DCE доменов.
type DomainType byte

type epochFunc func() time.Time

type hwAddrFunc func() (net.HardwareAddr, error)
