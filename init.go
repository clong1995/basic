package basic

import (
	"basic/color"
	"fmt"
)

func Init() {
	fmt.Println(color.Magenta, "\n  _                  _         __     ___     ___  \n | |                (_)       /_ |   / _ \\   / _ \\ \n | |__    __ _  ___  _   ___   | |  | | | | | | | |\n | '_ \\  / _` |/ __|| | / __|  | |  | | | | | | | |\n | |_) || (_| |\\__ \\| || (__   | | _| |_| |_| |_| |\n |_.__/  \\__,_||___/|_| \\___|  |_|(_)\\___/(_)\\___/ \n                                                   ", color.Reset)
}
