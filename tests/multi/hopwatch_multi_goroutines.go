package main

import (
	"github.com/emicklei/hopwatch"
	"log"
)

func main() {
	ready := make(chan int)
	for id := 0 ; id < 4 ; id++ {
		log.Printf("spawn doit:%v",id) 
		go doit(id, ready)
	}	
	for j := 0 ; j < 4 ; j++ {
		who := <- ready
		log.Printf("done:%v", who)
	}
}

func doit(id int, ready chan int) {
	log.Printf("before break:%v",id)
	hopwatch.Display("id",id).Break()
	log.Printf("after break:%v",id)
	ready <- id
}

