package printf

import (
	"bytes"
	"fmt"
)

// 光标向上移动几行
func MoveCursorUpLines(count int) {
	fmt.Printf("\033[%dA\033[K", count)
}

// 光标向下移动几行
func MoveCursorDownLines(count int) {
	fmt.Printf(string(bytes.Repeat([]byte{'\n'}, count)))
}
