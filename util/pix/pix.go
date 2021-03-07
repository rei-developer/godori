package pix

import "unicode/utf8"

func Maker(text string, x string, y string) string {
	if ([]rune(text)[utf8.RuneCountInString(text)-1]-44032)%28 == 0 {
		return text + x
	} else {
		return text + y
	}
}
