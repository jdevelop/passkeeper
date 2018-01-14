package device

// Keyboard defines the methods to use for writing the character to the keyboard.
type Keyboard interface {
	WriteString(content string) (err error)
}
