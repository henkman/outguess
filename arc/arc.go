package arc

import "crypto/md5"

const (
	N = 256
)

type Stream struct {
	i, j uint8
	s    [N]uint8
}

func (as *Stream) Init() {
	for n := 0; n < N; n++ {
		as.s[n] = uint8(n)
	}
}

func (as *Stream) GetByte() uint8 {
	as.i++
	as.j += as.s[as.i]
	as.s[as.i], as.s[as.j] = as.s[as.j], as.s[as.i]
	return as.s[(as.s[as.i]+as.s[as.j])&0xFF]
}

func (as *Stream) GetWord() uint32 {
	val := uint32(as.GetByte()) << 24
	val |= uint32(as.GetByte()) << 16
	val |= uint32(as.GetByte()) << 8
	val |= uint32(as.GetByte())
	return val
}

func (as *Stream) AddRandom(dat []byte) {
	dl := len(dat)
	as.i--
	for n := 0; n < N; n++ {
		as.i++
		as.j += as.s[as.i] + dat[n%dl]
		as.s[as.i], as.s[as.j] = as.s[as.j], as.s[as.i]
	}
}

func (as *Stream) InitKey(typ, key []byte) {
	hash := md5.New()
	hash.Write(typ)
	hash.Write(key)
	digest := hash.Sum(nil)
	as.Init()
	as.AddRandom(digest)
}
