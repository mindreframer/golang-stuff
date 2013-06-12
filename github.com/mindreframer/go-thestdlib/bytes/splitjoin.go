package main

import (
    "bytes"
    "log"
    "strings"
)

func main() {
    languages := []byte("golang haskell ruby python")

    individualLanguages := bytes.Fields(languages)
    log.Printf("Fields split %q on whitespace into %q", languages, individualLanguages)

    vowelsAndSpace := "aeiouy "
    split := bytes.FieldsFunc(languages, func(r rune) bool {
        return strings.ContainsRune(vowelsAndSpace, r)
    })
    log.Printf("FieldsFunc split %q on vowels and space into %q", languages, split)

    space := []byte{' '}
    splitLanguages := bytes.Split(languages, space)
    log.Printf("Split split %q on a single space into %q", languages, splitLanguages)

    numberOfSubslices := 2 // Not number of splits
    singleSplit := bytes.SplitN(languages, space, numberOfSubslices)
    log.Printf("SplitN split %q on a single space into %d subslices: %q", languages, numberOfSubslices, singleSplit)

    splitAfterLanguages := bytes.SplitAfter(languages, space)
    log.Printf("SplitAfter split %q AFTER a single space (keeping the space) into %q", languages, splitAfterLanguages)

    splitAfterNLanguages := bytes.SplitAfterN(languages, space, numberOfSubslices)
    log.Printf("SplitAfterN split %q AFTER a single space (keeping the space) into %d subslices: %q", languages, numberOfSubslices, splitAfterNLanguages)

    languagesBackTogether := bytes.Join(individualLanguages, space)
    log.Printf("Languages are back togeher again! %q == %q? %v", languagesBackTogether, languages, bytes.Equal(languagesBackTogether, languages))
}
