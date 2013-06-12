package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	dbDir             string
	migrationsDir     string
	configFile        string
	schemaFile        string
	environment       string
	keyPath           string
	driver            string
	source            string
	command           string
	workingDir        string
	testMode          bool
	migrationTemplate = template.Must(template.New("migration").Parse(_migrationTemplate))
)

func init() {
	flag.StringVar(&dbDir, "db-dir", "db", "Sets the db directory.")
	flag.StringVar(&migrationsDir, "migrations-dir", "db/migrations", "Sets the migrations directory")
	flag.StringVar(&configFile, "conf", "db/config.json", "Sets the path to the config.json file")
	flag.StringVar(&schemaFile, "schema", "db/schema.go", "Sets the generated schema file path")
	flag.StringVar(&environment, "env", "development", "Sets the environment")
	flag.StringVar(&keyPath, "key-path", "", `Sets the key path of the driver and source fields relative to the environment field in the config.json. For example, if the driver and source fields were located at <environment>:{"postgres":{"hood":{"driver": ... }}} the key path would be "postgres.hood"`)
	flag.StringVar(&driver, "driver", "", "Sets the driver. If the driver and source fields are set, the config.json will be ignored.")
	flag.StringVar(&source, "source", "", "Sets the source. If the driver and source fields are set, the config.json will be ignored.")
	flag.BoolVar(&testMode, "test", false, "If set to true, just the configuration is printed without running a command.")
	flag.Parse()

	// Get driver and source from config, if flags were not set
	if driver == "" || source == "" {
		readConf()
	}

	// Get working dir
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	workingDir = wd
	command = flag.Arg(0)
}

func readConf() {
	// Read file
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.Open(path.Join(wd, configFile))
	if err != nil {
		return
		// panic(err)
	}
	defer f.Close()

	// Parse JSON
	var v map[string]interface{}
	err = json.NewDecoder(f).Decode(&v)
	if err != nil {
		panic(err)
	}

	// Append environment to key path
	fqKeyPath := environment
	if len(keyPath) > 0 {
		fqKeyPath += "." + keyPath
	}

	// Get driver and source fields
	root := v
	for _, key := range strings.Split(fqKeyPath, ".") {
		if node, ok := root[key].(map[string]interface{}); ok {
			root = node
		}
	}

	if node, ok := root["driver"].(string); ok {
		driver = node
	}
	if node, ok := root["source"].(string); ok {
		source = node
	}
}

func main() {
	// Only print config in test mode
	if testMode {
		fmt.Printf("db-dir:\t\t'%s'\n", dbDir)
		fmt.Printf("migrations-dir:\t'%s'\n", migrationsDir)
		fmt.Printf("conf:\t\t'%s'\n", configFile)
		fmt.Printf("schema:\t\t'%s'\n", schemaFile)
		fmt.Printf("env:\t\t'%s'\n", environment)
		fmt.Printf("key-path:\t'%s'\n", keyPath)
		fmt.Printf("driver:\t\t'%s'\n", driver)
		fmt.Printf("source:\t\t'%s'\n", source)
		return
	}

	// Run command
	switch command {
	case "create:config":
		cmdCreateConfig()
	case "create:migration":
		cmdCreateMigration(flag.Arg(1))
	case "db:migrate":
		cmdMigrate()
	case "db:rollback":
		cmdRollback()
	case "db:reset":
		cmdReset()
	default:
		log.Println("invalid command")
	}
}

func cmdCreateConfig() {
	p := path.Join(workingDir, configFile)
	err := os.MkdirAll(path.Dir(p), 0777)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(p, []byte(confTmpl), 0666)
	if err != nil {
		panic(err)
	}
	log.Printf("created db configuration '%s'", configFile)
}

func cmdCreateMigration(name string) {
	if name == "" {
		panic("invalid migration name")
	}
	// Write template
	ts := time.Now().Unix()

	p := path.Join(workingDir, migrationsDir, fmt.Sprintf("%d_%s.go", ts, name))
	err := os.MkdirAll(path.Dir(p), 0777)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = migrationTemplate.Execute(f, &struct {
		Timestamp int64
		Name      string
	}{
		Timestamp: ts,
		Name:      name,
	})

	if err != nil {
		panic(err)
	}

	// Write init.go file
	p = path.Join(workingDir, migrationsDir, "init.go")
	err = ioutil.WriteFile(p, []byte(initTmpl), 0666)
	if err != nil {
		panic(err)
	}
	log.Printf("created migration '%s'", p)
}

var initTmpl = `package main

type M struct{}

func main() {}
`

var confTmpl = `{
  "development": {
    "driver": "",
    "source": ""
  },
  "production": {
    "driver": "",
    "source": ""
  },
  "test": {
    "driver": "",
    "source": ""
  }
}`

func cmdMigrate() {
	runMigrations(math.MaxInt32)
}

func cmdRollback() {
	runMigrations(-1)
}

func cmdReset() {
	runMigrations(math.MinInt32)
}

func runMigrations(steps int32) {
	// Read the migrations dir
	info, err := ioutil.ReadDir(path.Join(workingDir, migrationsDir))
	if err != nil {
		panic(nil)
	}

	// Create a temp dir to copy the migrations to
	tmpDir, err := ioutil.TempDir("", "hood-migration-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	// Copy all migrations
	files := []string{}
	for _, file := range info {
		// Skip the init file
		if file.Name() == "init.go" {
			continue
		}

		// Skip non-go files
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		// Copy
		dstFile := path.Join(tmpDir, file.Name())
		_, err = copyFile(
			dstFile,
			path.Join(workingDir, migrationsDir, file.Name()),
		)
		if err != nil {
			panic(err)
		}
		files = append(files, dstFile)
	}

	// Copy the runner template
	main := path.Join(tmpDir, "runner.go")
	err = ioutil.WriteFile(main, []byte(_runnerTemplate), 0666)
	if err != nil {
		panic(err)
	}
	files = append(files, main)

	// Adjust schema path
	if !path.IsAbs(schemaFile) {
		schemaFile = path.Join(workingDir, schemaFile)
	}

	// Invoke runner
	cmd := exec.Command("go", "run")
	cmd.Args = append(cmd.Args, files...)
	cmd.Args = append(cmd.Args,
		"-driver", strconv.Quote(driver),
		"-source", strconv.Quote(source),
		"-schema", strconv.Quote(schemaFile),
		"-steps", fmt.Sprintf("%d", steps),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		panic(err)
	}
}

func copyFile(dst, src string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()
	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()
	return io.Copy(df, sf)
}
