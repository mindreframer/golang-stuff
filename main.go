// cmail is a command that runs another command and sends stdout and stderr
// to a specified email address at certain intervals.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/cmd"
)

var (
	flagSendMail      = "sendmail"
	flagTo            = os.Getenv("EMAIL")
	flagPeriod        = "1h"
	flagNoPeriod      = false
	flagNoPass        = false
	flagInc           = false
	flagSubjectPrefix = "[cmail] "

	fullProgram string
)

func init() {
	flag.StringVar(&flagSendMail, "sendmail", flagSendMail,
		"The command to use to send mail. The email content and headers\n"+
			"will be sent to stdin.")
	flag.StringVar(&flagTo, "to", flagTo,
		"The email address to send mail to. By default, this is set to the\n"+
			"value of the $EMAIL environment variable.")
	flag.StringVar(&flagSubjectPrefix, "subj", flagSubjectPrefix,
		"A subject prefix to use for all emails.")
	flag.StringVar(&flagPeriod, "period", flagPeriod,
		"The amount of time to wait between sending data gathered from\n"+
			"stdin. Value should be a duration defined by Go's\n"+
			"time.ParseDuration. e.g., '300ms', '1.5h', '1m'.")
	flag.BoolVar(&flagNoPass, "no-pass", flagNoPass,
		"If set, stdout/stderr will not be passed thru.")
	flag.BoolVar(&flagInc, "inc", flagInc,
		"If set, emails will contain incremental changes as opposed to\n"+
			"each email containing all data. Will also use less memory.")
	flag.BoolVar(&flagNoPeriod, "no-period", flagNoPeriod,
		"If set, only one email will be sent when the command finishes.")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	flagTo = strings.TrimSpace(flagTo)

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	if len(flagTo) == 0 {
		log.Println("I don't know who to send email to. Please use the\n" +
			"'-to' flag or set the EMAIL environment variable.\n")
		flag.Usage()
		os.Exit(1)
	}

	period, err := time.ParseDuration(flagPeriod)
	assert(err, "Invalid period '%s': %s.", flagPeriod, err)

	sigged := make(chan os.Signal)
	signal.Notify(sigged, os.Interrupt, os.Kill)

	var program *exec.Cmd
	var inlines <-chan string
	if flag.NArg() == 0 {
		fullProgram = "stdin"

		inlines = gobble(bufio.NewReader(os.Stdin))
	} else {
		program = exec.Command(flag.Arg(0), flag.Args()[1:]...)
		fullProgram = strings.Join(flag.Args(), " ")
		if len(fullProgram) > 200 {
			fullProgram = fullProgram[0:201]
		}

		stdout, err := program.StdoutPipe()
		assert(err, "Could not get stdout pipe: %s.", err)

		stderr, err := program.StderrPipe()
		assert(err, "Could not get stderr pipe: %s.", err)

		err = program.Start()
		assert(err, "Could not start program '%s': %s.", fullProgram, err)

		// Start goroutines for reading stdout and stderr, then mux them.
		stdoutLines := gobble(bufio.NewReader(stdout))
		stderrLines := gobble(bufio.NewReader(stderr))
		inlines = muxer(stdoutLines, stderrLines)
	}

	// Start the goroutine responsible for sending emails.
	// The send is also responsible for quitting the program.
	// (When all emails remaining have been sent.)
	send := sender()

	// Keep track of all lines emitted to stdout/stderr.
	// If the duration passes, stop and send whatever we have.
	// If the user interrupts the program, stop and send whatever we have.
	// If EOF is read on both stdout and stderr, send what we have.
	// We exit the program by closing the `send` channel, which will force
	// any remaining emails left to be sent.
	outlines := make([]string, 0)
	addMsg := func(msg string) {
		outlines = append(outlines, []string{"\n", "\n", msg + "\n"}...)
	}
	killed := false // set if user interrupted
	for {
		select {
		case <-time.After(period):
			if !flagNoPeriod {
				send <- outlines
				outlines = outlines[:0]
			}
		case <-sigged:
			if program != nil {
				program.Process.Kill()
				killed = true
				// continue reading stdout/stderr until program really quits.
			} else { // reading stdin, so send what we've got now.
				send <- outlines
				close(send)
				select {}
			}
		case line, ok := <-inlines:
			if !ok {
				// Program completed successfully!
				if killed {
					addMsg("Program interrupted.")
				} else {
					addMsg("Program completed successfully.")
				}
				send <- outlines
				close(send)
				select {}
			}
			outlines = append(outlines, line)
		}
	}
}

// muxer takes a list of incoming string channels, and muxes them all into
// the result channel. The channel returned is closed if and only if all
// channels in 'ins' have been closed.
func muxer(ins ...<-chan string) <-chan string {
	wg := new(sync.WaitGroup)
	combined := make(chan string, 500)
	for _, in := range ins {
		wg.Add(1)
		in := in
		go func() {
			for line := range in {
				combined <- line
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(combined)
	}()
	return combined
}

// sender receives chunks of lines on the result channel, and sends those
// chunks of lines via email.
//
// sender is the only goroutine that should exit the program in normal
// operation, which happens when there are no more chunks of lines to read.
func sender() chan<- []string {
	toSend := make([]string, 500)
	send := make(chan []string)
	go func() {
		for newLines := range send {
			switch {
			case flagInc:
				if len(newLines) == 0 {
					emailLines([]string{"Nothing to report."})
				} else {
					emailLines(newLines)
				}
			default:
				toSend = append(toSend, newLines...)
				emailLines(toSend)
			}
		}
		os.Exit(0)
	}()
	return send
}

// gobble reads lines from any buffered reader, sends the lines on the result
// channel, and closes the channel when EOF is read.
//
// gobble will quit the program with an error message if the input source
// cannot be read.
func gobble(buf *bufio.Reader) <-chan string {
	lines := make(chan string)
	go func() {
		for {
			line, err := buf.ReadString('\n')
			if err != nil && err != io.EOF {
				fatal("Could not read line: %s", err)
			}
			if !flagNoPass {
				fmt.Print(line)
			}
			lines <- line
			if err == io.EOF {
				close(lines)
				break
			}
		}
	}()
	return lines
}

// emailLines sends a chunk of lines via email.
func emailLines(lines []string) {
	var c *cmd.Command
	subj := fmt.Sprintf("%s%s", flagSubjectPrefix, fullProgram)

	if flagSendMail == "mailx" {
		c = cmd.New(flagSendMail, "-s", subj, flagTo)
		fmt.Fprint(c.BufStdin, strings.Join(lines, ""))
	} else {
		c = cmd.New(flagSendMail, "-t")
		date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700")
		fmt.Fprintf(c.BufStdin,
			`Subject: %s
From: %s
To: %s
Date: %s

%s`, subj, flagTo, flagTo, date, strings.Join(lines, ""))
	}

	if err := c.Run(); err != nil {
		log.Printf("Error sending mail '%s -t': %s.", flagSendMail, err)
	}
}

func usage() {
	log.Printf("Usage: %s [flags] command [args]\n\n", path.Base(os.Args[0]))
	log.Printf("cmail sends data read from `command` periodically, and/or\n" +
		"when EOF is reached.\n\n")

	flag.VisitAll(func(fl *flag.Flag) {
		var def string
		if len(fl.DefValue) > 0 {
			def = fmt.Sprintf(" (default: %s)", fl.DefValue)
		}
		log.Printf("-%s%s\n", fl.Name, def)
		log.Printf("    %s\n", strings.Replace(fl.Usage, "\n", "\n    ", -1))
	})
}

func assert(err error, format string, v ...interface{}) {
	if err != nil {
		fatal(format, v...)
	}
}

func fatal(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
