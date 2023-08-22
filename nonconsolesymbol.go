package colorable

import (
	"bytes"
	"io"

	"github.com/zijiren233/stream"
)

var _ io.Writer = &NonConsoleSymbolWriter{}

type ConsoleSymbol byte

const (
	T_A ConsoleSymbol = 'A' // 光标上移 \ cursor up
	T_B ConsoleSymbol = 'B' // 光标下移 \ cursor down
	T_C ConsoleSymbol = 'C' // 光标右移 \ cursor right
	T_D ConsoleSymbol = 'D' // 光标左移 \ cursor left
	T_E ConsoleSymbol = 'E' // 光标下移到下一行的开头 \ cursor down to the beginning of the next line
	T_F ConsoleSymbol = 'F' // 光标上移到上一行的开头 \ cursor up to the beginning of the previous line
	T_G ConsoleSymbol = 'G' // 光标移动到当前行的指定列 \ cursor moves to the specified column of the current line
	T_H ConsoleSymbol = 'H' // 光标移动到指定行和列 \ cursor moves to the specified row and column
	T_J ConsoleSymbol = 'J' // 清屏或清除从光标到屏幕的一部分 \ clear the screen or clear part of the screen from the cursor
	T_K ConsoleSymbol = 'K' // 清除从光标到行尾的一部分 \ clear part of the line from the cursor to the end of the line
	T_S ConsoleSymbol = 'S' // 向上滚动指定的行数 \ scroll up the specified number of lines
	T_T ConsoleSymbol = 'T' // 向下滚动指定的行数 \ scroll down the specified number of lines
	T_m ConsoleSymbol = 'm' // 设置颜色和格式 \ set color and format
	T_s ConsoleSymbol = 's' // 保存光标位置 \ save cursor position
	T_u ConsoleSymbol = 'u' // 恢复光标位置 \ restore cursor position
)

type NonConsoleSymbolConf struct {
	Ignore    []ConsoleSymbol
	IgnoreAll bool
}

type NonConsoleSymbolWriter struct {
	target io.Writer
	conf   NonConsoleSymbolConf
	buf    *bytes.Buffer
}

func NewNonConsoleSymbolWriter(w io.Writer, conf NonConsoleSymbolConf) io.Writer {
	return &NonConsoleSymbolWriter{
		target: w,
		conf:   conf,
		buf:    bytes.NewBuffer(nil),
	}
}

func (w *NonConsoleSymbolWriter) Write(p []byte) (n int, err error) {
	r := bytes.NewReader(p)
	var c byte
	for {
		switch {
		case w.buf.Len() == 0:
			if c, err = r.ReadByte(); err != nil {
				if err == io.EOF {
					err = nil
				}
				return
			}
			n++
			if c != 0x1b {
				if _, err = w.target.Write([]byte{c}); err != nil {
					return
				}
				continue
			}
			if err = w.buf.WriteByte(c); err != nil {
				return
			}
			fallthrough
		case w.buf.Len() == 1:
			if c, err = r.ReadByte(); err != nil {
				if err == io.EOF {
					err = nil
				}
				return
			}
			n++
			if c != 0x5b {
				_, err = w.target.Write([]byte{0x1b, c})
				w.buf.Reset()
				if err != nil {
					return
				}
				continue
			}
			if err = w.buf.WriteByte(c); err != nil {
				return
			}
			fallthrough
		default:
			for {
				if c, err = r.ReadByte(); err != nil {
					if err == io.EOF {
						err = nil
					}
					return
				}
				n++
				if err = w.buf.WriteByte(c); err != nil {
					return
				}
				if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' {
					if w.conf.IgnoreAll {
						w.buf.Reset()
						break
					}
					csc := ConsoleSymbol(c)
					ignore := false
					for _, cs := range w.conf.Ignore {
						if cs == csc {
							ignore = true
							break
						}
					}
					if ignore {
						w.buf.Reset()
						break
					}
					newSlice := make([]byte, w.buf.Len())
					copy(newSlice, w.buf.Bytes())
					_, err = w.target.Write(newSlice)
					w.buf.Reset()
					if err != nil {
						return
					}
					break
				}
			}
		}
	}
}

type NonConsoleSymbolReader struct {
	source *stream.Reader
	conf   NonConsoleSymbolConf
	buf    *bytes.Buffer
}

func NewNonConsoleSymbolReader(r io.Reader, conf NonConsoleSymbolConf) *NonConsoleSymbolReader {
	return &NonConsoleSymbolReader{
		source: stream.NewReader(r, stream.BigEndian),
		conf:   conf,
		buf:    bytes.NewBuffer(nil),
	}
}

func (r *NonConsoleSymbolReader) Read(p []byte) (n int, err error) {
	n = copy(p[n:], r.buf.Bytes())
	r.buf.Truncate(r.buf.Len() - n)
	var c byte
	maxI := len(p) - 1
	for {
		if n > maxI {
			return
		}
		if c, err = r.source.ReadByte(); err != nil {
			return
		}
		if c != 0x1b {
			p[n] = c
			n++
			continue
		}
		if err = r.buf.WriteByte(c); err != nil {
			return
		}
		if c, err = r.source.ReadByte(); err != nil {
			return
		}
		if err = r.buf.WriteByte(c); err != nil {
			return
		}
		if c != 0x5b {
			tn := copy(p[n:], r.buf.Bytes())
			n += tn
			r.buf.Truncate(r.buf.Len() - tn)
			continue
		}
		for {
			if c, err = r.source.ReadByte(); err != nil {
				return
			}
			if err = r.buf.WriteByte(c); err != nil {
				return
			}
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' {
				if r.conf.IgnoreAll {
					r.buf.Reset()
					break
				}
				csc := ConsoleSymbol(c)
				ignore := false
				for _, cs := range r.conf.Ignore {
					if cs == csc {
						ignore = true
						break
					}
				}
				if ignore {
					r.buf.Reset()
					break
				}
				tn := copy(p[n:], r.buf.Bytes())
				n += tn
				r.buf.Truncate(r.buf.Len() - tn)
				break
			}
		}
	}
}
