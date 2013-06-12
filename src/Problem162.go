package main

import (
	"fmt"
	"math/big"
	"time"
)

var memo map[[5]int]*big.Rat

func prob(zero, one, a, leading bool, digits int) *big.Rat {

	if answer, ok := memo[keygen(zero, one, a, leading, digits)]; ok {
		return answer
	}

	if digits == 0 {
		if zero && one && a {
			return big.NewRat(1, 1)
		} else {
			return big.NewRat(0, 1)
		}
	}

	answer := big.NewRat(0, 1)

	uno := big.NewRat(1, 16)
	if leading == true {
		uno = uno.Mul(uno, prob(true, one, a, leading, digits-1))

	} else {
		uno = uno.Mul(uno, prob(zero, one, a, leading, digits-1))
	}
	answer = answer.Add(answer, uno)

	dos := big.NewRat(1, 16)
	dos = dos.Mul(dos, prob(zero, true, a, true, digits-1))
	answer = answer.Add(answer, dos)

	tres := big.NewRat(1, 16)
	tres = tres.Mul(tres, prob(zero, one, true, true, digits-1))
	answer = answer.Add(answer, tres)

	rest := big.NewRat(13, 16)
	rest = rest.Mul(rest, prob(zero, one, a, true, digits-1))
	answer = answer.Add(answer, rest)

	memo[keygen(zero, one, a, leading, digits)] = answer

	return answer

}

func keygen(a, b, c, d bool, n int) [5]int {
	var A, B, C, D int
	if a {
		A = 1
	} else {
		A = 0
	}
	if b {
		B = 1
	} else {
		B = 0
	}
	if c {
		C = 1
	} else {
		C = 0
	}
	if d {
		D = 1
	} else {
		D = 0
	}
	return [5]int{A, B, C, D, n}

}

func main() {
	starttime := time.Now()

	memo = make(map[[5]int]*big.Rat)

	all := big.NewRat(16, 1)
	all = all.Mul(all, all) //16^2
	all = all.Mul(all, all) //16^4
	all = all.Mul(all, all) //16^8
	all = all.Mul(all, all) //16^16

	ours := all.Mul(all, prob(false, false, false, false, 16)).Num()

	fmt.Println(ours.String(), "make sure to conver to hex!")

	fmt.Println("Elapsed time:", time.Since(starttime))
}
