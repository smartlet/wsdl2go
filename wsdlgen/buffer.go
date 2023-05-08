package wsdlgen

import (
	"fmt"
	"io"
	"sync"
	"unicode/utf8"
)

const (
	defaultBufferMinimum = 2048
	defaultBufferMaximum = defaultBufferMinimum * 4
)

type Buffer struct {
	off int
	len int
	cap int
	buf []byte
}

func NewBuffer(c int) *Buffer {
	return &Buffer{
		off: 0,
		len: 0,
		cap: c,
		buf: make([]byte, c),
	}
}

func (b *Buffer) Line(format string, args ...any) {
	fmt.Fprintf(b, format, args...)
	b.WriteByte('\n')
}

func (b *Buffer) Len() int {
	return b.len
}

func (b *Buffer) Cap() int {
	return b.cap
}

func (b *Buffer) Buff() []byte {
	return b.buf[:b.len]
}

func (b *Buffer) Bytes() []byte {
	return b.buf[b.off:b.len]
}

func (b *Buffer) String() string {
	return string(b.buf[b.off:b.len])
}

func (b *Buffer) Reset() {
	b.off = 0
	b.len = 0
}

func (b *Buffer) Read(p []byte) (int, error) {
	if b.off >= b.len {
		return 0, io.EOF
	}
	n := copy(p, b.buf[b.off:b.len])
	b.off += n
	return n, nil
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.off >= b.len {
		return 0, io.EOF
	}
	c := b.buf[b.off]
	b.off++
	return c, nil
}

func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.off = int(offset)
	case io.SeekCurrent:
		b.off += int(offset)
	case io.SeekEnd:
		b.off = b.len - int(offset)
	}
	return int64(b.off), nil
}

func (b *Buffer) Close() error {
	return nil
}

func (b *Buffer) Grow(n int) {
	buf := b.buf
	b.cap += n
	b.buf = make([]byte, b.cap)
	copy(b.buf, buf)
}

func (b *Buffer) Write(p []byte) (int, error) {
	n := len(p)
	if b.len+n > b.cap {
		b.Grow(b.cap + n)
	}
	b.len += copy(b.buf[b.len:], p)
	return n, nil
}

func (b *Buffer) WriteString(s string) (n int, err error) {
	n = len(s)
	if b.len+n > b.cap {
		b.Grow(b.cap + n)
	}
	b.len += copy(b.buf[b.len:], s)
	return n, nil
}

func (b *Buffer) WriteByte(c byte) error {
	if b.len >= b.cap {
		b.Grow(b.cap + 1)
	}
	b.buf[b.len] = c
	b.len++
	return nil
}

func (b *Buffer) WriteRune(r rune) (n int, err error) {
	// Compare as uint32 to correctly handle negative runes.
	if uint32(r) < utf8.RuneSelf {
		b.WriteByte(byte(r))
		return 1, nil
	}
	if b.len+utf8.UTFMax > b.cap {
		b.Grow(b.cap + utf8.UTFMax)
	}
	n = utf8.EncodeRune(b.buf[b.len:b.len+utf8.UTFMax], r)
	b.len += n
	return n, nil
}

var _ io.Reader = (*Buffer)(nil)
var _ io.Writer = (*Buffer)(nil)
var _ io.Seeker = (*Buffer)(nil)
var _ io.Closer = (*Buffer)(nil)

var pool = sync.Pool{
	New: func() any {
		return NewBuffer(defaultBufferMinimum)
	},
}

func borrowBuffer() *Buffer {
	buf := pool.Get().(*Buffer)
	buf.Reset()
	return buf
}

func returnBuffer(b *Buffer) {
	if b.cap < defaultBufferMaximum {
		pool.Put(b)
	}
}
