package dye

// String Возвращение последовательности в виде строки.
func (s sequence) String() string { return string(s) }

// Byte Возвращение последовательности в виде среза байт.
func (s sequence) Byte() []byte { return []byte(s) }
