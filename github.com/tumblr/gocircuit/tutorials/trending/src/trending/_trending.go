package trending

import (
	"circuit/kit/sched/limiter"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"circuit/use/n"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
	"tumblr/firehose"
)

// Post is a data structure for keeping track of a post's popularity
type Post struct {
	ID    int64
	Name  string
	Score int32
}

// SortablePosts is a type-wrapper around a slice of posts that is sortable for
// the purposes of using sort.Sort
type SortablePosts []*Post

func (s SortablePosts) Len() int {
	return len(s)
}

func (s SortablePosts) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

func (s SortablePosts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Section: Implementation of mapper logic

// StartMapper.Start is a worker function that launches a mapper worker.
// firehoseConfig specifies the credentials for connecting to the Tumblr Firehose.
// reducer is an array of circuit cross-runtime pointers, listing all available reducers.
func (StartMapper) Start(firehoseConfig *firehose.Request, reducer []circuit.X) {
	circuit.Daemonize(func() {
		f := firehose.Redial(firehoseConfig)
		var n int64
		// Repeat forever: Read an event from the Firehose and pass is on to an appropriate reducer
		for {
			event := f.Read()
			p := &Post{}
			if event.Post != nil {
				p.ID = event.Post.ID
				p.Name = event.Post.BlogName
			}
			if event.Like != nil {
				p.ID = event.Like.RootPostID
			}
			p.Score = 1
			// XXX panic-protect
			reducer[int(p.ID%int64(len(reducer)))].Call("Add", p)
			n++
			if n%100 == 0 {
				println("Consumed", n, "events from the firehose")
			}
		}
	})
}

// Requisite boilerplate for registering StartMapper.Start with the circuit type system
type StartMapper struct{}

func init() { circuit.RegisterFunc(StartMapper{}) }

// Section: Implementation of reducer logic

// Reducer consumes new post events (for a certain subspace of post IDs) and
// maintains a top ten ranking of posts
type Reducer struct {
	sync.Mutex
	m   map[int64]*Post // Map from post IDs to posts
	top []*Post         // Cached copy of the top ten posts on thsi reducer
}

func init() {
	// Values, for which cross-runtime pointers will be created and passed
	// outside of the local runtime, must have their types registered with
	// the circuit type system using RegisterValue.
	circuit.RegisterValue(&Reducer{})
} // Boilerplate

// Add consumes a new post event
func (r *Reducer) Add(p *Post) {
	r.Lock()
	defer r.Unlock()
	q, ok := r.m[p.ID]
	if !ok {
		q = p
		q.Score = 0
		r.m[p.ID] = q
	}
	q.Score += p.Score
	if q.Name == "" {
		q.Name = p.Name
	}
}

// maintainTop periodically computes the top ten posts and caches the results
func (r *Reducer) maintainTop() {
	for {
		time.Sleep(10 * time.Second)
		r.Lock()
		var sp SortablePosts
		for _, p := range r.m {
			sp = append(sp, p)
		}
		sort.Sort(sp)
		sp = sp[:min(10, len(sp))]
		for i, p := range sp {
			q := *p
			sp[i] = &q
		}
		r.top = sp
		r.Unlock()
	}
}

// Top returns a cached copy of the top ten posts
func (r *Reducer) Top() []*Post {
	r.Lock()
	defer r.Unlock()
	return r.top
}

// StartReducer.Start is a worker function that initializes a new Reducer and
// returns a cross-runtime pointer to it
func (StartReducer) Start() circuit.X {
	r := &Reducer{}                      // Create a new reducer object
	r.m = make(map[int64]*Post)          // Create the map holding the posts, indexed by ID
	circuit.Listen("reducer-service", r) // Register Reducer as public service that can be accessed on this worker
	circuit.Daemonize(func() {
		r.maintainTop() // Start a background goroutine that maintains the top ten posts on this reducer
	})
	return circuit.Ref(r) // Make the pointer to the Reducer object exportable and return it
}

type StartReducer struct{} // cruft
func init()                { circuit.RegisterFunc(StartReducer{}) } // cruft

// StartAggregator.Start is a worker function that starts an infinite loop,
// which polls all reducers for their local top ten posts, computes the global
// top ten posts, and prints them out.
func (StartAggregator) Start(reducerAnchor string) {
	circuit.Daemonize(func() {
		for {
			time.Sleep(2 * time.Second)

			// Read anchor directory containing all live reducers
			d, err := anchorfs.OpenDir(reducerAnchor)
			if err != nil {
				println("opendir:", err.Error())
				continue
			}
			// List all anchor files; they correspond to circuit workers hosting Reducer objects
			_, files, err := d.Files()
			if err != nil {
				println("files:", err.Error())
				continue
			}
			// Fetch top ten posts from each reducer, in parallel
			var (
				l   limiter.Limiter
				lk  sync.Mutex
				top SortablePosts
			)
			println("Starting parallel aggregation")
			l.Init(10) // At most 10 concurrent reducer requests at a time
			for _, f_ := range files {
				println("f=", f_.Owner().String())
				f := f_ // Explain...
				l.Go(func() { getReducerTop(f.Owner(), &lk, &top) })
			}
			l.Wait()
			top = top[:min(10, len(top))]
			println("Completed aggregation of", len(top), "best posts")
			// Print the global top ten
			fmt.Printf("Top ten, %s:\n", time.Now().Format(time.UnixDate))
			for i, p := range top {
				fmt.Printf("#% 2d: % 30s id=%d\n", i, p.Name, p.ID)
			}
		}
	})
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func getReducerTop(addr circuit.Addr, lk *sync.Mutex, top *SortablePosts) {
	// Obtain cross-runtime pointer to reducer service on the worker owning file f
	x, err := circuit.TryDial(addr, "reducer-service")
	if err != nil {
		println("dial", addr.String(), "error", err.Error())
		// Skip reducers that seem to be dead
		return
	}
	// Catch panics due to dead worker and return empty list of top ten posts in this case
	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintf(os.Stderr, "%s.Top panic: %#v\n", x.String(), p)
		}
	}()
	// Fetch top ten posts
	rtop := x.Call("Top")[0].([]*Post)
	lk.Lock()
	defer lk.Unlock()
	println("Reducer", addr.String(), "contributed", len(rtop), "posts")
	(*top) = append(*top, rtop...)
}

type StartAggregator struct{} // cruft
func init()                   { circuit.RegisterFunc(StartAggregator{}) } // cruft

// Main is the body of the Go program that starts the trending circuit
func Main() {
	// Start the aggregator worker.
	// Once started, it looks for reducers in the given anchor directory.
	// If no such reducers exits, it waits and retries. It is thus not a
	// problem that the reducers have not been started already.
	println("Kicking aggregator")
	_, addr, err := circuit.Spawn(aggregatorHost, []string{"/tutorial/aggregator"}, StartAggregator{}, "/tutorial/reducer")
	if err != nil {
		panic(err)
	}
	println(addr.String())

	// Start the reducers
	println("Kicking reducers")
	reducer := make([]circuit.X, len(reducerHost))
	for i, h := range reducerHost {
		retrn, addr, err := circuit.Spawn(h, []string{"/tutorial/reducer"}, StartReducer{})
		if err != nil {
			panic(err)
		}
		reducer[i] = retrn[0].(circuit.X)
		println(addr.String())
	}

	// Start the mappers and give them cross-runtime pointers to the already-started reducers.
	println("Kicking mappers")
	for _, h := range mapperHost {
		_, addr, err := circuit.Spawn(h, []string{"/tutorial/mapper"}, StartMapper{}, testFirehose, reducer)
		if err != nil {
			panic(err)
		}
		println(addr.String())
	}
}

var mapperHost = []circuit.Host{
	n.ParseHost("localhost"),
	n.ParseHost("localhost"),
	n.ParseHost("localhost"),
	n.ParseHost("localhost"),
	n.ParseHost("localhost"),
}

var reducerHost = []circuit.Host{
	n.ParseHost("localhost"),
	n.ParseHost("localhost"),
	n.ParseHost("localhost"),
	n.ParseHost("localhost"),
	n.ParseHost("localhost"),
}

var aggregatorHost = n.ParseHost("localhost")
