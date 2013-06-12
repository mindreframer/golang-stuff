package main

import (
    "fmt"
    "log"
)

type Tuple struct {
    Left, Right interface{}
}

func (t Tuple) String() string {
    log.Printf("in Stringer interface method for Tuple")
    return fmt.Sprintf("(%#v, %#v)", t.Left, t.Right)
}

type Tuple2 struct {
    Left, Right interface{}
}

func (t Tuple2) Error() string {
    log.Printf("in Error interface method for Tuple2")
    return "lol it's an error!"
}

func (t Tuple2) String() string {
    log.Printf("in Stringer interface method for Tuple2")
    return fmt.Sprintf("(%#v, %#v)", t.Left, t.Right)
}

func main() {
    fmt.Printf("%s\n", Tuple{1, 2})
    fmt.Printf("%s\n", Tuple2{1.5, 2.1})
    fmt.Printf("%v\n", Tuple{"Bruce Wayne", "Batman"})
}
