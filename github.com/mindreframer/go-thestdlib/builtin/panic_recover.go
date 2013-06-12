package main

import (
    "errors"
    "log"
)

func handlePanic(f func()) {
    defer func() {
        if r := recover(); r != nil {
            if str, ok := r.(string); ok {
                log.Printf("got a string error: %s", str)
                return
            }

            if err, ok := r.(error); ok {
                log.Printf("got an error error: %s", err.Error())
                return
            }

            log.Printf("got a different kind of error: %v", r)
        }
    }()
    f()
}

func main() {
    handlePanic(func() {
        panic("string error")
    })

    handlePanic(func() {
        panic(errors.New("error error"))
    })

    handlePanic(func() {
        panic(10)
    })
}
