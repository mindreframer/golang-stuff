package csp

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// Test abstraction for the programs of section 3 that send runes from
// a coroutine west to a coroutine east, possibly transforming the
// string in the process.
func testWestEastProgram(routine func(a, b chan rune), input, expected string, t *testing.T) {
	west, east := make(chan rune), make(chan rune)

	go func() {
		for _, r := range input {
			west <- r
		}
		close(west)
	}()

	go routine(west, east)

	received := make([]rune, 0, 50)
	for r := range east {
		received = append(received, r)
	}
	if string(received) != expected {
		t.Error(string(received))
	}
}

func TestCOPY(t *testing.T) {
	testWestEastProgram(S31_COPY, "Hello ** World***", "Hello ** World***", t)
}

func TestSQUASH_EXT(t *testing.T) {
	testWestEastProgram(S32_SQUASH_EXT, "Hello ** World***", "Hello ↑ World↑*", t)
}

func TestDISASSEMBLE(t *testing.T) {
	cardfile := make(chan []rune)
	X := make(chan rune)

	in := [][]rune{
		[]rune("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		[]rune("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
	}

	expected := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb "

	go func() {
		for _, card := range in {
			cardfile <- card
		}
		close(cardfile)
	}()

	go S33_DISASSEMBLE(cardfile, X)

	actual := make([]rune, 0, 162)
	for r := range X {
		actual = append(actual, r)
	}

	if string(actual) != expected {
		t.Errorf("Got %v", actual)
	}
}

func stringsMustBeEqual(a, b string, t *testing.T) {
	if a != b {
		t.Errorf("Expected '%v' (len %v), got '%v' (len %v)", a, len([]rune(a)), b, len([]rune(b)))
	}
}

func mustBeNLines(lines [][]rune, n int, t *testing.T) {
	if len(lines) != n {
		t.Errorf("Expected %v lines, got %v", n, len(lines))
	}
}

func TestASSEMBLE(t *testing.T) {
	X, lineprinter := make(chan rune), make(chan []rune)

	// 125 a's, 100 b's, the missing 25 at the end should be padded as spaces.
	line1 := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`
	line2 := `bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb`
	in := []rune(line1 + line2)
	go func() {
		for _, r := range in {
			X <- r
		}
		close(X)
	}()

	go S34_ASSEMBLE(X, lineprinter)

	actual := make([][]rune, 0, 2)
	for line := range lineprinter {
		if len(line) == 0 {
			break
		}
		actual = append(actual, line)
	}

	mustBeNLines(actual, 2, t)
	stringsMustBeEqual(string(actual[0]), line1, t)
	stringsMustBeEqual(string(actual[1]), line2+`                         `, t)
}

func TestConway(t *testing.T) {
	cardfile, lineprinter := make(chan []rune), make(chan []rune)

	card1 := `**aa*aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`
	card2 := `b**bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb*b*`

	expected := []rune(strings.Replace(card1+" "+card2+" ", "**", "↑", -1))
	expected1 := string(expected[:125]) // first line
	expected2 := string(expected[125:]) // second line
	// "the last line should be completed with spaces if necessary"
	expected2 += strings.Repeat(" ", 125-len(expected[125:]))

	go func() {
		cardfile <- []rune(card1)
		cardfile <- []rune(card2)
		close(cardfile)
	}()

	go S36_Conway(cardfile, lineprinter)

	actual := make([][]rune, 0, 2)
	for line := range lineprinter {
		if len(line) == 0 {
			break
		}
		actual = append(actual, line)
	}

	mustBeNLines(actual, 2, t)
	stringsMustBeEqual(string(actual[0]), expected1, t)
	stringsMustBeEqual(string(actual[1]), expected2, t)
}

func TestFac(t *testing.T) {
	f := S42_facM(8)
	for _, pair := range [][]int{[]int{0, 1}, []int{1, 1}, []int{4, 24}, []int{8, 40320}} {
		f <- pair[0]
		res := <-f
		if res != pair[1] {
			t.Errorf("Expected %v! == %v, but got %v\n", pair[0], pair[1], res)
		}
	}
}

func TestIntSet(t *testing.T) {
	expect := func(b bool, op string) {
		if !b {
			t.Errorf("%v failed", op)
		}
	}

	set := S43_NewIntSet()
	hasChan := make(chan bool)

	set.Has(0, hasChan)
	expect(!(<-hasChan), "empty !has(0)")
	set.Has(1, hasChan)
	expect(!(<-hasChan), "empty !has(1)")
	set.Has(100, hasChan)
	expect(!(<-hasChan), "empty !has(100)")

	set.Insert(34523, nil)
	set.Has(0, hasChan)
	expect(!(<-hasChan), "{34523} !has(0)")
	set.Has(34523, hasChan)
	expect(<-hasChan, "{34523} has(34523)")

	// Parallel use, Insert() must lock.
	runtime.GOMAXPROCS(runtime.NumCPU())
	n := 10

	ack := make(chan int)
	for i := 0; i < n; i++ {
		go set.Insert(i, ack)
	}
	// Wait until all are inserted.
	for i := 0; i < n; i++ {
		<-ack
	}

	for i := 0; i < n; i++ {
		set.Has(i, hasChan)
		expect(<-hasChan, "parallel insertions")
	}
}

func TestIntSetScan(t *testing.T) {
	set := S43_NewIntSet()

	n := 10
	expected := make([]int, 0, n)
	for i := 23; i < 23+n; i++ {
		expected = append(expected, i)
	}

	ack := make(chan int)
	for _, n := range expected {
		set.Insert(n, ack)
	}
	for i := 0; i < n; i++ {
		<-ack
	}

	actual := make([]int, 0, n)
	elements := set.Scan()
	for el := range elements {
		actual = append(actual, el)
	}
	sort.Ints(actual)

	if len(actual) != len(expected) {
		t.Errorf("Expected %v elements, got %v: %v", len(expected),
			len(actual), actual)
	}
}

func TestParIntSet(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	insert, has, scan, least := S45_ParIntSet(10)

	// The empty set does not contain any number,
	hasResponseChan := make(chan bool)
	for _, num := range []int{0, 10} {
		has <- S45_HasQuery{num, hasResponseChan}
		if <-hasResponseChan {
			t.Errorf("ParSet {} contains %v", num)
		}
	}

	// After insertion, it contains that number but no other ones..
	for _, num := range []int{0, 7, 5} {
		insert <- num

		has <- S45_HasQuery{num, hasResponseChan}
		if h := <-hasResponseChan; !h {
			t.Errorf("ParSet doesn't contain %v", num)
		}
	}
	has <- S45_HasQuery{1, hasResponseChan}
	if <-hasResponseChan {
		t.Errorf("ParSet {0} contains 1")
	}

	// Scan.
	scanRcvr := make(chan int)
	scan <- scanRcvr
	expected := []int{0, 5, 7}
	i := 0
	for n := range scanRcvr {
		if n != expected[i] {
			t.Errorf("Expected %v as %vth number in Scan.", expected[i], i)
		}
		i++
	}

	// Successively ask for the least member and check that it was removed.
	i = 0
	leastResp := make(chan S45_LeastResponse)
	for {
		least <- leastResp
		l := <-leastResp
		if l.NoneLeft {
			break
		}
		if l.Least != expected[i] {
			t.Errorf("Expected %v as %vth least number, got %v.", expected[i], i, l.Least)
		}

		// The least operation is defined as removing the returned
		// element from the set.
		has <- S45_HasQuery{expected[i], hasResponseChan}
		if <-hasResponseChan {
			t.Errorf("Just removed %v from the set, but has says it's still there.", expected[i])
		}

		i++
	}
	if i != 2 {
		t.Errorf("Expected to get %v least members, got %v.", len(expected), i)
	}
}

// 5.1
func TestBuffer(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	consumer, producer := S51_Buffer(10)

	for i := 0; i < 10; i++ {
		select {
		case producer <- i:
			// empty
		case <-time.After(2 * time.Second):
			t.Errorf("Producer couldn't send %vth value", i)
			return
		}
	}
	// We should get here without blocking since the buffer stores the
	// ten portions.

	received := make([]int, 10)
	for i := 0; i < 10; i++ {
		select {
		case received[i] = <-consumer:
			if received[i] != i {
				t.Errorf("Received %v as %vth number.", received[i], i)
			}
		case <-time.After(2 * time.Second):
			t.Errorf("Consumer couldn't receive %vth value", i)
			return
		}
	}
}

func TestSemaphore(t *testing.T) {
	s := S52_NewSemaphore()

	// We cannot decrement before an increment.
	select {
	case s.dec <- struct{}{}:
		t.Error("Shouldn't be able to dec before inc")
	case <-time.After(2 * time.Second):
		// ok, dec blocked
	}

	// After an inc, dec will work.
	s.inc <- struct{}{}
	select {
	case s.dec <- struct{}{}:
		// ok
	case <-time.After(2 * time.Second):
		t.Error("Should be able to dec after inc")
	}
}

// That's not actually a test, we just let the scenario run for 10
// seconds so we can observe the log.
func TestDiningPhilosophers(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	S53_DiningPhilosophers(10 * time.Second)
}

func TestPrimeSieve(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	numPrimes := 100

	// Copied from http://primes.utm.edu/lists/small/10000.txt
	first100PrimesStr := "2 3 5 7 11 13 17 19 23 29 31 37 41 43 47 53 59 61 67 71 73 79 83 89 97 101 103 107 109 113 127 131 137 139 149 151 157 163 167 173 179 181 191 193 197 199 211 223 227 229 233 239 241 251 257 263 269 271 277 281 283 293 307 311 313 317 331 337 347 349 353 359 367 373 379 383 389 397 401 409 419 421 431 433 439 443 449 457 461 463 467 479 487 491 499 503 509 521 523 541"
	first100Primes := strings.Split(first100PrimesStr, " ")

	primes := make([]int, 0, numPrimes)
	primeChan := make(chan int)

	doneChan := make(chan bool)
	go func() {
		for p := range primeChan {
			if p == -1 {
				doneChan <- true
				return
			}
			primes = append(primes, p)
		}
	}()

	S61_SIEVE(numPrimes, primeChan)
	<-doneChan

	l := len(first100Primes)
	if len(primes) != l {
		t.Errorf("Expected %v primes, but got %v.", l, len(primes))
	}

	// As SIEVE runs concurrently on possibly multiple cores, the
	// primes can arrive out of order.
	sort.Ints(primes)

	for i := 0; i < l; i++ {
		if strconv.Itoa(primes[i]) != first100Primes[i] {
			t.Errorf("Expected %v as %vnth prime, got %v.", first100Primes[i], i, primes[i])
		}
	}
}

func TestMatrixMultiply(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	A := [][]float64{
		[]float64{1, 2, 3},
		[]float64{4, 5, 6},
		[]float64{7, 8, 9},
	}
	matrix := S62_NewMatrix(A)

	for _, testcase := range []struct{ other, expected [][]float64 }{
		{
			other: [][]float64{
				[]float64{1, 1, 1},
				[]float64{1, 1, 1},
				[]float64{1, 1, 1},
			},
			expected: [][]float64{
				[]float64{12, 15, 18},
				[]float64{12, 15, 18},
				[]float64{12, 15, 18},
			},
		},
		{
			other: [][]float64{
				[]float64{1, 1, 1},
				[]float64{2, 2, 2},
				[]float64{3, 3, 3},
			},
			expected: [][]float64{
				[]float64{12, 15, 18},
				[]float64{24, 30, 36},
				[]float64{36, 45, 54},
			},
		},
		{
			other: [][]float64{
				[]float64{1, 2, 3},
				[]float64{1, 2, 3},
				[]float64{1, 2, 3},
			},
			expected: [][]float64{
				[]float64{30, 36, 42},
				[]float64{30, 36, 42},
				[]float64{30, 36, 42},
			},
		},
	} {
		result := make([][]float64, 3)
		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			result[i] = make([]float64, 3)
			wg.Add(1)
			go func(i int) {
				for j := 0; j < 3; j++ {
					val := <-matrix.SOUTH[i]
					result[j][i] = val
				}
				wg.Done()
			}(i)
		}

		for i := 0; i < 3; i++ {
			go func(i int) {
				for j := 0; j < 3; j++ {
					matrix.WEST[j] <- testcase.other[i][j]
				}
			}(i)
		}

		wg.Wait()

		if !matricesEqual(result, testcase.expected) {
			t.Errorf("Expected \n%v, got \n%v",
				printMatrix(testcase.expected), printMatrix(result))
		}
	}
}

func matricesEqual(a, b [][]float64) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func printMatrix(m [][]float64) string {
	var b bytes.Buffer
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[0]); j++ {
			fmt.Fprintf(&b, "%v ", m[i][j])
		}
		fmt.Fprintf(&b, "\n")
	}

	return b.String()
}
