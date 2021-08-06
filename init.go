package basic

import (
	"basic/color"
	"fmt"
)

func Init() {
	//输出文字图形
	//http://patorjk.com/software/taag/#p=display&f=Big&t=Type%20Something%20
	//basic 1.0.0
	fmt.Println(color.Magenta, "\n  _                     _          __        ___        __ \n | |                   (_)        /_ |      / _ \\      /_ |\n | |__     __ _   ___   _    ___   | |     | | | |      | |\n | '_ \\   / _` | / __| | |  / __|  | |     | | | |      | |\n | |_) | | (_| | \\__ \\ | | | (__   | |  _  | |_| |  _   | |\n |_.__/   \\__,_| |___/ |_|  \\___|  |_| (_)  \\___/  (_)  |_|\n                                                           ", color.Reset)
}
