package main

import (
    "debug/elf"
    "log"
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func printHeader(fh *elf.FileHeader) {
    log.Printf("fh.Class: %s", fh.Class)
    log.Printf("fh.Data: %s", fh.Data)
    log.Printf("fh.Version: %s", fh.Version)
    log.Printf("fh.OSABI: %s", fh.OSABI)
    log.Printf("fh.ABIVersion: %#x", fh.ABIVersion)
    log.Printf("fh.ByteOrder: %s", fh.ByteOrder)
    log.Printf("fh.Type: %s", fh.Type)
    log.Printf("fh.Machine: %s", fh.Machine)
}

func printSection(s *elf.Section) {
    log.Printf("section [Type: %s, Flags, %s, Addr: %#x, Offset: %#x, Size: %#x, Link: %#x, Info: %#x, Addralign: %#x, Entsize: %#x]", s.Type, s.Flags, s.Addr, s.Offset, s.Size, s.Link, s.Info, s.Addralign, s.Entsize)
}

func printProgramHeader(p *elf.Prog) {
    log.Printf("program header [Type: %s, Flags: %s, Off: %#x, Vaddr: %#x, Filesz: %#x, Memsz: %#x, Align: %#x]", p.Type, p.Flags, p.Off, p.Vaddr, p.Filesz, p.Memsz, p.Align)
}

func printSections(s []*elf.Section) {
    log.Printf("file has %d sections", len(s))
    for _, section := range s {
        printSection(section)
    }
}

func printProgs(p []*elf.Prog) {
    log.Printf("file has %d program headers", len(p))
    for _, prog := range p {
        printProgramHeader(prog)
    }
}

func printImportedLibraries(libs []string, err error) {
    if err != nil {
        log.Printf("failed getting imported libraries: %s", err)
    } else {
        log.Printf("file imports %d libraries: %s", len(libs), libs)
    }
}

func printSymbols(symbols []elf.Symbol, err error) {
    if err != nil {
        log.Printf("no symbols: %s", err)
    } else {
        // Grab about 1% of the symbols
        symbolSelection := make([]string, 0, 20)
        for _, symbol := range symbols {
            if rand.Float32() <= 0.01 {
                symbolSelection = append(symbolSelection, symbol.Name)
            }
        }
        log.Printf("there are %d symbols, printing %d of them", len(symbols), len(symbolSelection))
        log.Printf("a selection of symbols: %v", symbols)
    }
}

func printImportedSymbols(importedSymbols []elf.ImportedSymbol, err error) {
    if err != nil {
        log.Printf("no imported symbols: %s", err)
    } else {
        importedSymbolSelection := make([]string, 0, 20)
        for _, symbol := range importedSymbols {
            if rand.Float32() <= 0.1 {
                importedSymbolSelection = append(importedSymbolSelection, symbol.Name+" from "+symbol.Library+",")
            }
        }
        log.Printf("there are %d imported symbols, printing %d of them", len(importedSymbols), len(importedSymbolSelection))
        log.Printf("a selection of imported symbols: %v", importedSymbolSelection)
    }
}

func printFileInformation(f *elf.File) {
    printHeader(&f.FileHeader)
    printSections(f.Sections)
    printProgs(f.Progs)
    printImportedLibraries(f.ImportedLibraries())
    printSymbols(f.Symbols())
    printImportedSymbols(f.ImportedSymbols())
}

func main() {
    file, err := elf.Open("bash.elf")
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer file.Close()
    printFileInformation(file)
}
