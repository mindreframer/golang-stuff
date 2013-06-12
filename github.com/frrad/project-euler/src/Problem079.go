package main

import (
	"./eulerlib"
	"fmt"
	"math/rand"
	"time"
)

const mergesame = 1 //Could be <1 a priori

func works(key string, secret string) bool {
	a := -1
	for i := len(secret) - 1; i > -1; i-- {
		if secret[i] == key[0] {
			a = i
		}
	}
	if a == -1 {
		return false

	}
	b := -1

	for i := len(secret) - 1; i > a; i-- {
		if secret[i] == key[1] {
			b = i
		}
	}

	if b == -1 {
		return false

	}

	for i := len(secret) - 1; i > b; i-- {
		if secret[i] == key[2] {
			return true
		}
	}

	return false

}

func merge(string1 string, string2 string) string {

	if len(string1) == 0 {
		return string2

	}

	if len(string2) == 0 {
		return string1
	}

	if len(string1) == 1 && len(string2) == 1 {
		if rand.Float32() < .5 {
			return string1 + string2
		}
		return string2 + string1
	}

	for i := 0; i < len(string1); i++ {
		for j := 0; j < len(string2); j++ {
			if string1[i] == string2[j] && rand.Float32() < mergesame {

				return merge(string1[:i], string2[:j]) + string(string1[i]) + merge(string1[i+1:], string2[j+1:])

			}
		}
	}

	a := rand.Int() % len(string1)
	b := rand.Int() % len(string2)
	return merge(string1[:a], string2[:b]) + merge(string1[a:], string2[b:])

}

func allwork(keys []string, secret string) bool {
	if len(keys) == 1 {
		return works(keys[0], secret)
	}
	return works(keys[0], secret) && allwork(keys[1:], secret)
}

func main() {

	rand.Seed(time.Now().UnixNano())

	keys := eulerlib.Import("problemdata/keylog.txt")

	eulerlib.BubbleSortVec(keys)

	keys = eulerlib.RemoveDuplicatesVec(keys)

	best := 100

	for {
		current := ""

		for !allwork(keys, current) {
			for i := 0; i < 2*len(keys); i++ {
				try := rand.Int() % len(keys)
				if !works(keys[try], current) {
					current = merge(current, keys[try])
				}
			}

		}

		if len(current) < best {
			best = len(current)
			fmt.Println(current, best)
		}
	}

}
