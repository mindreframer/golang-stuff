package goson

//put qoutes around the input string
func quote(s []byte) []byte {
	q := make([]byte, len(s)+2)
	q[0] = '"'
	q[len(q)-1] = '"'
	for i := 1; i < len(q)-1; i++ {
		q[i] = s[i-1]
	}
	return q
}
