package main

import (
    "debug/elf"
    "log"
)

func printDwarfInformation(f *elf.File) {
    dwarf, err := f.DWARF()
    if err != nil {
        log.Printf("failed getting DWARF info: %s", err)
        return
    }

    rd := dwarf.Reader()
    for {
        entry, err := rd.Next()
        if err != nil {
            log.Printf("failed getting next DWARF entry: %s", err)
            return
        }
        if entry == nil {
            // All done
            return
        }
        log.Printf("got entry with tag: %s, and offset %d", entry.Tag, entry.Offset)
        for _, field := range entry.Field {
            log.Printf("\t%s: %v", field.Attr, field.Val)
        }
    }
}

func main() {
    file, err := elf.Open("hello")
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer file.Close()
    printDwarfInformation(file)
}
