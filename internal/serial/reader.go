package serial

import "machine"

type buffer struct {
	data [128]byte
	n    int
}

func (b *buffer) add(c byte) bool {
	if c == '\n' {
		return true
	}
	if b.n < len(b.data) {
		b.data[b.n] = c
		b.n++
	}
	return false
}

func (b *buffer) string() string {
	return string(b.data[:b.n])
}

func (b *buffer) reset() {
	b.n = 0
}

type Reader struct {
	buf buffer
}

func (r *Reader) Read(serial machine.Serialer) bool {
	for serial.Buffered() > 0 {
		b, _ := serial.ReadByte()
		if r.buf.add(b) {
			return true
		}
	}
	return false
}

func (r *Reader) Line() string {
	s := r.buf.string()
	r.buf.reset()
	return s
}
