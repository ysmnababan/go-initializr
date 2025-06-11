package pkg

type Rune rune

func (c Rune) IsAlphabet() bool {
	return c <= 90 && c >= 65 || c >= 97 && c <= 122
}
func (c Rune) IsSpace() bool {
	return c == ' '
}
func (c Rune) IsSlash() bool {
	return c == '/'
}
