package oidb_0x5eb

import "encoding/binary"

const (
	_locationBytesLen         = 12 // country+province+city == 4+4+4
	_locationBirthdayBytesLen = 4  // year+month+day == 2+1+1
)

// ParseLocation parse location bytes to id
func ParseLocation(bs []byte) (country, province, city uint32) {
	if len(bs) != _locationBytesLen {
		return
	}
	country = binary.BigEndian.Uint32(bs[:4])
	province = binary.BigEndian.Uint32(bs[4:8])
	city = binary.BigEndian.Uint32(bs[8:])
	return
}

// ParseBirthday parse birthday bytes to number
func ParseBirthday(bs []byte) (year, month, day uint32) {
	if len(bs) != _locationBirthdayBytesLen {
		return
	}
	year = uint32(binary.BigEndian.Uint16(bs[:2]))
	month = uint32(bs[2])
	day = uint32(bs[3])
	return
}
