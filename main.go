// Command template-go is a minimal example program.
package main

import (
	"fmt"

	"github.com/napolitain/template-go/internal/greet"
)

func main() {
	fmt.Println(greet.Hello("world"))
}
