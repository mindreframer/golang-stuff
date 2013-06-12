package main

import "log"

func main() {
    m := make(map[string]int)
    log.Println(m)

    m["one"] = 1
    log.Println(m)

    m["two"] = 2
    log.Println(m)

    delete(m, "one")
    log.Println(m)

    delete(m, "one")
    log.Println(m)

    m = nil
    delete(m, "two") // Will panic
}
