package main

import "log"

func main() {
    // Empty slice, with capacity of 10
    ints := make([]int, 0, 10)
    log.Printf("ints: %v", ints)

    ints2 := append(ints, 1, 2, 3)

    log.Printf("ints2: %v", ints2)
    log.Printf("Slice was at %p, it's probably still at %p", ints, ints2)

    moreInts := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
    ints3 := append(ints2, moreInts...)

    log.Printf("ints3: %v", ints3)
    log.Printf("Slice was at %p, and it moved to %p", ints2, ints3)

    ints4 := []int{1, 2, 3}
    log.Printf("ints4: %v", ints4)
    // The idiomatic way to append to a slice,
    // just assign to the same variable again
    ints4 = append(ints4, 4, 5, 6)
    log.Printf("ints4: %v", ints4)
}
