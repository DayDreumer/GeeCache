package GeeCache

type ByteView struct {
	b []byte
}

func (bv ByteView) Len() int {
	return len(bv.b)
}

func (bv ByteView) ByteSlice() []byte {
	return cloneByte(bv.b)
}

func (bv ByteView) String() string {
	return string(bv.b)
}

func cloneByte(b []byte) []byte {
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
}
