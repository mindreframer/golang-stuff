package main

import (
    "bytes"
    "encoding/xml"
    "io"
    "log"
)

type Name struct {
    First, Last string `xml:",omitempty"`
}

type Author struct {
    Id   int `xml:"id,attr"`
    Name Name
}

type BlogPost struct {
    XMLName  xml.Name `xml:"Post"`
    Id       int      `xml:"id,attr"`
    Author   Author
    Title    string
    Subtitle string   `xml:",omitempty"`
    Tags     []string `xml:"Tags>Tag"`
    Body     string   `xml:"Content"`
    Notes    string   `xml:",comment"`
}

func encode(w io.Writer) {
    post := BlogPost{
        Id: 10,
        Author: Author{
            Id: 5,
            Name: Name{
                First: "Alan",
                Last:  "Kay",
            },
        },
        Title: "It's All About Messages",
        Tags:  []string{"object-oriented", "programming", "oop"},
        Body:  "It's not about objects, it's about messages",
        Notes: "He's the boss",
    }

    encoder := xml.NewEncoder(w)
    err := encoder.Encode(post)
    if err != nil {
        log.Fatalf("failed encoding to a stream: %s", err)
    }
}

func decode(r io.Reader) {
    var post BlogPost
    decoder := xml.NewDecoder(r)
    err := decoder.Decode(&post)
    if err != nil {
        log.Fatalf("failed decoding from stream: %s", err)
    }
    log.Printf("%#v", post)
}

func pretty() {
    post := BlogPost{
        Id: 5,
        Author: Author{
            Id: 2,
            Name: Name{
                First: "Daniel",
                Last:  "Huckstep",
            },
        },
        Title: "Go, The Standard Library",
        Tags:  []string{"golang", "programming", "reference"},
        Body:  "I <strong>like</strong> programming Go, it's so much fun!",
        Notes: "Need to write more often...",
    }
    data, err := xml.MarshalIndent(post, "", "\t")
    if err != nil {
        log.Fatalf("failed pretty printing: %s", err)
    }
    log.Printf("pretty print:%s", data)
}

func tokens() {
    doc := []byte(`<post id="5"><title>Batman</title><author>Daniel Huckstep</author></post>`)
    decoder := xml.NewDecoder(bytes.NewReader(doc))
    for {
        token, err := decoder.Token()
        switch err {
        case nil:
            // Nothing to see here
        case io.EOF:
            log.Println("done parsing tokens")
            return
        default:
            log.Fatalf("got error getting token: %s", err)
        }

        switch tok := token.(type) {
        case xml.StartElement:
            log.Printf("found start element: %s", tok.Name)
        case xml.EndElement:
            log.Printf("found end element: %s", tok.Name)
        case xml.CharData:
            log.Printf("found chardata element: %s", tok)
        case xml.Comment:
            log.Printf("found comment element: %s", tok)
        case xml.ProcInst:
            log.Printf("found processing instruction: %s", tok.Target)
        case xml.Directive:
            log.Printf("found directive: %s", tok)
        default:
            panic("not reached")
        }
    }
}

func main() {
    pretty()
    var buffer bytes.Buffer
    encode(&buffer)
    log.Printf("encoded post to %s", buffer.String())
    decode(&buffer)
    tokens()
}
