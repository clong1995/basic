package basic

import (
	"fmt"
	"github.com/clong1995/basic/color"
)

func Init() {
	//输出文字图形
	//http://patorjk.com/software/taag/#p=display&f=Big&t=Type%20Something%20
	//basic 1.0.0
	fmt.Println(color.Magenta, "\n  _                     _          __       __        ___  \n | |                   (_)        /_ |     /_ |      / _ \\ \n | |__     __ _   ___   _    ___   | |      | |     | | | |\n | '_ \\   / _` | / __| | |  / __|  | |      | |     | | | |\n | |_) | | (_| | \\__ \\ | | | (__   | |  _   | |  _  | |_| |\n |_.__/   \\__,_| |___/ |_|  \\___|  |_| (_)  |_| (_)  \\___/ \n                                                           ", color.Reset)
}
