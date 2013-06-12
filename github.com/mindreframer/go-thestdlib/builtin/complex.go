package main

import "log"

func main() {
    c1 := 1.5 + 0.5i
    c2 := complex(1.5, 0.5)
    log.Printf("c1: %v", c1)
    log.Printf("c2: %v", c2)
    log.Printf("c1 == c2: %v", c1 == c2)
    log.Printf("c1 real: %v", real(c1))
    log.Printf("c1 imag: %v", imag(c1))
    log.Printf("c1 + c2: %v", c1+c2)
    log.Printf("c1 - c2: %v", c1-c2)
    log.Printf("c1 * c2: %v", c1*c2)
    log.Printf("c1 / c2: %v", c1/c2)
    log.Printf("c1 type: %T", c1)

    c3 := complex(float32(1.5), float32(0.5))
    log.Printf("c3 type: %T", c3)
}
