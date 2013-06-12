package main

import (
	"./euler"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type card struct {
	value int
	suit  string
}

func makeCard(letters string) card {
	value, err := strconv.Atoi(letters[:1])
	if err != nil {
		if letters[:1] == "T" {
			value = 10
		}
		if letters[:1] == "J" {
			value = 11
		}
		if letters[:1] == "Q" {
			value = 12
		}
		if letters[:1] == "K" {
			value = 13
		}
		if letters[:1] == "A" {
			value = 14
		}
	}

	return card{value, letters[1:]}
}

func hasPair(hand []card) int {
	for i, card := range hand {
		for _, match := range hand[i+1:] {
			if match.value == card.value {
				return card.value
			}
		}
	}
	return 0
}

func removeCard(value int, hand []card) []card {
	handcopy := make([]card, len(hand))
	copy(handcopy, hand)
	for i, card := range handcopy {
		if card.value == value {
			return append(handcopy[:i], handcopy[i+1:]...)
		}
	}
	return handcopy
}

func has3(hand []card) int {
	for i, card1 := range hand {
		for j, card2 := range hand[i+1:] {
			for _, card3 := range hand[i+j+2:] {
				if card1.value == card2.value && card2.value == card3.value {
					return card1.value
				}
			}
		}
	}
	return 0
}

func hasHouse(hand []card) (three int, two int) {
	three = has3(hand)
	if three == 0 {
		return 0, 0
	}

	a, b := has2pair(hand)

	if a == 0 {
		return 0, 0
	}
	if a == three {
		return three, b
	}
	if b == three {
		return three, a
	}
	return 0, 0 //four of a kind

}

func hasflush(hand []card) bool {
	suit := hand[0].suit
	for _, card := range hand {
		if suit != card.suit {
			return false
		}
	}
	return true

}

func has2pair(hand []card) (pair1 int, pair2 int) {
	pair1 = hasPair(hand)
	if pair1 == 0 {
		return 0, 0
	}
	pair2 = hasPair(removeCard(pair1, removeCard(pair1, hand)))
	if pair2 == 0 {
		return 0, 0
	}
	return
}

func hasStraight(hand []card) int {
	valuelist := make([]int, 0)
	for _, card := range hand {
		valuelist = append(valuelist, card.value)
	}
	valuelist = euler.BubbleSortInts(valuelist)
	for i := 0; i < len(valuelist)-1; i++ {
		if valuelist[i]-1 != valuelist[i+1] {
			return 0
		}
	}

	return valuelist[0]
}

func highCard(hand []card) int {
	valuelist := make([]int, 0)
	for _, card := range hand {
		valuelist = append(valuelist, card.value)
	}

	valuelist = euler.BubbleSortInts(valuelist)

	return valuelist[0]
}

func wins(hand1, hand2 []card) bool {

	//Someone has a straight flush
	if (hasflush(hand1) && hasStraight(hand1) != 0) ||
		(hasflush(hand2) && hasStraight(hand2) != 0) {
		fmt.Println("Straight flush, what do?!", hand1, hand2)
		//This case isn't in the file, so I didn't bother 
	}

	a, b := has2pair(hand1)
	c, d := has2pair(hand2)

	//someone has four of a kind
	if (a != 0 && a == b) || (c != 0 && d == c) {
		if a != 0 && c == 0 {
			return true
		}
		if c != 0 && a == 0 {
			return false
		} else {
			fmt.Println("4 of a kind, what do?", hand1, hand2)
			//Again, this case doesn't appear
		}
	}

	a, b = hasHouse(hand1)
	c, d = hasHouse(hand2)

	//someone has full house
	if a+b+c+d != 0 {
		if a != 0 && c == 0 {
			return true
		}
		if c != 0 && a == 0 {
			return false
		} else {
			fmt.Println("Two FH, what do?", hand1, hand2)
		}
	}

	//Flushes
	if hasflush(hand1) && !hasflush(hand2) {
		return true
	}
	if hasflush(hand2) && !hasflush(hand1) {
		return false
	}
	if hasflush(hand1) && hasflush(hand2) {
		fmt.Println("Two flushes!!", hand1, hand2)
	}

	//Straights
	if hasStraight(hand1)+hasStraight(hand2) > 0 {
		return hasStraight(hand1) > hasStraight(hand2)
	}

	//Three of a kind
	if has3(hand1)+has3(hand2) != 0 {
		return has3(hand1) > has3(hand2)
	}

	a, b = has2pair(hand1)
	c, d = has2pair(hand2)

	//Two pair
	if a+b+c+d != 0 { //someone has two pair 

		if a != 0 && c == 0 {
			return true
		}
		if c != 0 && a == 0 {
			return false
		} else {
			fmt.Println("Two pair, what do?", hand1, hand2)
			//Again, this case doesn't appear
		}

	}

	//One pair
	if hasPair(hand1)+hasPair(hand2) != 0 {
		return hasPair(hand1) > hasPair(hand2)
	}

	//High card
	return highCard(hand1) > highCard(hand2)
}

func main() {
	starttime := time.Now()

	data := euler.Import("problemdata/poker.txt")

	total := 0

	for _, line := range data {

		hand1 := make([]card, 5)
		for i, card := range strings.Split(line[:14], " ") {
			hand1[i] = makeCard(card)
		}

		hand2 := make([]card, 5)
		for i, card := range strings.Split(line[15:], " ") {
			hand2[i] = makeCard(card)
		}

		if wins(hand1, hand2) {
			total++
		}

	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
