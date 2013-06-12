package client

import (
	"fmt"
	"github.com/zond/god/common"
	"github.com/zond/god/setop"
)

func ExampleSetExpression() {
	conn := MustConn("127.0.0.1:9191")
	conn.Kill()
	conn.SubPut([]byte("myfriends"), []byte("alice"), common.EncodeFloat64(10))
	conn.SubPut([]byte("myfriends"), []byte("bob"), common.EncodeFloat64(5))
	conn.SubPut([]byte("yourfriends"), []byte("bob"), common.EncodeFloat64(6))
	conn.SubPut([]byte("yourfriends"), []byte("charlie"), common.EncodeFloat64(4))
	fmt.Printf("name score\n")
	for _, friend := range conn.SetExpression(setop.SetExpression{
		Code: "(U:FloatSum myfriends yourfriends)",
	}) {
		fmt.Printf("%v %v\n", string(friend.Key), common.MustDecodeFloat64(friend.Values[0]))
	}
	// Output:
	// name score
	// alice 10
	// bob 11
	// charlie 4
}

func ExampleTreeMirror() {
	conn := MustConn("127.0.0.1:9191")
	conn.Kill()
	conn.SubAddConfiguration([]byte("myfriends"), "mirrored", "yes")
	conn.SubPut([]byte("myfriends"), []byte("alice"), common.EncodeFloat64(10))
	conn.SubPut([]byte("myfriends"), []byte("bob"), common.EncodeFloat64(5))
	conn.SubPut([]byte("myfriends"), []byte("charlie"), common.EncodeFloat64(6))
	conn.SubPut([]byte("myfriends"), []byte("denise"), common.EncodeFloat64(4))
	fmt.Printf("name score\n")
	for _, friend := range conn.MirrorReverseSlice([]byte("myfriends"), nil, nil, true, true) {
		fmt.Printf("%v %v\n", common.MustDecodeFloat64(friend.Key), string(friend.Value))
	}
	// Output:
	// name score
	// 10 alice
	// 6 charlie
	// 5 bob
	// 4 denise
}
