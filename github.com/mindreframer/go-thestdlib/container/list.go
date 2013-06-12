package main

import (
    "container/list"
    "log"
)

const size = 5

func Do(l *list.List, f func(interface{})) {
    // Standard list iterating straight from their example
    for e := l.Front(); e != nil; e = e.Next() {
        f(e.Value)
    }
}

func printList(l *list.List) {
    elements := make([]interface{}, 0, l.Len())
    Do(l, func(i interface{}) {
        elements = append(elements, i)
    })
    log.Printf("%v", elements)
}

func main() {
    l := list.New()
    printList(l) // []
    for i := 0; i < size; i++ {
        l.PushBack(i)
    }
    printList(l) // [0 1 2 3 4]

    l = l.Init()
    for i := 0; i < size; i++ {
        l.PushFront(i)
    }
    printList(l) // [4 3 2 1 0]

    f := l.Front()
    e := f.Next().Next()
    e = l.InsertAfter(10, e)
    printList(l)         // [4 3 2 10 1 0]
    log.Println(l.Len()) // 6

    l.Remove(e.Next())
    printList(l)         // [4 3 2 10 0]
    log.Println(l.Len()) // 5
}
