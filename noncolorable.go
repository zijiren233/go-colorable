package colorable

import (
	"io"
)

// NewNonColorable returns new instance of Writer which removes escape sequence from Writer.
func NewNonColorableWriter(w io.Writer) io.Writer {
	return NewNonConsoleSymbolWriter(
		w,
		NonConsoleSymbolConf{
			Ignore: []ConsoleSymbol{T_m},
		},
	)
}

func NewNonColorableReader(r io.Reader) io.Reader {
	return NewNonConsoleSymbolReader(
		r,
		NonConsoleSymbolConf{
			Ignore: []ConsoleSymbol{T_m},
		},
	)
}
