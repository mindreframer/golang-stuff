package main

import (
    "debug/macho"
    "log"
    "math/rand"
)

const (
    MH_NOUNDEFS uint32 = 1 << iota /* the object file has no undefined
       references */
    MH_INCRLINK /* the object file is the output of an
       incremental link against a base file
       and can't be link edited again */
    MH_DYLDLINK /* the object file is input for the
       dynamic linker and can't be staticly
       link edited again */
    MH_BINDATLOAD /* the object file's undefined
       references are bound by the dynamic
       linker when loaded. */
    MH_PREBOUND /* the file has its dynamic undefined
       references prebound. */
    MH_SPLIT_SEGS /* the file has its read-only and
       read-write segments split */
    MH_LAZY_INIT /* the shared library init routine is
       to be run lazily via catching memory
       faults to its writeable segments
       (obsolete) */
    MH_TWOLEVEL /* the image is using two-level name
       space bindings */
    MH_FORCE_FLAT /* the executable is forcing all images
       to use flat name space bindings */
    MH_NOMULTIDEFS /* this umbrella guarantees no multiple
       defintions of symbols in its
       sub-images so the two-level namespace
       hints can always be used. */
    MH_NOFIXPREBINDING /* do not have dyld notify the
       prebinding agent about this
       executable */
    MH_PREBINDABLE /* the binary is not prebound but can
       have its prebinding redone.
       only used when MH_PREBOUND is not set. */
    MH_ALLMODSBOUND /* indicates that this binary binds to
       all two-level namespace modules of
       its dependent libraries. only used
       when MH_PREBINDABLE and MH_TWOLEVEL
       are both set. */
    MH_SUBSECTIONS_VIA_SYMBOLS /* safe to divide up the sections into
       sub-sections via symbols for dead
       code stripping */
    MH_CANONICAL /* the binary has been canonicalized
       via the unprebind operation */
    MH_WEAK_DEFINES /* the final linked image contains
       external weak symbols */
    MH_BINDS_TO_WEAK /* the final linked image uses
       weak symbols */
    MH_ALLOW_STACK_EXECUTION /* When this bit is set, all stacks
       in the task will be given stack
       execution privilege.  Only used in
       MH_EXECUTE filetypes. */
    MH_ROOT_SAFE /* When this bit is set, the binary
       declares it is safe for use in
       processes with uid zero */
    MH_SETUID_SAFE /* When this bit is set, the binary
       declares it is safe for use in
       processes when issetugid() is true */
    MH_NO_REEXPORTED_DYLIBS /* When this bit is set on a dylib,
       the static linker does not need to
       examine dependent dylibs to see
       if any are re-exported */
    MH_PIE /* When this bit is set, the OS will
       load the main executable at a
       random address.  Only used in
       MH_EXECUTE filetypes. */
    MH_DEAD_STRIPPABLE_DYLIB /* Only for use on dylibs.  When
       linking against a dylib that
       has this bit set, the static linker
       will automatically not create a
       LC_LOAD_DYLIB load command to the
       dylib if no symbols are being
       referenced from the dylib. */
    MH_HAS_TLV_DESCRIPTORS /* Contains a section of type
       S_THREAD_LOCAL_VARIABLES */
    MH_NO_HEAP_EXECUTION /* When this bit is set, the OS will
       run the main executable with
       a non-executable heap even on
       platforms (e.g. i386) that don't
       require it. Only used in MH_EXECUTE
       filetypes. */
)

func printHeader(fh *macho.FileHeader) {
    log.Printf("fh.Magic: %#x", fh.Magic)
    log.Printf("fh.CPU: %s", fh.Cpu)
    log.Printf("fh.SubCPU: %#x", fh.SubCpu)

    log.Printf("fh.Type: %#x", fh.Type)
    switch fh.Type {
    case macho.TypeExec:
        log.Println("file is an executable")
    case macho.TypeObj:
        log.Println("file is an object")
    default:
        panic("not reachable")
    }

    log.Printf("fh.Ncmd: %d", fh.Ncmd)
    log.Printf("fh.Cmdsz: %d", fh.Cmdsz)
    log.Printf("fh.Flags: %#b", fh.Flags)

    switch fh.Flags & MH_NOUNDEFS {
    case 0:
        log.Println("MH_NOUNDEFS flag is not set")
    default:
        log.Println("object has no undefined references")
    }

    switch fh.Flags & MH_INCRLINK {
    case 0:
        log.Println("MH_INCRLINK flag is not set")
    default:
        log.Println("the object file is the output of an incremental link against a base file and can't be link edited again")
    }

    switch fh.Flags & MH_DYLDLINK {
    case 0:
        log.Println("MH_DYLDLINK flag is not set")
    default:
        log.Println("the object file is input for the dynamic linker and can't be staticly link edited again")
    }

    switch fh.Flags & MH_SETUID_SAFE {
    case 0:
        log.Println("MH_SETUID_SAFE flag is not set")
    default:
        log.Println("executable is setuid safe")
    }
}

func printSection(s *macho.Section) {
    log.Printf("section %s", s.Name)
    log.Printf("\tSeg %s", s.Seg)
    log.Printf("\tAddr %#x", s.Addr)
    log.Printf("\tSize %d", s.Size)
    log.Printf("\tOffset %d", s.Offset)
    log.Printf("\tAlign %d", s.Align)
    log.Printf("\tReloff %s", s.Seg)
    log.Printf("\tNreloc %d", s.Nreloc)
    log.Printf("\tFlags %b", s.Flags)
}

func printSections(sections []*macho.Section) {
    for _, section := range sections {
        printSection(section)
    }
}

func printSymtab(symtab *macho.Symtab) {
    if symtab == nil {
        log.Println("no symbol table")
    }

    log.Printf("symtab.Cmd: %s", symtab.Cmd)
    log.Printf("symtab.Len: %d", symtab.Len)
    log.Printf("symtab.Symoff: %d", symtab.Symoff)
    log.Printf("symtab.Nsyms: %d", symtab.Nsyms)
    log.Printf("symtab.Stroff: %d", symtab.Stroff)
    log.Printf("symtab.Strsize: %d", symtab.Strsize)
    log.Printf("symtab has %d symbols", len(symtab.Syms))

    // Grab about 2.5% of the symbols
    symbols := make([]string, 0, len(symtab.Syms)/40)
    for _, symbol := range symtab.Syms {
        if rand.Float32() <= 0.025 {
            symbols = append(symbols, symbol.Name)
        }
    }
    log.Printf("a selection of the symbols: %v", symbols)
}

func printDysymtab(dysymtab *macho.Dysymtab) {
    log.Printf("dysymtab.Cmd: %s", dysymtab.Cmd)
    log.Printf("dysymtab.Len: %d", dysymtab.Len)
    log.Printf("len(dysymtab.IndirectSyms): %d", len(dysymtab.IndirectSyms))
}

func printImportedLibraries(importedLibraries []string, err error) {
    if err != nil {
        log.Printf("failed getting imported libraries: %s", err)
        return
    }
    log.Printf("file imports %d libraries: %s", len(importedLibraries), importedLibraries)
}

func printFileInformation(f *macho.File) {
    log.Printf("ByteOrder: %s", f.ByteOrder)
    printHeader(&f.FileHeader)

    // Also f.FileHeader.Ncmd
    log.Printf("file has %d load commands", len(f.Loads))
    log.Printf("file has %d sections", len(f.Sections))

    printSections(f.Sections)
    printSymtab(f.Symtab)
    printDysymtab(f.Dysymtab)
    printImportedLibraries(f.ImportedLibraries())
}

func main() {
    file, err := macho.Open("bash.macho")
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer file.Close()
    printFileInformation(file)
}
