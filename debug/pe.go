package main

import (
    "debug/pe"
    "log"
)

func printFileHeader(fh pe.FileHeader) {
    log.Printf("fh.Machine: %d", fh.Machine)
    log.Printf("fh.NumberOfSections: %d", fh.NumberOfSections)
    log.Printf("fh.TimeDateStamp: %d", fh.TimeDateStamp)
    log.Printf("fh.PointerToSymbolTable: %#x", fh.PointerToSymbolTable)
    log.Printf("fh.NumberOfSymbols: %d", fh.NumberOfSymbols)
    log.Printf("fh.SizeOfOptionalHeader: %d", fh.SizeOfOptionalHeader)
    log.Printf("fh.Characteristics: %#x", fh.Characteristics)
}

func printSection(s *pe.Section) {
    log.Printf("section %s", s.Name)
    log.Printf("\tVirtualSize: %d", s.VirtualSize)
    log.Printf("\tVirtualAddress: %d", s.VirtualAddress)
    log.Printf("\tSize: %d", s.Size)
    log.Printf("\tOffset: %d", s.Offset)
    log.Printf("\tPointerToRelocations: %d", s.PointerToRelocations)
    log.Printf("\tPointerToLineNumbers: %d", s.PointerToLineNumbers)
    log.Printf("\tNumberOfRelocations: %d", s.NumberOfRelocations)
    log.Printf("\tNumberOfLineNumbers: %d", s.NumberOfLineNumbers)
    log.Printf("\tCharacteristics: %d", s.Characteristics)
}

func printSections(sections []*pe.Section) {
    for _, section := range sections {
        printSection(section)
    }
}

func printImportedLibraries(importedLibraries []string, err error) {
    if err != nil {
        log.Printf("failed getting imported libraries: %s", err)
        return
    }
    log.Printf("file imports %d libraries: %s", len(importedLibraries), importedLibraries)
}

func printImportedSymbols(importedSymbols []string, err error) {
    if err != nil {
        log.Printf("failed getting imported symbols: %s", err)
        return
    }
    log.Printf("file imports %d symbols: %s", len(importedSymbols), importedSymbols)
}

func printFileInformation(f *pe.File) {
    printFileHeader(f.FileHeader)
    printSections(f.Sections)
    printImportedLibraries(f.ImportedLibraries())
    printImportedSymbols(f.ImportedSymbols())
}

func main() {
    file, err := pe.Open("Hello.exe")
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer file.Close()
    printFileInformation(file)
}
