package main

import (
    "flag"
    "log"
    "strings"
)

var (
    cmdFlags = map[string]*flag.FlagSet{
        "hello":   flag.NewFlagSet("hello", flag.ExitOnError),
        "goodbye": flag.NewFlagSet("goodbye", flag.ExitOnError),
    }
    subject = cmdFlags["hello"].String("subject", "World", "the subject to say hello to")
    dots    = cmdFlags["goodbye"].Int("dots", 3, "How many dots to print")
)

func hello(subject string) {
    log.Printf("Hello, %s!", subject)
}

func goodbye(dots int) {
    space := ", "
    if dots > 0 {
        space = strings.Repeat(".", dots)
    }
    log.Printf("Goodbye%scruel world!", space)
}

func main() {
    flag.Parse()
    for _, cmd := range flag.Args() {
        flags, ok := cmdFlags[cmd]
        if !ok {
            log.Fatalf("no command %q found", cmd)
        }
        flags.Parse(flag.Args()[1:])
        switch cmd {
        case "hello":
            hello(*subject)
        case "goodbye":
            goodbye(*dots)
        }
        break
    }
}
