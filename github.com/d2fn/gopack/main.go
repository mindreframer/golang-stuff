package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	GopackDir          = ".gopack"
	GopackTestProjects = ".gopack/test-projects"
	VendorDir          = ".gopack/vendor"
)

const (
	Blue     = uint8(94)
	Green    = uint8(92)
	Red      = uint8(31)
	Gray     = uint8(90)
	EndColor = "\033[0m"
)

var (
	pwd        string
	showColors = true
)

func main() {
	if os.Getenv("GOPACK_SKIP_COLORS") == "1" {
		showColors = false
	}

	fmtcolor(104, "/// g o p a c k ///")
	fmt.Println()
	// localize GOPATH
	setupEnv()
	p, err := AnalyzeSourceTree(".")
	if err != nil {
		fail(err)
	}
	d := LoadDependencyModel(".")
	failWith(d.Validate(p))
	// prepare dependencies
	d.VisitDeps(
		func(d *Dep) {
			fmtcolor(Gray, "updating %s\n", d.Import)
			d.goGetUpdate()
			fmtcolor(Gray, "pointing %s at %s %s\n", d.Import, d.CheckoutType(), d.CheckoutSpec)
			d.switchToBranchOrTag()
		})
	// run the specified command
	cmd := exec.Command("go", os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fail(err)
	}
}

// set GOPATH to the local vendor dir
func setupEnv() {
	dir, err := os.Getwd()
	pwd = dir
	if err != nil {
		fail(err)
	}
	vendor := fmt.Sprintf("%s/%s", pwd, VendorDir)
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		fail(err)
	}
}

func fmtcolor(c uint8, s string, args ...interface{}) {
	if showColors {
		fmt.Printf("\033[%dm", c)
	}

	if len(args) > 0 {
		fmt.Printf(s, args...)
	} else {
		fmt.Printf(s)
	}

	if showColors {
		fmt.Printf(EndColor)
	}
}

func logcolor(c uint8, s string, args ...interface{}) {
	log.Printf("\033[%dm", c)
	if len(args) > 0 {
		log.Printf(s, args...)
	} else {
		log.Printf(s)
	}
	log.Printf(EndColor)
}

func failf(s string, args ...interface{}) {
	fmtcolor(Red, s, args...)
	os.Exit(1)
}

func fail(a ...interface{}) {
	fmt.Printf("\033[%dm", Red)
	fmt.Print(a)
	fmt.Printf(EndColor)
	os.Exit(1)
}

func failWith(errors []*ProjectError) {
	if len(errors) > 0 {
		fmt.Printf("\033[%dm", Red)
		for _, e := range errors {
			fmt.Printf(e.String())
		}
		fmt.Printf(EndColor)
		fmt.Println()
		os.Exit(len(errors))
	}
}
