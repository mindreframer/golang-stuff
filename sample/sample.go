package main

import (
	"fmt"
	"github.com/emilsjolander/goson"
)

type Repo struct {
	Name  string
	URL   string
	Stars int
	Forks int
}

type User struct {
	Name  string
	Repos []Repo
}

func main() {

	user := &User{
		Name: "Emil Sjölander",
		Repos: []Repo{
			Repo{
				Name:  "goson",
				URL:   "https://github.com/emilsjolander/goson",
				Stars: 0,
				Forks: 0,
			},
			Repo{
				Name:  "StickyListHeaders",
				URL:   "https://github.com/emilsjolander/StickyListHeaders",
				Stars: 722,
				Forks: 197,
			},
			Repo{
				Name:  "android-FlipView",
				URL:   "https://github.com/emilsjolander/android-FlipView",
				Stars: 157,
				Forks: 47,
			},
		},
	}

	result, err := goson.Render("user", goson.Args{"User": user})

	if err != nil {
		panic(err)
	}

	fmt.Println(string(result))
}
