package printer

import "fmt"

func PrintPlain(output []byte) {
	fmt.Print(string(output))
}
