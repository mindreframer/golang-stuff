package main

import (
    "bytes"
    "log"
)

func contains(s, sub []byte) {
    if bytes.Contains(s, sub) {
        log.Printf("%s contains %s", s, sub)
    } else {
        log.Printf("%s does NOT contain %s", s, sub)
    }
}

func count(s, sep []byte) {
    log.Printf("%s contains %d instance(s) of %s", s, bytes.Count(s, sep), sep)
}

func hasPrefix(s, prefix []byte) {
    if bytes.HasPrefix(s, prefix) {
        log.Printf("%s has the prefix %s", s, prefix)
    } else {
        log.Printf("%s does NOT have the prefix %s", s, prefix)
    }
}

func hasSuffix(s, suffix []byte) {
    if bytes.HasSuffix(s, suffix) {
        log.Printf("%s has the suffix %s", s, suffix)
    } else {
        log.Printf("%s does NOT have the suffix %s", s, suffix)
    }
}

func index(s, sep []byte) {
    if i := bytes.Index(s, sep); i == -1 {
        log.Printf("%s does NOT appear in %s", sep, s)
    } else {
        log.Printf("%s appears at index %d in %s", sep, i, s)
    }
}

func indexAny(s []byte, chars string) {
    if i := bytes.IndexAny(s, chars); i == -1 {
        log.Printf("No unicode characters in %q appear in %s", chars, s)
    } else {
        log.Printf("A unicode character in %q appears at index %d in %s", chars, i, s)
    }
}

func indexByte(s []byte, b byte) {
    if i := bytes.IndexByte(s, b); i == -1 {
        log.Printf("%q does NOT appear in %s", b, s)
    } else {
        log.Printf("%q appears at index %d in %s", b, i, s)
    }
}

func indexFunc(s []byte, f func(rune) bool) {
    if i := bytes.IndexFunc(s, f); i == -1 {
        log.Printf("Something controlled by %#v does NOT appear in %s", f, s)
    } else {
        log.Printf("Something controlled by %#v appears at index %d in %s", f, i, s)
    }
}

func indexRune(s []byte, r rune) {
    if i := bytes.IndexRune(s, r); i == -1 {
        log.Printf("Rune %d does NOT appear in %s", r, s)
    } else {
        log.Printf("Rune %d appears at index %d in %s", r, i, s)
    }
}

func lastIndex(s, sep []byte) {
    if i := bytes.LastIndex(s, sep); i == -1 {
        log.Printf("%s does NOT appear in %s", sep, s)
    } else {
        log.Printf("%s appears last at index %d in %s", sep, i, s)
    }
}

func lastIndexAny(s []byte, chars string) {
    if i := bytes.LastIndexAny(s, chars); i == -1 {
        log.Printf("No unicode characters in %q appear in %s", chars, s)
    } else {
        log.Printf("A unicode character in %q appears last at index %d in %s", chars, i, s)
    }
}

func lastIndexFunc(s []byte, f func(rune) bool) {
    if i := bytes.LastIndexFunc(s, f); i == -1 {
        log.Printf("Something controlled by %#v does NOT appear in %s", f, s)
    } else {
        log.Printf("Something controlled by %#v appears at index %d in %s", f, i, s)
    }
}

func main() {
    golang := []byte("golang")
    haskell := []byte("haskell")
    lang := []byte("lang")
    gos := []byte("go")

    contains(golang, lang)
    contains(golang, haskell)

    count(golang, lang)
    count(haskell, []byte("l"))

    hasPrefix(golang, gos)
    hasPrefix(haskell, gos)

    hasSuffix(golang, lang)
    hasSuffix(haskell, lang)

    index(golang, lang)
    index(golang, gos)
    index(haskell, lang)

    indexAny(golang, "lang")
    indexAny(haskell, "lang")
    indexAny(haskell, "go")

    indexByte(golang, 'h')
    indexByte(golang, 'l')
    indexByte(haskell, 'l')

    g := rune('g')
    indexFunc(golang, func(r rune) bool { return r == g })
    indexFunc(haskell, func(r rune) bool { return r == g })

    indexRune(golang, rune('o'))
    indexRune(haskell, rune('l'))

    lastIndex(golang, []byte("g"))
    lastIndex(haskell, []byte("l"))

    lastIndexAny(golang, "abcdefg")
    lastIndexAny(haskell, "lmnop")

    lastIndexFunc(golang, func(r rune) bool { return r == g })
    lastIndexFunc(haskell, func(r rune) bool { return r == g })
}
