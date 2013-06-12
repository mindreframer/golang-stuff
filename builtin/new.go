package main

import "log"

type Actor struct {
    Name string
}

type Movie struct {
    Title  string
    Actors []*Actor
}

func main() {
    ip := new(int)
    log.Printf("ip type: %T, ip: %v, *ip: %v", ip, ip, *ip)

    m := new(Movie)
    log.Printf("m type: %T, m: %v, *m: %v", m, m, *m)
}
