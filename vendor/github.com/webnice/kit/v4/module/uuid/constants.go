package uuid

import "os"

// Размер UUID в байтах.
const size = 16

const (
	// Варианты макета UUID
	variantNCS       = VariantType(0) // NCS.
	variantRFC4122   = VariantType(1) // RFC4122.
	variantMicrosoft = VariantType(2) // Microsoft.
	variantFuture    = VariantType(3) // Future.

	// DCE домены для генерации UUID основанной на POSIX UID/GID.
	domainPerson = DomainType(0) // Person.
	domainGroup  = DomainType(1) // Group.
	domainOrg    = DomainType(2) // Org.

	// Начало эпохи (15 октября 1582) и unix эпоха (1 января 1970).
	epochStart = 122192928000000000
)

// Основной экземпляр объекта пакета uuid.
var singleton *impl

// Проверка наличия переменных на уровне компилятора и предотвращение предупреждения "переменные не используются".
var _, _, _, _ = Version, Variant, Namespace, Domain

var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}
	posixUID   = uint32(os.Getuid())
	posixGID   = uint32(os.Getgid())

	// NULL Специальный UUID состоящий из нулей, длинной 128 бит или 16 байт.
	NULL UUID = &uuid{data: [size]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}

	// Version Поддерживаемые версии UUID.
	Version = version{
		Unknown: VersionType{value: 0, name: `Unknown`}, // Неизвестная версия.
		V1:      VersionType{value: 1, name: `V1`},      // V1 версия 1.
		V2:      VersionType{value: 2, name: `V2`},      // V2 версия 2.
		V3:      VersionType{value: 3, name: `V3`},      // V3 версия 3.
		V4:      VersionType{value: 4, name: `V4`},      // V4 версия 4.
		V5:      VersionType{value: 5, name: `V5`},      // V5 версия 5.
		V6:      VersionType{value: 6, name: `V6`},      // V6 версия 6.
	}

	// Variant Варианты макета UUID.
	Variant = variant{
		NCS:       variantNCS,
		RFC4122:   variantRFC4122,
		Microsoft: variantMicrosoft,
		Future:    variantFuture,
	}

	// Namespace Пространства имён.
	Namespace = namespace{
		DNS:  NamespaceType(singleton.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8").(*uuid).data),
		URL:  NamespaceType(singleton.FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8").(*uuid).data),
		OID:  NamespaceType(singleton.FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8").(*uuid).data),
		X500: NamespaceType(singleton.FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8").(*uuid).data),
	}

	// Domain DCE домены.
	Domain = domain{
		Person: domainPerson,
		Group:  domainGroup,
		Org:    domainOrg,
	}
)
