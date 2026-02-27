package uuid

import (
	"io"

	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
)

// V1 Создание UUID версии 1 основанной на текущем времени и MAC адресе.
func (ui *impl) V1() (ret UUID, err error) {
	var (
		timeNow      uint64
		clockSeq     uint16
		hardwareAddr []byte
		u            *uuid
	)

	u = new(uuid)
	if timeNow, clockSeq, err = ui.getClockSequence(); err != nil {
		ret = NULL
		return
	}
	binary.BigEndian.PutUint32(u.data[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u.data[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u.data[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u.data[8:], clockSeq)
	if hardwareAddr, err = ui.getHardwareAddr(); err != nil {
		return NULL, err
	}
	copy(u.data[10:], hardwareAddr)
	u.SetVersion(Version.V1)
	u.SetVariant(Variant.RFC4122)
	ret = u

	return
}

// V2 Создание UUID версии 2 основанной на POSIX UID/GID идентификаторах пользователя и группы соответственно.
func (ui *impl) V2(domain DomainType) (ret UUID, err error) {
	var u *uuid

	if ret, err = ui.V1(); err != nil {
		ret = NULL
		return
	}
	u = ret.(*uuid)
	switch domain {
	case domainPerson:
		binary.BigEndian.PutUint32(u.data[:], posixUID)
	case domainGroup:
		binary.BigEndian.PutUint32(u.data[:], posixGID)
	}
	u.data[9] = byte(domain)
	u.SetVersion(Version.V2)
	u.SetVariant(Variant.RFC4122)
	ret = u

	return
}

// V3 Создание UUID версии 3 основанной на MD5 алгоритме, пространстве имён и переданном названии.
func (ui *impl) V3(namespace NamespaceType, name string) (ret UUID) {
	var u *uuid

	u = ui.newFromHash(md5.New(), namespace, name)
	u.SetVersion(Version.V3)
	u.SetVariant(Variant.RFC4122)
	ret = u

	return
}

// V4 Создание UUID версии 4 основанной на генераторе случайных чисел.
func (ui *impl) V4() (ret UUID) {
	var (
		u   *uuid
		err error
	)

	u = &uuid{}
	if _, err = io.ReadFull(ui.rand, u.data[:]); err != nil {
		ret = NULL
		return
	}
	u.SetVersion(Version.V4)
	u.SetVariant(Variant.RFC4122)
	ret = u

	return
}

// V5 Создание UUID версии 5 основанной на SHA-1 хэшировании от пространства имён и названия.
func (ui *impl) V5(namespace NamespaceType, name string) (ret UUID) {
	var u *uuid

	u = ui.newFromHash(sha1.New(), namespace, name)
	u.SetVersion(Version.V5)
	u.SetVariant(Variant.RFC4122)
	ret = u

	return
}

// V6 Создание UUID версии 6.
func (ui *impl) V6() (ret UUID) {
	var (
		err         error
		u           *uuid
		timeNow     uint64
		clockSeq    uint16
		timeHighMid uint64
		timeHigh    uint32
		timeMid     uint16
		timeLow     uint16
		buf         []byte
	)

	u = new(uuid)
	if timeNow, clockSeq, err = ui.getClockSequence(); err != nil {
		ret = NULL
		return
	}
	timeHighMid = timeNow >> 12
	timeHigh = uint32(timeHighMid >> 16)
	timeMid = uint16(timeHighMid & 0xffff)
	timeLow = uint16(timeNow & 0xfff)
	timeLow |= 0x6000
	binary.BigEndian.PutUint32(u.data[0:], timeHigh)
	binary.BigEndian.PutUint16(u.data[4:], timeMid)
	binary.BigEndian.PutUint16(u.data[6:], timeLow)
	binary.BigEndian.PutUint16(u.data[8:], clockSeq|0x8000) // concat UUID variant.
	buf = make([]byte, 6)
	if _, err = io.ReadFull(ui.rand, buf); err != nil {
		ret = NULL
		return
	}
	copy(u.data[10:], buf[:])
	u.SetVersion(Version.V6)
	u.SetVariant(Variant.RFC4122)
	ret = u

	return
}
