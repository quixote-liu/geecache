package geecache

type ByteViews struct {
	b []byte
}

func (v ByteViews) Len() int {
	return len(v.b)
}

func (v ByteViews) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (v ByteViews) String() string {
	return string(v.b)
}
