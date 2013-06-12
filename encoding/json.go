package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "os"
)

type BlogPost struct {
    // Marshal as "writer" instead of Author
    Author string `json:"writer,omitempty"`
    // Will get marshalled as "Title"
    Title string
    Body  string `json:"body"`
    // Don't marshal this field at all
    Published bool `json:"-"`
}

// This would marshal just fine,
// but let's write out own marshaller.
type Pair struct {
    X, Y int
}

func (p Pair) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf(`"%d|%d"`, p.X, p.Y)), nil
}

func (p *Pair) UnmarshalJSON(data []byte) error {
    _, err := fmt.Sscanf(string(data), `"%d|%d"`, &p.X, &p.Y)
    return err
}

func encodeTo(w io.Writer, i interface{}) {
    encoder := json.NewEncoder(w)
    if err := encoder.Encode(i); err != nil {
        log.Fatalf("failed encoding to writer: %s", err)
    }
}

func encode(i interface{}) []byte {
    data, err := json.Marshal(i)
    if err != nil {
        log.Fatalf("failed encoding: %s", data)
    }
    return data
}

func decode(data string) interface{} {
    var i interface{}
    err := json.Unmarshal([]byte(data), &i)
    if err != nil {
        log.Fatalf("failed decoding: %s", err)
    }
    return i
}

func simple() {
    log.Printf("encoded %d to %s", 1, encode(1))
    log.Printf("encoded %f to %s", 1.5, encode(1.5))
    log.Printf("encoded %s to %s", "Hello, World!", encode("Hello, World!"))

    log.Printf("decoded %f from %s", decode("1"), "1")
    log.Printf("decoded %v from %s", decode(`["foo","bar","baz"]`), `["foo","bar","baz"]`)
}

func custom() {
    pair := Pair{5, 10}
    encoded := encode(pair)
    log.Printf("encoded %v to %s", pair, encoded)

    var pair2 Pair
    if err := json.Unmarshal(encoded, &pair2); err != nil {
        log.Fatalf("failed decoding Pair: %s", err)
    }
    log.Printf("decoded %#v from %s", pair2, `"1|2"`)
}

func structExample() {
    post := BlogPost{
        // Since Author is empty, it won't be written out
        Title:     "Being Awesome At Go",
        Body:      "Read this book!",
        Published: true,
    }
    encodeTo(os.Stdout, post)

    post = BlogPost{
        Author:    "Daniel Huckstep",
        Title:     "Being Awesome At Go",
        Body:      "Read this book!",
        Published: true,
    }
    encodeTo(os.Stdout, post)
}

func streamDecode() {
    var buffer bytes.Buffer
    post := BlogPost{
        Author:    "Daniel Huckstep",
        Title:     "Being Awesome At Go",
        Body:      "Read this book!",
        Published: true,
    }
    encodeTo(&buffer, post)

    decoder := json.NewDecoder(&buffer)
    var newPost BlogPost
    if err := decoder.Decode(&newPost); err != nil {
        log.Printf("decoding failed: %s", err)
    }
    log.Printf("decoded %#v", newPost)
}

func pretty() {
    post := BlogPost{
        Author:    "Daniel Huckstep",
        Title:     "Being Awesome At Go",
        Body:      "Read this book!",
        Published: true,
    }
    data, err := json.MarshalIndent(post, "", "\t")
    if err != nil {
        log.Fatalf("failed marshal with indent: %s", err)
    }
    log.Printf("pretty print:\n%s", data)
}

func main() {
    simple()
    custom()
    structExample()
    streamDecode()
    pretty()
}
