package connection

// address type enum
const (
	AddressTypeIPV4   AddressType = 0x01
	AddressTypeDomain AddressType = 0x03
	AddressTypeIPV6   AddressType = 0x04
)

//AddressType request address type
type AddressType byte

// Valid is valid type
func (tpy AddressType) Valid() bool {
	switch tpy {
	case AddressTypeIPV4, AddressTypeDomain, AddressTypeIPV6:
		return true
	default:
		return false
	}
}

// String return string form
func (tpy AddressType) String() string {
	switch tpy {
	case AddressTypeIPV4:
		return "ipv4"
	case AddressTypeDomain:
		return "domain"
	case AddressTypeIPV6:
		return "ipv6"
	default:
		return "unknown"
	}
}
