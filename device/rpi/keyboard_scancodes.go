package rpi

import "fmt"

var scanKeyMap = make(map[rune][2]byte)

func init() {

	for sym := 'A'; sym <= 'Z'; sym++ {
		symIdx := byte(sym-'A') + 4
		scanKeyMap[sym] = [...]byte{symIdx, 0x02}
		scanKeyMap[sym-'A'+'a'] = [...]byte{symIdx, 0x0}
	}

	var numShift = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '('}
	for sym := '1'; sym <= '9'; sym++ {
		symIdx := byte(sym-'1') + 0x1e
		scanKeyMap[sym] = [...]byte{symIdx, 0x0}
		scanKeyMap[numShift[sym-'1']] = [...]byte{symIdx, 0x02}
	}

	scanKeyMap['0'] = [...]byte{0x27, 0x0}
	scanKeyMap[')'] = [...]byte{0x27, 0x02}

	var pairs = [...][]rune{
		{'-', '_'},
		{'=', '+'},
		{'[', '{'},
		{']', '}'},
		{'\\', '|'},
		nil,
		{';', ':'},
		{'\'', '"'},
		{'`', '~'},
		{',', '<'},
		{'.', '>'},
		{'/', '?'},
	}

	for sym := byte(0x2d); sym <= 0x38; sym++ {
		symIdx := sym - 0x2d
		if pairs[symIdx] == nil {
			continue
		}
		scanKeyMap[pairs[symIdx][0]] = [...]byte{sym, 0x0}
		scanKeyMap[pairs[symIdx][1]] = [...]byte{sym, 0x2}
	}

	scanKeyMap[' '] = [...]byte{0x2c, 0x0}

}

func ResolveScanKey(key rune) ([]byte, error) {
	if res, ok := scanKeyMap[key]; ok {
		return res[:], nil
	}
	return nil, fmt.Errorf("Can't resolve key %c", key)
}
