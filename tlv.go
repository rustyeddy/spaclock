package main

//
// TLV stands for Type Length Value (or vector), the first byte is the
// length, the second byte is the type and the remaining values,
// depending on the type.
//
// There is a micro controller extension, that says all messages with
// the first bit of the type set are single byte messages. All
// messages with the 2nd MSB set are two byte messages.
//
// Illustration:
//
// - type: 0x00 - 0x3f (  0 -  63) three bytes (128 - 255 are 3 bytes)
// - type: 0x40 - 0x7f ( 64 - 127) three bytes
// - type: 0x80 - 0xbf (128 - 191) two bytes
// - type: 0xc0 - 0xff (192 - 255) single byte (type and value are
//
// Reserved Types:
//
// - type: 0x00 = hello, len=1, value=X (x == don't care)
// - type: 0xff = RESET 0xff03dead
//
const (
	tlvTypeHello   byte = 0x00
	tlvTypeMessage byte = 0x01
	tlvTypeTime    byte = 0x02
	tlvTypeDate    byte = 0x03
	tlvTypeTempf   byte = 0x04
)

type TLV struct {
	Buffer []byte
}

// NewTLV will
func NewTLV(t byte, l int, str string) TLV {
	nbuf := make([]byte, l)
	nbuf[0] = t
	nbuf[1] = byte(l)
	nbuf = append(nbuf, []byte(str)...)
	return TLV{nbuf}
}

func (tlv *TLV) Type() (ttype byte) {
	return tlv.Buffer[0]
}

func (tlv *TLV) Len() int {
	blen := len(tlv.Buffer)
	return blen
}

func (tlv *TLV) Value() []byte {
	return tlv.Buffer[2:]
}

func (tlv *TLV) String() string {
	return string(tlv.Buffer[2:])
}

func (tlv *TLV) Marshal(buf []byte) {
	tlv = &TLV{buf}
}
