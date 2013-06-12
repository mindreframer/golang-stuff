package main

import (
    "debug/elf"
    "debug/gosym"
    "log"
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func printSyms(syms []gosym.Sym) {
    selection := make([]string, 0, 24)
    for _, sym := range syms {
        if sym.Name != "" {
            if rand.Float32() <= 0.005 {
                selection = append(selection, sym.Name)
            }
        }
    }
    log.Printf("there are %d symbols, printing %d of them", len(syms), len(selection))
    log.Printf("a selection of symbols: %v", selection)
}

func printFuncs(funcs []gosym.Func) {
    selection := make([]string, 0, 24)
    for _, f := range funcs {
        if rand.Float32() <= 0.005 {
            selection = append(selection, f.Name)
        }
    }
    log.Printf("there are %d functions, printing %d of them", len(funcs), len(selection))
    log.Printf("a selection of functions: %v", selection)
}

func printFiles(files map[string]*gosym.Obj) {
    selection := make([]string, 0, 24)
    for name := range files {
        if rand.Float32() <= 0.02 {
            selection = append(selection, name)
        }
    }
    log.Printf("there are %d files, printing %d of them", len(files), len(selection))
    log.Printf("a selection of files: %v", selection)
}

func getSectionData(f *elf.File, name string) []byte {
    section := f.Section(name)
    if section == nil {
        log.Fatalf("failed getting section %s", name)
    }
    data, err := section.Data()
    if err != nil {
        log.Fatalf("failed getting section %s data: %s", name, err)
    }
    return data
}

func processGoInformation(f *elf.File) {
    gosymtab := getSectionData(f, ".gosymtab")
    gopclntab := getSectionData(f, ".gopclntab")

    lineTable := gosym.NewLineTable(gopclntab, f.Section(".text").Addr)
    table, err := gosym.NewTable(gosymtab, lineTable)
    if err != nil {
        log.Fatalf("failed making table: %s", err)
    }

    printSyms(table.Syms)
    printFuncs(table.Funcs)
    printFiles(table.Files)
}

func main() {
    file, err := elf.Open("doozerd")
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer file.Close()
    processGoInformation(file)
}
