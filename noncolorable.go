package colorable

import (
	"bytes"
	"io"
)

// NonColorable holds writer but removes escape sequence.
type NonColorableWriter struct {
	out io.Writer
}

type NonColorableReader struct {
	in io.Reader
}

// NewNonColorable returns new instance of Writer which removes escape sequence from Writer.
func NewNonColorableWriter(w io.Writer) io.Writer {
	return &NonColorableWriter{out: w}
}

func NewNonColorableReader(r io.Reader) io.Reader {
	return &NonColorableReader{in: r}
}

// Write writes data on console
func (w *NonColorableWriter) Write(data []byte) (n int, err error) {
	return NonColorWrite(w.out, data)
}

func (w *NonColorableReader) Read(data []byte) (n int, err error) {
	return NonColorRead(w.in, data)
}

func NonColorWrite(out io.Writer, data []byte) (n int, err error) {
	er := bytes.NewReader(data)
	var plaintext bytes.Buffer
loop:
	for {
		c1, err := er.ReadByte()
		if err != nil {
			plaintext.WriteTo(out)
			break loop
		}
		if c1 != 0x1b {
			plaintext.WriteByte(c1)
			continue
		}
		_, err = plaintext.WriteTo(out)
		if err != nil {
			break loop
		}
		c2, err := er.ReadByte()
		if err != nil {
			break loop
		}
		if c2 != 0x5b {
			continue
		}

		for {
			c, err := er.ReadByte()
			if err != nil {
				break loop
			}
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' {
				break
			}
		}
	}

	return len(data), nil
}

func NonColorRead(in io.Reader, data []byte) (n int, err error) {
	n, err = in.Read(data)
	er := bytes.NewReader(data[:n])
	n = 0
loop:
	for {
		c1, err := er.ReadByte()
		if err != nil {
			break loop
		}
		if c1 != 0x1b {
			data[n] = c1
			n++
			continue
		}
		c2, err := er.ReadByte()
		if err != nil {
			break loop
		}
		if c2 != 0x5b {
			continue
		}

		for {
			c, err := er.ReadByte()
			if err != nil {
				break loop
			}
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' {
				break
			}
		}
	}
	return
}
