package rpi

import (
	"testing"
)

func TestScancodesInit(t *testing.T) {

	var tests = []struct {
		char       rune
		code, mods byte
	}{
		{'A', 0x04, 0x02},
		{'a', 0x04, 0x00},
		{'Z', 0x1d, 0x02},
		{'z', 0x1d, 0x00},
		{'1', 0x1e, 0x00},
		{'!', 0x1e, 0x02},
		{'2', 0x1f, 0x00},
		{'@', 0x1f, 0x02},
		{'9', 0x26, 0x00},
		{'(', 0x26, 0x02},
		{'0', 0x27, 0x00},
		{')', 0x27, 0x02},
		{' ', 0x2C, 0x00},
		{'=', 0x2e, 0x00},
		{'+', 0x2e, 0x02},
		{'[', 0x2f, 0x00},
		{'{', 0x2f, 0x02},
		{']', 0x30, 0x00},
		{'}', 0x30, 0x02},
		{'/', 0x38, 0x00},
		{'?', 0x38, 0x02},
	}

	for _, test := range tests {
		m, err := ResolveScanKey(test.char)
		if err != nil {
			t.Error(err)
		}
		if m[0] != test.code {
			t.Errorf("Expected %x actual %d", test.code, m[0])
		}
		if m[1] != test.mods {
			t.Errorf("Expected %x actual %d", test.mods, m[1])
		}
	}

}
