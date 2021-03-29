package main

import (
	"embed"
	"fmt"
)

//go:embed dir hello1.txt
var dir embed.FS

func main() {
	// f, err := dir.Open("hello.txt")
	// f, err := dir.Open("subdir/hello.txt")
	// f, err := dir.Open("dir/subdir/hello.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// c, _ := ioutil.ReadAll(f)
	// fmt.Println("c : ", string(c))

	c1, _ := dir.ReadFile("hello1.txt")
	fmt.Println("c1 : ", string(c1))

	c2, _ := dir.ReadFile("dir/subdir/hello.txt")
	fmt.Println("c2 : ", string(c2))
}
