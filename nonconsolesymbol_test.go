package colorable

import (
	"bytes"
	"testing"
)

func TestNonConsoleSymbolWriter(t *testing.T) {
	var buf bytes.Buffer
	want := "hello"
	NewNonConsoleSymbolWriter(&buf, NonConsoleSymbolConf{
		Ignore: []ConsoleSymbol{T_m, T_J},
	}).Write([]byte("\x1b[0m" + want + "\x1b[2J"))
	got := buf.String()
	if got != want {
		t.Fatalf("want %q but %q", want, got)
	}

	buf.Reset()
	NewNonConsoleSymbolWriter(&buf, NonConsoleSymbolConf{
		Ignore: []ConsoleSymbol{T_m},
	}).Write([]byte("\x1b[0m"))
	got = buf.String()
	if got != "" {
		t.Fatalf("want %q but %q", "", got)
	}
}

func TestNonConsoleSymbolReader(t *testing.T) {
	var buf bytes.Buffer
	want := "hello"
	_, err := buf.ReadFrom(NewNonConsoleSymbolReader(bytes.NewReader([]byte("\x1b[0m"+want+"\x1b[2J")), NonConsoleSymbolConf{
		Ignore: []ConsoleSymbol{T_m, T_J},
	}))
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Fatalf("want %q but %q", want, got)
	}
}
