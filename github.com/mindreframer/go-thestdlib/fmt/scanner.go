package main

import (
    "fmt"
)

type Tuple struct {
    Left, Right int
}

func (t Tuple) Format(f fmt.State, c rune) {
    switch c {
    case 'P':
        fmt.Fprintf(f, "(%#v, %#v)", t.Left, t.Right)
    }
}

func (t *Tuple) Scan(state fmt.ScanState, verb rune) error {
    switch verb {
    case 'P':
        n, err := fmt.Fscanf(state, "(%d, %d)", &t.Left, &t.Right)
        if err != nil {
            return err
        }
        if n != 2 {
            return fmt.Errorf("scanned %d things, expected 2", n)
        }
    }
    return nil
}

func main() {
    var i int
    var f float32
    var t Tuple

    fmt.Printf("%d %P %f\n", i, t, f)
    fmt.Sscanf("5 (1, 2) 2.5", "%d %P %f", &i, &t, &f)
    fmt.Printf("%d %P %f\n", i, t, f)
}
