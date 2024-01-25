package printf

import (
	"fmt"
	"testing"
	"time"
)

func TestTextColor(t *testing.T) {
	fmt.Println(SetTextColor("hello", TextBlue))
	fmt.Println(SetTextColor("hello", TextWhite))
	fmt.Println(SetTextColor("hello", TextRed))
	fmt.Println("=======")
	FpTextColor(fmt.Println, "hello", TextGreen)
	fmt.Println("==========")
	fmt.Println("==========")
	fmt.Println("==========")
	fmt.Println("==========")
	fmt.Println("==========")
	time.Sleep(time.Second)
	MoveCursorUpLines(3)
	time.Sleep(time.Second)
	FpTextColor(fmt.Println, "world", TextGreen)
	time.Sleep(time.Second)
	MoveCursorDownLines(1)
}
