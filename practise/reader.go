package practise

import (
	"fmt"
	"io"
	"strings"
)

func ReaderTest() {
	r := strings.NewReader("Hello, Reader!")
	b := make([]byte, 8)
	var full []byte
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
		full = append(full, b[:n]...)
	}
	fmt.Printf("full[] = %q\n", full)
}
