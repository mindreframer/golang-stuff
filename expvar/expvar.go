package main

import (
    "expvar"
    "flag"
    "log"
    "time"
)

var (
    times      = flag.Int("times", 1, "times to say hello")
    name       = flag.String("name", "World", "thing to say hello to")
    helloTimes = expvar.NewInt("hello")
)

func init() {
    expvar.Publish("time", expvar.Func(now))
}

func now() interface{} {
    return time.Now().Format(time.RFC3339Nano)
}

func hello(times int, name string) {
    helloTimes.Add(int64(times))
    for i := 0; i < times; i++ {
        log.Printf("Hello, %s!", name)
    }
}

func printVars() {
    log.Println("expvars:")
    expvar.Do(func(kv expvar.KeyValue) {
        switch kv.Key {
        case "memstats":
            // Do nothing, this is a big output.
        default:
            log.Printf("\t%s -> %s", kv.Key, kv.Value)
        }
    })
}

func main() {
    flag.Parse()
    printVars()
    hello(*times, *name)
    printVars()
    hello(*times, *name)
    printVars()
}
