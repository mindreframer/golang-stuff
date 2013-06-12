// The examples from Tony Hoare's seminal 1978 paper "Communicating
// sequential processes" implemented in Go.
//
// Go's design was strongly influenced by Hoare's paper [1]. Although
// Go differs significantly from the example language used in the
// paper, the examples still translate rather easily. The biggest
// difference apart from syntax is that Go models the conduits of
// concurrent communication explicitly as channels, while the
// processes of Hoare's language send messages directly to each other,
// similar to Erlang. Hoare hints at this possibility in section 7.3,
// but with the limitation that "each port is connected to exactly one
// other port in another process", in which case it would be a mostly
// syntactic difference.
//
// [1]
// http://blog.golang.org/2010/07/share-memory-by-communicating.html
//
// Implementing these examples, and the careful reading of the paper
// required to do so, were a very enlightening experience. I found the
// iterative array of 4.2, the concurrent routines changing their
// behavior execution of 4.5, and the highly concurrent matrix
// multiplication of 6.2 to be particularly interesting.
//
// I tried to name routines and variables like in the paper, which
// explains the now outdated upper-case names. Similarly, I tried to
// make the function signatures as similar to the paper as possible,
// so we mostly work directly with channels, where one would hide this
// implementation detail in real-world code.
//
// The ugly name prefixes like "S33_" make godoc order the types by
// section of the paper.
//
// Most of the examples have tests, although I have not taken a lot of
// care to test corner cases. The test of the S53_DiningPhilosophers
// is not really a test, it simply runs the routine for ten seconds so
// you can observe the philosophers behavior (use `go test -v`).
//
// Thomas Kappler <tkappler@gmail.com>
package csp

import (
	"fmt"
	"sync"
	"time"
)

// 3.1 COPY
//
// > "Problem: Write a process X to copy characters output by process west
// to process, east."
//
// In Go, the communication channel between two processes (goroutines) is
// explicitly represented via the chan type. So we model west and east as
// channels of runes. In an endless loop, we read from west and write the
// result directly to east.
//
// As an addition to the paper's example, we stop when west is closed,
// otherwise we would just hang at this point. To indicate this to the
// client, we close the east channel.
func S31_COPY(west, east chan rune) {
	for r := range west {
		east <- r
	}
	close(east)
}

// 3.2 SQUASH
//
// > "Problem: Adapt the previous program [COPY] to replace every pair of
// consecutive asterisks "**" by an upward arrow "↑". Assume that the final
// character input is not an asterisk."
//
// If we get an asterisk from west, we receive the next character as well
// and then decide whether to send the upward arrow or the last two
// characters as-is. Go's UTF8 support allows to treat the arrow like any
// other character.
func S32_SQUASH(west, east chan rune) {
	for c := range west {
		if c != '*' {
			east <- c
		} else {
			c2 := <-west
			if c2 != '*' {
				east <- c
				east <- c2
			} else {
				east <- '↑'
			}
		}
	}
	close(east)
}

// Hoare adds a remark to 3.2 SQUASH: "(2) As an exercise, adapt this
// process to deal sensibly with input which ends with an odd number of
// asterisks." This version handles this case by sending a single asterisk
// from west to east if west did not supply another character after a
// timeout, or if west was closed in the meantime.
func S32_SQUASH_EXT(west, east chan rune) {
	for c := range west {
		if c != '*' {
			east <- c
		} else {
			select {
			case <-time.After(10 * time.Second):
				east <- c
			case c2, ok := <-west:
				if !ok { // west closed
					east <- c
					break
				}
				if c2 != '*' {
					east <- c
					east <- c2
				} else {
					east <- '↑'
				}

			}
		}
	}
	close(east)
}

// 3.3 DISASSEMBLE
//
// > "Problem: to read cards from a cardfile and output to process X the
// stream of characters they contain. An extra space should be inserted at
// the end of each card."
//
// Trivially translated to Go. We don't need to care about the indices 1 to
// 80, range handles this for us.
func S33_DISASSEMBLE(cardfile chan []rune, X chan rune) {
	for cardimage := range cardfile {
		for _, r := range cardimage {
			X <- r
		}
		X <- ' '
	}
	close(X)
}

// 3.4 ASSEMBLE
//
// > "Problem: To read a stream of characters from process X and print them
// in lines of 125 characters on a lineprinter. The last line should be
// completed with spaces if necessary."
func S34_ASSEMBLE(X chan rune, lineprinter chan []rune) {
	linelen := 125
	lineimage := make([]rune, linelen)
	i := 0
	for c := range X {
		if c == 0 {
			break
		}
		lineimage[i] = c
		i++
		if i == linelen {
			c := make([]rune, linelen)
			copy(c, lineimage)
			lineprinter <- c
			i = 0
		}
	}

	// Print the last line padded with spaces.
	if i > 0 {
		for j := i; j < linelen; j++ {
			lineimage[j] = ' '
		}
		lineprinter <- lineimage
	}

	lineprinter <- nil
}

// 3.5 Reformat
//
// > "Problem: Read a sequence of cards of 80 characters each, and print
// the characters on a lineprinter at 125 characters per line. Every card
// should be followed by an extra space, and the last line should be
// completed with spaces if necessary."
//
// This is a great example of how easily concurrent processes can be
// combined in the manner Unix pipes. No extra code is required to let the
// data flow through the two routines we wrote earlier.
func S35_Reformat(cardfile, lineprinter chan []rune) {
	pipe := make(chan rune)
	go S33_DISASSEMBLE(cardfile, pipe)
	S34_ASSEMBLE(pipe, lineprinter)
}

// 3.6 Conway's Problem
//
// > "Problem: Adapt the above program to replace every pair of consecutive
// asterisks by an upward arrow."
//
// The implementation in four lines is a testament to the expressive
// power of modeling programs as communicating sequential processes.
func S36_Conway(cardfile, lineprinter chan []rune) {
	pipe1, pipe2 := make(chan rune), make(chan rune)
	go S33_DISASSEMBLE(cardfile, pipe1)
	go S32_SQUASH_EXT(pipe1, pipe2)
	S34_ASSEMBLE(pipe2, lineprinter)
}

// 4. Subroutines and Data Representations
//
// > "A coroutine acting as a subroutine is a process operating
// concurrently with its user process in a parallel command:
// [subr::SUBROUTINE||X::USER]. [...] The USER will call the subroutine by
// a pair of commands: subr!(arguments); ...; subr?(results). Any commands
// between these two will be executed concurrently with the subroutine."
//
// Here the paper's influence on Go comes out clearly: coroutines are
// goroutines and launching them via "!" is the "go" command. Only reading
// the results is quite different in Go because of the explicit
// representation of the conduit between coroutine and main routine, the
// channel.

// 4.1  Function: Division With Remainder
//
// > "Problem: Construct a process to represent a function-type subroutine,
// which accepts a positive dividend and divisor, and returns their integer
// quotient and remainder. Efficiency is of no concern."
func S41_DIV(x, y int, res chan struct{ quot, rem int }) {
	quot := 0
	rem := x
	for rem >= y {
		rem -= y
		quot += 1
	}
	res <- struct{ quot, rem int }{quot, rem}
}

// 4.2 Recursion: Factorial
//
// > "Problem: Compute a factorial by the recursive method, to a given
// limit."
//
// This example is fascinating. It introduces the "iterative array" which
// kept me puzzled for a while, but made for a great a-ha moment when I got
// it. It's an array of coroutines, so that for a given integer input i,
// the coroutine at index i knows how to deal with it. By addressing its
// neighbor coroutines at i-1 and i+1, a coroutine can communicate with the
// others to break down the overall problem. That sounds pretty abstract
// and will be clearer in the following examples.
//
// To compute n!, we use the simple recurrence `n! = n * (n-1)!`, with `0!
// = 1! = 1`. In our iterative array of goroutines, when routine i receives
// the value x, it sends x-1 up the chain, i.e., to the right in the
// iterative array. This continues until the value is 0 or 1 and we hit the
// base case of the recursion. The coroutines then pass values back down
// the chain, i.e. leftwards, starting with 1. When routine i receives a
// result passed back down the chain, it multiplies it with x and passes on
// the result. When it arrives at routine 0, n! is computed.
//
// The caller doesn't see this process. We only need to expose goroutine 0,
// which can compute any factorial by passing the value up the chain and
// waiting for the result. Any factorial up to the limit of the length of
// the iterative array, that is, which has to be given when creating the
// factorial iterative array.
//
// Go models communication between goroutines with explicit channel values,
// so in this implementation the iterative array is an array of channels.
// When we create it, we launch the corresponding goroutines at the same
// time.
func S42_facM(limit int) chan int {
	fac := make([]chan int, limit+1, limit+1)

	fac[0] = make(chan int)

	for i := 1; i <= limit; i++ {
		fac[i] = make(chan int)
		go func(i int) {
			for {
				n := <-fac[i-1]
				if n == 0 || n == 1 {
					fac[i-1] <- 1
				} else {
					fac[i] <- n - 1
					r := <-fac[i]
					fac[i-1] <- n * r
				}
			}
		}(i)
	}

	return fac[0]
}

// 4.3 Data Representation: Small Set of Integers
//
// > "Problem: To represent a set of not more than 100 integers as a
// process, S, which accepts two kinds of instruction from its calling
// process X: (1) S!insert(n), insert the integer n in the set, and (2)
// S!has(n); ... ; S?b, b is set true if n is in the set, and false
// otherwise."
type S43_IntSet struct {
	content   []int
	writeLock sync.RWMutex
}

// If s contains n, return its index, otherwise return the next free index.
func (s *S43_IntSet) search(n int) int {
	for i, el := range s.content {
		if el == n {
			return i
		}
	}
	return len(s.content)
}

func S43_NewIntSet() S43_IntSet {
	return S43_IntSet{content: make([]int, 0, 100)}
}

// Send true on res if s contains n, otherwise false.
//
// The caller needs to pass in the result channel because otherwise we'd
// need to return a channel here to make the operation asynchronous. But
// what channel? If we make a new one every time, it's wasteful. If we have
// only one, we need to lock access to it, and we cannot close it, so the
// caller might wait indefinitely on it (although that would be an error in
// the client).
func (s *S43_IntSet) Has(n int, res chan<- bool) {
	go func() {
		s.writeLock.RLock()
		i := s.search(n)
		res <- (i < len(s.content))
		s.writeLock.RUnlock()
	}()
}

// Insert the number n into the set, if there is still room. Note that in
// Hoare's specification the client has no way of knowing whether the set
// is full or not, safe for testing whether the insertion worked with
// has().
//
// The client can also not know when the insertion is complete. Parallel
// insertions and Has() queries are protected by a mutex. But there is no
// guarantee that the Insert() has even started to run, i.e., actually
// acquired the lock. To protect against seeing stale data, the client can
// pass in an ack channel and block on it.
func (s *S43_IntSet) Insert(n int, ack chan int) {
	go func() {
		s.writeLock.Lock()
		i := s.search(n)
		size := len(s.content)
		// If i is < size, n is already in the set, see Has().
		if i == size && size <= 100 {
			s.content = append(s.content, n)
		}
		s.writeLock.Unlock()
		if ack != nil {
			ack <- 1
		}
	}()
}

// 4.4 Scanning a Set
//
// > "Problem: Extend the solution to 4.3 by providing a fast method for
// scanning all members of the set without changing the value of the set."
//
// The implementation below looks quite different from Hoare's. Go's
// channels in combination with `range` make the implementation trivial. An
// implementation closer to the pseudocode in the paper might return a chan
// int and a chan bool "noneleft" to the caller, sending a signal on
// noneleft after the iteration.
func (s *S43_IntSet) Scan() chan int {
	res := make(chan int)
	go func() {
		s.writeLock.RLock()
		for _, c := range s.content {
			res <- c
		}
		s.writeLock.RUnlock()
		close(res)
	}()
	return res
}

// 4.5 Recursive Data Representation: Small Set of Integers
//
// > "Problem: Same as above, but an array of processes is to be used to
// achieve a high degree of parallelism. Each process should contain at
// most one number. When it contains no number, it should answer "false" to
// all inquiries about membership. On the first insertion, it changes to a
// second phase of behavior, in which it deals with instructions from its
// predecessor, passing some of them on to its successor. The calling
// process will be named S(0). For efficiency, the set should be sorted,
// i.e. the ith process should contain the ith largest number."
//
// We use the iterative array technique from 4.2 again here.
//
// I found this exercise to be the trickiest one. The *least*
// operation had me puzzled for a while. I tried to make it work,
// analogous to the other three, with a single channel used both to
// communicate with the client and to communicate internally between
// the goroutines.
func S45_ParIntSet(limit int) (chan int, chan S45_HasQuery, chan chan int, chan S45_LeastQuery) {
	insert := make([]chan int, limit+1)
	has := make([]chan S45_HasQuery, limit+1)
	scan := make([]chan chan int, limit+1)
	least := make([]chan S45_LeastResponse, limit+1)
	// Only the first one of these will actually be created and used to
	// communicate with the client, but we make an array to be able to
	// use the same code for all goroutines.
	leastQuery := make([]chan S45_LeastQuery, limit+1)

	for i := 1; i <= limit; i++ {
		insert[i] = make(chan int)
		has[i] = make(chan S45_HasQuery)
		scan[i] = make(chan chan int)
		least[i] = make(chan S45_LeastResponse)
		if i == 1 {
			leastQuery[i] = make(chan S45_LeastQuery)
		}

		go func(i int) {
			// This goroutine stores n.
			var n int

		EMPTY:
			for {
				select {
				case q := <-has[i]:
					q.Response <- false
				case n = <-insert[i]:
					break EMPTY
				case c := <-scan[i]:
					close(c)
				case least[i] <- S45_LeastResponse{NoneLeft: true}:
					// (empty)
				case q := <-leastQuery[i]:
					q <- S45_LeastResponse{NoneLeft: true}
				}
			}

			// NONEMPTY:
			for {
				select {
				case q := <-has[i]:
					if q.N <= n {
						q.Response <- q.N == n
					} else {
						if i == limit {
							// We've reached the limit of the set.
							q.Response <- false
						} else {
							// We don't have q, pass on the request.
							has[i+1] <- q
						}
					}
				case m := <-insert[i]:
					// If m is larger than our number n, pass it on,
					// otherwise pass ours on and keep m.
					if m < n {
						if i < limit {
							insert[i+1] <- n
						}
						n = m
					} else if m > n && i < limit {
						insert[i+1] <- m
					}
				case c := <-scan[i]:
					// Send our value and pass on the response channel.
					c <- n
					scan[i+1] <- c
				case least[i] <- S45_LeastResponse{n, false}:
					// Get the least from the next goroutine.
					nextL := <-least[i+1]

					// Shift one to the left in our concurrent list.
					if nextL.NoneLeft {
						goto EMPTY
					} else {
						n = nextL.Least
					}
				case leastQ := <-leastQuery[i]:
					// This case is only for the client-facing goroutine 1.

					// Get the least from the next goroutine.
					nextL := <-least[i+1]

					// Send our number to the client, and the reply from
					// the next goroutine whether there are more to come.
					leastQ <- S45_LeastResponse{n, nextL.NoneLeft}

					// Shift one to the left in our concurrent list.
					if nextL.NoneLeft {
						goto EMPTY
					} else {
						n = nextL.Least
					}
				}
			}
		}(i)
	}

	return insert[1], has[1], scan[1], leastQuery[1]
}

type S45_HasQuery struct {
	N        int
	Response chan bool
}

type S45_LeastQuery chan S45_LeastResponse

type S45_LeastResponse struct {
	Least    int
	NoneLeft bool
}

// 5.1 Bounded Buffer
//
// > "Problem: Construct a buffering process X to smooth variations in the
// speed of output of portions by a producer process and input by a
// consumer process. The consumer contains pairs of commands X!more( );
// X?p, and the producer contains commands of the form X!p. The buffer
// should contain up to ten portions."
//
// This is exactly what Go's buffered channels provide. We will do it
// manually here. Even that would be trivial using a `select` containing
// both the producer receive and the consumer send, so we'll do it without
// select. The implementation strictly follows Hoare's pseudo-code from the
// paper.
//
// We do actually use `select`, but only for single channel operations, in
// order to exploit the semantics of its `default` case: try the channel
// operation, but don't block if the other end isn't ready, instead just
// skip it and continue.
func S51_Buffer(bufSize int) (consumer chan int, producer chan int) {
	buffer := make([]int, bufSize)
	consumer = make(chan int)
	producer = make(chan int)

	in, out := 0, 0
	go func() {
		for {
			if in < out+10 {
				// We have room in the buffer, check the producer.
				select {
				case i := <-producer:
					buffer[in%bufSize] = i
					in++
				default: // don't block
				}
			}

			if out < in {
				// We have something in the buffer, check the consumer.
				select {
				case consumer <- buffer[out%bufSize]:
					out++
				default: // don't block
				}
			}

			// Sleep for a bit here to avoid busy waiting?
		}
	}()

	return consumer, producer
}

// 5.2 Integer Semaphore
//
// > "Problem: To implement an integer semaphore, S, shared among an array
// X(i:I..100) of client processes. Each process may increment the
// semaphore by S!V() or decrement it by S!P(), but the latter command must
// be delayed if the value of the semaphore is not positive."
//
// This is a nice one to write using `select`. We use two channels, inc and
// dec, for the two operations the semaphore offers. If dec isn't possible
// because the semaphore is 0, the client's channel send blocks.
//
// The number of clients, 100 in the paper, doesn't matter for the inc and
// dec operations since, in contrast to Hoare's pseudocode, we don't need
// to explicitly scan all clients in the channel receive operations. We
// would need to know the number if we wanted to keep track of active
// clients and shut down the semaphore once they are all finished. I didn't
// implement this, but it would be easy to do with an `activeClients int`
// variable and a `done` channel that decrements it.
//
// A problem I faced initially was that I put both the inc and the dec
// receive operations into one select. This would require a guard on the
// dec receive, something like `case val > 0 && <- dec:`, but Go doesn't
// support that. So I would receive from dec, thereby acknowledging the
// operation to the caller, before knowing that the decrement was legal.
// The solution is obvious in retrospect: try the dec receive in a separate
// select, protected by a `val > 0` guard.
type S52_Semaphore struct {
	inc, dec chan struct{}
}

func S52_NewSemaphore() *S52_Semaphore {
	s := &S52_Semaphore{
		inc: make(chan struct{}),
		dec: make(chan struct{}),
	}

	val := 0

	go func() {
		// We need at least one increment before we can react to dec.
		<-s.inc
		val++

		for {
			select {
			case <-s.inc:
				val++
			case <-s.dec:
				val--
				// If val is 0, we need an inc before we can continue.
				if val == 0 {
					<-s.inc
					val++
				}
			}
		}
	}()

	return s
}

func (s *S52_Semaphore) Inc() {
	s.inc <- struct{}{}
}

func (s *S52_Semaphore) Dec() {
	s.dec <- struct{}{}
}

// 5.3 Dining Philosophers (Problem due to E.W. Dijkstra)
//
// > "Problem: Five philosophers spend their lives thinking and eating. The
// philosophers share a common dining room where there is a circular table
// surrounded by five chairs, each belonging to one philosopher. In the
// center of the table there is a large bowl of spaghetti, and the table is
// laid with five forks (see Figure 1). On feeling hungry, a philosopher
// enters the dining room, sits in his own chair, and picks up the fork on
// the left of his place. Unfortunately, the spaghetti is so tangled that
// he needs to pick up and use the fork on his right as well. When he has
// finished, he puts down both forks, and leaves the room. The room should
// keep a count of the number of philosophers in it."
//
// The dining philosophers are famous in Computer Science because they
// illustrate the problem of deadlock. As Hoare explains, "The solution
// given above does not prevent all five philosophers from entering the
// room, each picking up his left fork, and starving to death because he
// cannot pick up his right fork."
func S53_DiningPhilosophers(runFor time.Duration) {
	// The room is a goroutine that listens on a channel to signal "enter"
	// and one to signal "exit".
	enterRoom := make(chan int)
	exitRoom := make(chan int)
	room := func() {
		occupancy := 0
		for {
			select {
			case i := <-enterRoom:
				if occupancy < 4 {
					occupancy++
				} else {
					// If all philosophers sit down to eat, they starve.
					// Wait for someone to leave.
					fmt.Printf("%v wants to enter, but must wait.\n", i)
					<-exitRoom
					// Enter the room, occupancy stays the same.
					fmt.Printf("%v can finally enter!\n", i)
				}
			case <-exitRoom:
				occupancy--
			}
		}
	}

	// The forks are goroutines listening to pickup and putdown channels
	// like the room, but we need one channel per philosopher to
	// distinguish them so that we can match pickup and putdown of a fork.
	pickup := make([]chan int, 5)
	putdown := make([]chan int, 5)
	for i := 0; i < 5; i++ {
		pickup[i] = make(chan int)
		putdown[i] = make(chan int)
	}
	fork := func(i int) {
		for {
			select {
			case <-pickup[i]:
				<-putdown[i]
			case <-pickup[abs(i-1)%5]:
				<-putdown[abs(i-1)%5]
			}
		}
	}

	// Thinking and eating are sleeps followed by a log message so we know
	// what's going on.
	think := func(i int) {
		time.Sleep(2 * time.Second)
		fmt.Printf("%v thought.\n", i)
	}
	eat := func(i int) {
		time.Sleep(1 * time.Second)
		fmt.Printf("%v ate.\n", i)
	}

	// A philosopher leads a simple life.
	philosopher := func(i int) {
		for {
			think(i)
			enterRoom <- i
			pickup[i] <- i
			pickup[(i+1)%5] <- i
			eat(i)
			putdown[i] <- i
			putdown[(i+1)%5] <- i
			exitRoom <- i
		}
	}

	// Launch the scenario.
	go room()
	for i := 0; i < 5; i++ {
		go fork(i)
	}
	for i := 0; i < 5; i++ {
		go philosopher(i)
	}

	time.Sleep(runFor)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 6.1 Prime Numbers: The Sieve of Eratosthenes
//
// > "Problem: To print in ascending order all primes less than 10000. Use
// an array of processes, SIEVE, in which each process inputs a prime from
// its predecessor and prints it. The process then inputs an ascending
// stream of numbers from its predecessor and passes them on to its
// successor, suppressing any that are multiples of the original prime."
//
// Here I ran into a problem that I suspect is with the algorithm itself.
// It uses one process aka goroutine per prime number. The pseudocode in
// the paper instantiates 101 processes, enough for the first 100 primes.
// However, it sends the numbers up to 10.000 to the first process, which
// contain the first 1229 primes. Maybe I overlooked something about how
// the paper's pseudo-implementation handles this, but in my Go
// implementation I get a deadlock after these first 100 primes.
//
// Not wanting to spend too much time on this, I changed the function
// signature to accept a numPrimes parameter giving the number of
// primes to generate.
func S61_SIEVE(numPrimes int, primes chan int) {
	sieve := make([]chan int, numPrimes)
	sieve[numPrimes-1] = make(chan int)
	done := make(chan bool)

	for i := 0; i < numPrimes-1; i++ {
		sieve[i] = make(chan int)

		go func(i int) {
			p, ok := <-sieve[i]
			if !ok {
				return
			}

			primes <- p

			mp := p // mp is a multiple of p
			for m := range sieve[i] {
				for m > mp {
					mp += p
				}
				if m < mp {
					sieve[i+1] <- m
				}
			}
		}(i)
	}

	go func() {
		p := <-sieve[numPrimes-1]
		primes <- p
		done <- true
	}()

	// Send 2, then all odd numbers up to upto.
	sieve[0] <- 2
	n := 3
SENDNUMBERS:
	for {
		select {
		case sieve[1] <- n: // empty
		case <-done:
			break SENDNUMBERS
		}
		n += 2
	}

	primes <- -1
}

// A matrix for use in example 6.2, matrix multiplication.
type S62_Matrix struct {
	A                   [][]float64
	WEST, SOUTH         []chan float64
	eastward, southward [][]chan float64
}

// A constant source of zeros from the top.
func (m S62_Matrix) NORTH(col int) {
	for {
		m.southward[0][col] <- 0.0
	}
}

// A sink on the right.
func (m S62_Matrix) EAST(row int) {
	rightmost := len(m.eastward[row]) - 1
	for _ = range m.eastward[row][rightmost] {
		// do nothing, just consume
	}
}

// A concurrent routine for a matrix cell that's not on an edge.
func (m S62_Matrix) CENTER(row, col int) {
	for x := range m.eastward[row][col-1] {
		m.eastward[row][col] <- x
		sum := <-m.southward[row-1][col]
		m.southward[row][col] <- (m.A[row-1][col-1]*x + sum)
	}
}

// 6.2 An Iterative Array: Matrix Multiplication
//
// > "Problem: A square matrix A of order 3 is given. Three streams
// are to be input, each stream representing a column of an array IN.
// Three streams are to be output, each representing a column of" the
// product matrix IN × A."
//
// Make a new matrix with the given values (rows, then columns). This
// constructor will launch the goroutines for NORTH, EAST and CENTER
// as described in the paper. The client can then send the values of
// row i of IN to WEST[i] and read column j of IN × A from SOUTH[j].
// See the test for an example.
func S62_NewMatrix(values [][]float64) S62_Matrix {
	numRows := len(values)
	numCols := len(values[0])

	m := S62_Matrix{A: values}

	m.eastward = make([][]chan float64, numRows+1)
	for i := 0; i < numRows+1; i++ {
		m.eastward[i] = make([]chan float64, numCols+1)
		for j := 0; j < numCols+1; j++ {
			m.eastward[i][j] = make(chan float64)
		}
	}

	m.southward = make([][]chan float64, numRows+1)
	for i := 0; i < numRows+1; i++ {
		m.southward[i] = make([]chan float64, numCols+1)
		for j := 0; j < numCols+1; j++ {
			m.southward[i][j] = make(chan float64)
		}
	}

	m.WEST = make([]chan float64, numRows)
	for row := 1; row <= numRows; row++ {
		m.WEST[row-1] = m.eastward[row][0]
	}

	m.SOUTH = m.southward[numRows][1:]

	for col := 1; col <= numCols; col++ {
		go m.NORTH(col)
	}

	for row := 1; row <= numRows; row++ {
		go m.EAST(row)
	}

	for row := 1; row <= numRows; row++ {
		for col := 1; col <= numCols; col++ {
			go m.CENTER(row, col)
		}
	}

	return m
}
