package main

import (
    "bytes"
    "log"
    "unicode"
)

func main() {
    quickBrownFox := []byte("The quick brown fox jumped over the lazy dog")

    title := bytes.Title(quickBrownFox)
    log.Printf("Title turned %q into %q", quickBrownFox, title)

    allTitle := bytes.ToTitle(quickBrownFox)
    log.Printf("ToTitle turned %q to %q", quickBrownFox, allTitle)

    allTitleTurkish := bytes.ToTitleSpecial(unicode.TurkishCase, quickBrownFox)
    log.Printf("ToTitleSpecial turned %q into %q using the Turkish case rules", quickBrownFox, allTitleTurkish)

    lower := bytes.ToLower(title)
    log.Printf("ToLower turned %q into %q", title, lower)

    turkishCapitalI := []byte("İ")
    turkishLowerI := bytes.ToLowerSpecial(unicode.TurkishCase, turkishCapitalI)
    log.Printf("ToLowerSpecial turned %q into %q using the Turkish case rules", turkishCapitalI, turkishLowerI)

    upper := bytes.ToUpper(quickBrownFox)
    log.Printf("ToUpper turned %q to %q", quickBrownFox, upper)

    upperSpecial := bytes.ToUpperSpecial(unicode.TurkishCase, quickBrownFox)
    log.Printf("ToUpperSpecial turned %q into %q using the Turkish case rules", quickBrownFox, upperSpecial)
}
