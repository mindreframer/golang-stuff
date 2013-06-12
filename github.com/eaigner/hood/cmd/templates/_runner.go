package main

import (
	"flag"
	"github.com/eaigner/hood"
	"io/ioutil"
	"log"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var (
	driver     string
	source     string
	schemaPath string
	steps      int
)

func init() {
	flag.StringVar(&driver, "driver", "", "Sets the driver")
	flag.StringVar(&source, "source", "", "Sets the source")
	flag.StringVar(&schemaPath, "schema", "db/schema.go", "Sets the schema path")
	flag.IntVar(&steps, "steps", 0, "Sets the steps")
	flag.Parse()

	if driver == "" {
		panic("driver not set")
	}
	if source == "" {
		panic("source not set")
	}

	// Unquote
	d, err := strconv.Unquote(driver)
	if err != nil {
		panic(err)
	}
	driver = d

	s, err := strconv.Unquote(source)
	if err != nil {
		panic(err)
	}
	source = s

	sp, err := strconv.Unquote(schemaPath)
	if err != nil {
		panic(err)
	}
	schemaPath = sp
}

type M struct{}

type Migrations struct {
	Id      hood.Id
	Current int
}

func main() {
	// Print action
	if steps > 0 {
		log.Printf("applying migrations...")
	} else if steps == -1 {
		log.Printf("rolling back by 1...")
	} else if steps < 0 {
		log.Printf("reset. rolling back all migrations...")
	}

	// Parse migrations
	stamps := []int{}
	ups := map[int]reflect.Method{}
	downs := map[int]reflect.Method{}

	structVal := reflect.ValueOf(&M{})
	for i := 0; i < structVal.NumMethod(); i++ {
		method := structVal.Type().Method(i)
		if c := strings.Split(method.Name, "_"); len(c) >= 3 {
			stamp, _ := strconv.Atoi(c[len(c)-2])
			if c[len(c)-1] == "Up" {
				ups[stamp] = method
				stamps = append(stamps, stamp)
			} else {
				downs[stamp] = method
			}
		}
	}

	sort.Ints(stamps)

	// Open hood
	hd, err := hood.Open(driver, source)
	if err != nil {
		panic(err)
	}
	hd.Log = true

	// Create migration table if necessary
	tx := hd.Begin()
	tx.CreateTableIfNotExists(&Migrations{})
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	// Check if any previous migrations have been run
	var rows []Migrations
	err = hd.Find(&rows)
	if err != nil {
		panic(err)
	}
	if len(rows) > 1 {
		panic("invalid migrations table")
	}
	info := Migrations{}
	if len(rows) > 0 {
		info = rows[0]
	}

	// Apply
	cur := 0
	count := 0
	if steps > 0 {
		for _, stamp := range stamps {
			if stamp > info.Current {
				if cur++; cur <= steps {
					apply(stamp, stamp, &count, hd, &info, structVal, ups[stamp])
				}
			}
		}
	} else if steps < 0 {
		for i := len(stamps) - 1; i >= 0; i-- {
			stamp := stamps[i]
			next := 0
			if i > 0 {
				next = stamps[i-1]
			}
			if stamp <= info.Current {
				if cur--; cur >= steps {
					apply(stamp, next, &count, hd, &info, structVal, downs[stamp])
				}
			}
		}
	}

	if steps > 0 {
		log.Printf("applied %d migrations", count)
	} else if steps < 0 {
		log.Printf("rolled back %d migrations", count)
	}

	log.Printf("generating new schema... %s", schemaPath)

	dry := hood.Dry()
	for _, ts := range stamps {
		if ts <= info.Current {
			method := ups[ts]
			method.Func.Call([]reflect.Value{structVal, reflect.ValueOf(dry)})
		}
	}
	err = ioutil.WriteFile(schemaPath, []byte(dry.GoSchema()), 0666)
	if err != nil {
		panic(err)
	}
	err = exec.Command("go", "fmt", schemaPath).Run()
	if err != nil {
		panic(err)
	}
	log.Printf("wrote schema %s", schemaPath)
	log.Printf("done.")
}

func apply(stamp, current int, count *int, hd *hood.Hood, info *Migrations, structVal reflect.Value, method reflect.Method) {
	log.Printf("applying %s...", method.Name)
	txn := hd.Begin()
	method.Func.Call([]reflect.Value{structVal, reflect.ValueOf(txn)})
	info.Current = current
	txn.Save(info)
	err := txn.Commit()
	if err != nil {
		panic(err)
	} else {
		*count++
	}

}
