package printf

import "fmt"

/*
conf := 0   // 配置、终端默认设置
bg   := 0   // 背景色、终端默认设置
text := 31  // 前景色、红色
fmt.Printf("\n %c[%d;%d;%dm%s%c[0m\n\n", 0x1B, conf, bg, text, "testPrintColor", 0x1B)

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见



*/
const (
	TextBlack TextColor = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

type TextColor int

const (
	BackgroundDefault TextColor = iota
	BackgroundHighlight
	BackgroundUnderline = iota + 2
	BackgroundFlicker
	BackgroundReverseShow = iota + 3
	BackgroundNotShow
)

func SetColorWithBackground(msg string, conf, bg, textColor TextColor) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, conf, bg, textColor, msg, 0x1B)
}

func SetTextColor(msg string, textColor TextColor) string {
	return SetColorWithBackground(msg, BackgroundDefault, BackgroundDefault, textColor)
}

func FpTextColor(fp func(a ...any) (n int, err error), msg string, textColor TextColor) (n int, err error) {
	return fp(SetColorWithBackground(msg, BackgroundDefault, BackgroundDefault, textColor))
}
