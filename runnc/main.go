package main

/*
typedef struct foo {
    int a;
    int b;
} Foo;
*/
import "C"

import "fmt"

//export Foo
type Foo struct {
	a int
	b int
}

//export PrintInt
func PrintInt(x int) C.struct_foo {
	fmt.Println(x)
	s := C.struct_foo{1, 2}
	return s
}

func main() {}
