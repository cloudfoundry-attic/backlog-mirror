package main

import (
	"errors"
	"fmt"
)


func main() {
	fmt.Println("fail", errors.New("no actual mirroring!"))
}
