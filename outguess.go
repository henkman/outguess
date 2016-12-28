package outguess

import (
	"errors"
	"fmt"
	"io"

	"github.com/henkman/outguess/arc"
	"github.com/henkman/outguess/stegjpeg"
)

type iterator struct {
	as      arc.Stream
	skipmod uint32
	off     uint32
}

func (it *iterator) Current() uint32 {
	return it.off
}

func (it *iterator) Next() uint32 {
	it.off += (it.as.GetWord() % it.skipmod) + 1
	return it.off
}

func (it *iterator) Seed(seed uint16) {
	reseed := [2]byte{
		uint8(seed),
		uint8(seed >> 8),
	}
	it.as.AddRandom(reseed[:])
}

func (it *iterator) Adapt(bm []byte, bits uint32, datalen uint) {
	bo := bits - it.off
	var sa float32
	if bo > bits/32 {
		sa = 2
	} else {
		sa = 2 - float32((bits/32)-bo)/float32(bits/32)
	}
	it.skipmod = uint32(sa * float32(bits-it.off) / float32(8*datalen))
}

func steg_retrbyte(bm []byte, bits uint32, it *iterator) uint32 {
	var tmp, bit uint32
	i := it.Current()
	for bit = 0; bit < bits; bit++ {
		if bm[i/8]&(1<<(i&7)) != 0 {
			tmp |= 1 << bit
		}
		i = it.Next()
	}
	return tmp
}

func Get(file io.Reader, msg io.Writer, key []byte) error {
	if key == nil {
		key = []byte("Default key")
	}
	_, bm, err := stegjpeg.Decode(file)
	if err != nil {
		return err
	}
	var as arc.Stream
	as.InitKey([]byte("Encryption"), key)
	tas := as

	var it iterator
	it.skipmod = 32
	it.as.InitKey([]byte("Seeding"), key)
	it.off = it.as.GetWord() % it.skipmod

	var buf [4]byte
	for i, _ := range buf {
		buf[i] = byte(steg_retrbyte(bm.Bitmap, 8, &it)) ^ as.GetByte()
	}
	seed := uint16(buf[0]) | uint16(buf[1])<<8
	datalen := uint(buf[2]) | uint(buf[3])<<8
	if datalen > uint(bm.Bytes) {
		return errors.New(fmt.Sprintf(
			"Extracted datalen is too long: %d > %d\n", datalen, len(bm.Bitmap)))
	}
	encdata := make([]byte, datalen)
	it.Seed(seed)
	var n uint
	for datalen > 0 {
		it.Adapt(bm.Bitmap, bm.Bits, datalen)
		encdata[n] = byte(steg_retrbyte(bm.Bitmap, 8, &it))
		n++
		datalen--
	}
	for i, _ := range encdata {
		encdata[i] ^= tas.GetByte()
	}
	_, err = msg.Write(encdata)
	return err
}

func Put(file, msg io.Reader, key []byte, out io.Writer) error {
	return errors.New("writing not yet implemented")
}
