package color

import "fmt"

func Success(str string) {
	fmt.Println(Green, str, Reset)
}
func Fail(str string) {
	fmt.Println(Red, str, Reset)
}
