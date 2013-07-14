package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/zond/god/client"
	"github.com/zond/god/common"
	"github.com/zond/setop"
	"io"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	stringFormat = "string"
	floatFormat  = "float"
	intFormat    = "int"
	bigFormat    = "big"
)

var formats = []string{stringFormat, floatFormat, intFormat, bigFormat}

type action func(conn *client.Conn, args []string)

var ip = flag.String("ip", "127.0.0.1", "IP address to connect to")
var port = flag.Int("port", 9191, "Port to connect to")
var enc = flag.String("enc", stringFormat, fmt.Sprintf("What format to assume when encoding and decoding byte slices: %v", formats))

func encode(s string) []byte {
	switch *enc {
	case stringFormat:
		return []byte(s)
	case floatFormat:
		result, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return common.EncodeFloat64(result)
	case intFormat:
		result, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
		return common.EncodeInt64(result)
	case bigFormat:
		result, ok := new(big.Int).SetString(s, 10)
		if !ok {
			panic(fmt.Errorf("Bad BigInt format: %v", s))
		}
		return common.EncodeBigInt(result)
	}
	panic(fmt.Errorf("Unknown encoding: %v", *enc))
}
func decode(b []byte) string {
	switch *enc {
	case stringFormat:
		return string(b)
	case floatFormat:
		res, err := common.DecodeFloat64(b)
		if err != nil {
			return fmt.Sprint(b)
		}
		return fmt.Sprint(res)
	case intFormat:
		res, err := common.DecodeInt64(b)
		if err != nil {
			return fmt.Sprint(b)
		}
		return fmt.Sprint(res)
	case bigFormat:
		return fmt.Sprint(common.DecodeBigInt(b))
	}
	panic(fmt.Errorf("Unknown encoding: %v", *enc))
}

type actionSpec struct {
	cmd  string
	args []*regexp.Regexp
}

func newActionSpec(pattern string) (result *actionSpec) {
	result = &actionSpec{}
	parts := strings.Split(pattern, " ")
	result.cmd = parts[0]
	for _, r := range parts[1:] {
		result.args = append(result.args, regexp.MustCompile(r))
	}
	return
}

var actions = map[*actionSpec]action{
	newActionSpec("mirrorReverseSliceIndex \\S+ \\d+ \\d+"): mirrorReverseSliceIndex,
	newActionSpec("mirrorSliceIndex \\S+ \\d+ \\d+"):        mirrorSliceIndex,
	newActionSpec("mirrorReverseSlice \\S+ \\S+ \\S+"):      mirrorReverseSlice,
	newActionSpec("mirrorSlice \\S+ \\S+ \\S+"):             mirrorSlice,
	newActionSpec("mirrorSliceLen \\S+ \\S+ \\d+"):          mirrorSliceLen,
	newActionSpec("mirrorReverseSliceLen \\S+ \\S+ \\d+"):   mirrorReverseSliceLen,
	newActionSpec("reverseSliceIndex \\S+ \\d+ \\d+"):       reverseSliceIndex,
	newActionSpec("sliceIndex \\S+ \\d+ \\d+"):              sliceIndex,
	newActionSpec("reverseSlice \\S+ \\S+ \\S+"):            reverseSlice,
	newActionSpec("slice \\S+ \\S+ \\S+"):                   slice,
	newActionSpec("sliceLen \\S+ \\S+ \\d+"):                sliceLen,
	newActionSpec("reverseSliceLen \\S+ \\S+ \\d+"):         reverseSliceLen,
	newActionSpec("setOp .+"):                               setOp,
	newActionSpec("dumpSetOp \\S+ .+"):                      dumpSetOp,
	newActionSpec("put \\S+ \\S+"):                          put,
	newActionSpec("clear"):                                  clear,
	newActionSpec("dump"):                                   dump,
	newActionSpec("subDump \\S+"):                           subDump,
	newActionSpec("subSize \\S+"):                           subSize,
	newActionSpec("size"):                                   size,
	newActionSpec("count \\S+ \\S+ \\S+"):                   count,
	newActionSpec("mirrorCount \\S+ \\S+ \\S+"):             mirrorCount,
	newActionSpec("get \\S+"):                               get,
	newActionSpec("del \\S+"):                               del,
	newActionSpec("subPut \\S+ \\S+ \\S+"):                  subPut,
	newActionSpec("subGet \\S+ \\S+"):                       subGet,
	newActionSpec("subDel \\S+ \\S+"):                       subDel,
	newActionSpec("subClear \\S+"):                          subClear,
	newActionSpec("describeAll"):                            describeAll,
	newActionSpec("describe \\S+"):                          describe,
	newActionSpec("describeTree \\S+"):                      describeTree,
	newActionSpec("describeAllTrees"):                       describeAllTrees,
	newActionSpec("mirrorFirst \\S+"):                       mirrorFirst,
	newActionSpec("mirrorLast \\S+"):                        mirrorLast,
	newActionSpec("mirrorPrevIndex \\S+ \\d+"):              mirrorPrevIndex,
	newActionSpec("mirrorNextIndex \\S+ \\d+"):              mirrorNextIndex,
	newActionSpec("first \\S+"):                             first,
	newActionSpec("last \\S+"):                              last,
	newActionSpec("prevIndex \\S+ \\d+"):                    prevIndex,
	newActionSpec("nextIndex \\S+ \\d+"):                    nextIndex,
	newActionSpec("next \\S+"):                              next,
	newActionSpec("prev \\S+"):                              prev,
	newActionSpec("subMirrorNext \\S+ \\S+"):                subMirrorNext,
	newActionSpec("subMirrorPrev \\S+ \\S+"):                subMirrorPrev,
	newActionSpec("mirrorIndexOf \\S+ \\S+"):                mirrorIndexOf,
	newActionSpec("mirrorReverseIndexOf \\S+ \\S+"):         mirrorReverseIndexOf,
	newActionSpec("subNext \\S+ \\S+"):                      subNext,
	newActionSpec("subPrev \\S+ \\S+"):                      subPrev,
	newActionSpec("indexOf \\S+ \\S+"):                      indexOf,
	newActionSpec("reverseIndexOf \\S+ \\S+"):               reverseIndexOf,
	newActionSpec("configuration"):                          configuration,
	newActionSpec("subConfiguration \\S+"):                  subConfiguration,
	newActionSpec("configure \\S+ \\S+"):                    configure,
	newActionSpec("subConfigure \\S+ \\S+ \\S+"):            subConfigure,
}

func mustAtoi(s string) *int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return &i
}

func configuration(conn *client.Conn, args []string) {
	fmt.Println(conn.Configuration())
}

func subConfiguration(conn *client.Conn, args []string) {
	fmt.Println(conn.SubConfiguration([]byte(args[1])))
}

func configure(conn *client.Conn, args []string) {
	conn.AddConfiguration(args[1], args[2])
}

func subConfigure(conn *client.Conn, args []string) {
	conn.SubAddConfiguration([]byte(args[1]), args[2], args[3])
}

func subSize(conn *client.Conn, args []string) {
	fmt.Println(conn.SubSize([]byte(args[1])))
}

func size(conn *client.Conn, args []string) {
	fmt.Println(conn.Size())
}

func mirrorReverseSliceIndex(conn *client.Conn, args []string) {
	for _, item := range conn.MirrorReverseSliceIndex([]byte(args[1]), mustAtoi(args[2]), mustAtoi(args[3])) {
		fmt.Printf("%v: %v => %v\n", item.Index, decode(item.Key), string(item.Value))
	}
}

func mirrorSliceIndex(conn *client.Conn, args []string) {
	for _, item := range conn.MirrorSliceIndex([]byte(args[1]), mustAtoi(args[2]), mustAtoi(args[3])) {
		fmt.Printf("%v: %v => %v\n", item.Index, decode(item.Key), string(item.Value))
	}
}

func mirrorReverseSlice(conn *client.Conn, args []string) {
	for i, item := range conn.MirrorReverseSlice([]byte(args[1]), []byte(args[2]), []byte(args[3]), true, false) {
		fmt.Printf("%v: %v => %v\n", i, decode(item.Key), string(item.Value))
	}
}

func mirrorSlice(conn *client.Conn, args []string) {
	for i, item := range conn.MirrorSlice([]byte(args[1]), []byte(args[2]), []byte(args[3]), true, false) {
		fmt.Printf("%v: %v => %v\n", i, decode(item.Key), string(item.Value))
	}
}

func mirrorSliceLen(conn *client.Conn, args []string) {
	for _, item := range conn.MirrorSliceLen([]byte(args[1]), []byte(args[2]), true, *(mustAtoi(args[3]))) {
		fmt.Printf("%v => %v\n", decode(item.Key), string(item.Value))
	}
}

func mirrorReverseSliceLen(conn *client.Conn, args []string) {
	for _, item := range conn.MirrorReverseSliceLen([]byte(args[1]), []byte(args[2]), true, *(mustAtoi(args[3]))) {
		fmt.Printf("%v => %v\n", decode(item.Key), string(item.Value))
	}
}

func reverseSliceIndex(conn *client.Conn, args []string) {
	for _, item := range conn.ReverseSliceIndex([]byte(args[1]), mustAtoi(args[2]), mustAtoi(args[3])) {
		fmt.Printf("%v: %v => %v\n", item.Index, string(item.Key), decode(item.Value))
	}
}

func sliceIndex(conn *client.Conn, args []string) {
	for _, item := range conn.SliceIndex([]byte(args[1]), mustAtoi(args[2]), mustAtoi(args[3])) {
		fmt.Printf("%v: %v => %v\n", item.Index, string(item.Key), decode(item.Value))
	}
}

func reverseSlice(conn *client.Conn, args []string) {
	for i, item := range conn.ReverseSlice([]byte(args[1]), []byte(args[2]), []byte(args[3]), true, false) {
		fmt.Printf("%v: %v => %v\n", i, string(item.Key), decode(item.Value))
	}
}

func slice(conn *client.Conn, args []string) {
	for i, item := range conn.Slice([]byte(args[1]), []byte(args[2]), []byte(args[3]), true, false) {
		fmt.Printf("%v: %v => %v\n", i, string(item.Key), decode(item.Value))
	}
}

func sliceLen(conn *client.Conn, args []string) {
	for _, item := range conn.SliceLen([]byte(args[1]), []byte(args[2]), true, *(mustAtoi(args[3]))) {
		fmt.Printf("%v => %v\n", string(item.Key), decode(item.Value))
	}
}

func reverseSliceLen(conn *client.Conn, args []string) {
	for _, item := range conn.ReverseSliceLen([]byte(args[1]), []byte(args[2]), true, *(mustAtoi(args[3]))) {
		fmt.Printf("%v => %v\n", string(item.Key), decode(item.Value))
	}
}

func printSetOpRes(res setop.SetOpResult) {
	var vals []string
	for _, val := range res.Values {
		vals = append(vals, decode(val))
	}
	fmt.Printf("%v => %v\n", string(res.Key), vals)
}

func setOp(conn *client.Conn, args []string) {
	op, err := setop.NewSetOpParser(args[1]).Parse()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, res := range conn.SetExpression(setop.SetExpression{Op: op}) {
			printSetOpRes(res)
		}
	}
}

func dumpSetOp(conn *client.Conn, args []string) {
	op, err := setop.NewSetOpParser(args[2]).Parse()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, res := range conn.SetExpression(setop.SetExpression{Dest: []byte(args[1]), Op: op}) {
			printSetOpRes(res)
		}
	}
}

func mirrorReverseIndexOf(conn *client.Conn, args []string) {
	if index, existed := conn.MirrorReverseIndexOf([]byte(args[1]), []byte(args[2])); existed {
		fmt.Println(index)
	}
}

func mirrorIndexOf(conn *client.Conn, args []string) {
	if index, existed := conn.MirrorIndexOf([]byte(args[1]), []byte(args[2])); existed {
		fmt.Println(index)
	}
}

func reverseIndexOf(conn *client.Conn, args []string) {
	if index, existed := conn.ReverseIndexOf([]byte(args[1]), []byte(args[2])); existed {
		fmt.Println(index)
	}
}

func indexOf(conn *client.Conn, args []string) {
	if index, existed := conn.IndexOf([]byte(args[1]), []byte(args[2])); existed {
		fmt.Println(index)
	}
}

func show(conn *client.Conn) {
	fmt.Println(conn.Describe())
}

func mirrorCount(conn *client.Conn, args []string) {
	fmt.Println(conn.MirrorCount([]byte(args[1]), []byte(args[2]), []byte(args[3]), true, false))
}

func count(conn *client.Conn, args []string) {
	fmt.Println(conn.Count([]byte(args[1]), []byte(args[2]), []byte(args[3]), true, false))
}

func mirrorPrevIndex(conn *client.Conn, args []string) {
	if key, value, index, existed := conn.MirrorPrevIndex([]byte(args[1]), *(mustAtoi(args[2]))); existed {
		fmt.Printf("%v: %v => %v\n", index, decode(key), string(value))
	}
}

func mirrorNextIndex(conn *client.Conn, args []string) {
	if key, value, index, existed := conn.MirrorNextIndex([]byte(args[1]), *(mustAtoi(args[2]))); existed {
		fmt.Printf("%v: %v => %v\n", index, decode(key), string(value))
	}
}

func prevIndex(conn *client.Conn, args []string) {
	if key, value, index, existed := conn.PrevIndex([]byte(args[1]), *(mustAtoi(args[2]))); existed {
		fmt.Printf("%v: %v => %v\n", index, string(key), decode(value))
	}
}

func nextIndex(conn *client.Conn, args []string) {
	if key, value, index, existed := conn.NextIndex([]byte(args[1]), *(mustAtoi(args[2]))); existed {
		fmt.Printf("%v: %v => %v\n", index, string(key), decode(value))
	}
}

func prev(conn *client.Conn, args []string) {
	if key, value, existed := conn.Prev([]byte(args[1])); existed {
		fmt.Printf("%v => %v\n", string(key), decode(value))
	}
}

func next(conn *client.Conn, args []string) {
	if key, value, existed := conn.Next([]byte(args[1])); existed {
		fmt.Printf("%v => %v\n", string(key), decode(value))
	}
}

func mirrorFirst(conn *client.Conn, args []string) {
	if key, value, existed := conn.MirrorFirst([]byte(args[1])); existed {
		fmt.Println(decode(key), "=>", string(value))
	}
}

func mirrorLast(conn *client.Conn, args []string) {
	if key, value, existed := conn.MirrorLast([]byte(args[1])); existed {
		fmt.Println(decode(key), "=>", string(value))
	}
}

func first(conn *client.Conn, args []string) {
	if key, value, existed := conn.First([]byte(args[1])); existed {
		fmt.Println(string(key), "=>", decode(value))
	}
}

func last(conn *client.Conn, args []string) {
	if key, value, existed := conn.Last([]byte(args[1])); existed {
		fmt.Println(string(key), "=>", decode(value))
	}
}

func subMirrorNext(conn *client.Conn, args []string) {
	if key, value, existed := conn.SubMirrorNext([]byte(args[1]), []byte(args[2])); existed {
		fmt.Printf("%v => %v\n", decode(key), string(value))
	}
}

func subMirrorPrev(conn *client.Conn, args []string) {
	if key, value, existed := conn.SubMirrorPrev([]byte(args[1]), []byte(args[2])); existed {
		fmt.Printf("%v => %v\n", decode(key), string(value))
	}
}

func subNext(conn *client.Conn, args []string) {
	if key, value, existed := conn.SubNext([]byte(args[1]), []byte(args[2])); existed {
		fmt.Printf("%v => %v\n", string(key), decode(value))
	}
}

func subPrev(conn *client.Conn, args []string) {
	if key, value, existed := conn.SubPrev([]byte(args[1]), []byte(args[2])); existed {
		fmt.Printf("%v => %v\n", string(key), decode(value))
	}
}

func describeAll(conn *client.Conn, args []string) {
	for _, description := range conn.DescribeAllNodes() {
		fmt.Println(description.Describe())
	}
}

func describeAllTrees(conn *client.Conn, args []string) {
	fmt.Print(conn.DescribeAllTrees())
}

func describe(conn *client.Conn, args []string) {
	if bytes, err := hex.DecodeString(args[1]); err != nil {
		fmt.Println(err)
	} else {
		if result, err := conn.DescribeNode(bytes); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result.Describe())
		}
	}
}

func describeTree(conn *client.Conn, args []string) {
	if bytes, err := hex.DecodeString(args[1]); err != nil {
		fmt.Println(err)
	} else {
		if result, err := conn.DescribeTree(bytes); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	}
}

func get(conn *client.Conn, args []string) {
	if value, existed := conn.Get([]byte(args[1])); existed {
		fmt.Printf("%v\n", decode(value))
	}
}

func subGet(conn *client.Conn, args []string) {
	if value, existed := conn.SubGet([]byte(args[1]), []byte(args[2])); existed {
		fmt.Printf("%v\n", decode(value))
	}
}

func clear(conn *client.Conn, args []string) {
	conn.Clear()
}

func dump(conn *client.Conn, args []string) {
	dump, wait := conn.Dump()
	linedump(dump, wait)
}

func subDump(conn *client.Conn, args []string) {
	dump, wait := conn.SubDump([]byte(args[1]))
	linedump(dump, wait)
}

func linedump(dump chan [2][]byte, wait *sync.WaitGroup) {
	defer func() {
		close(dump)
		wait.Wait()
	}()
	reader := bufio.NewReader(os.Stdin)
	var pair []string
	var line string
	var err error
	for line, err = reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		pair = strings.Split(strings.TrimSpace(line), "=")
		if len(pair) == 2 {
			dump <- [2][]byte{[]byte(pair[0]), encode(pair[1])}
		} else {
			return
		}
	}
	if err != io.EOF {
		fmt.Println(err)
	}
}

func put(conn *client.Conn, args []string) {
	conn.Put([]byte(args[1]), encode(args[2]))
}

func subPut(conn *client.Conn, args []string) {
	conn.SubPut([]byte(args[1]), []byte(args[2]), encode(args[3]))
}

func subClear(conn *client.Conn, args []string) {
	conn.SubClear([]byte(args[1]))
}

func subDel(conn *client.Conn, args []string) {
	conn.SubDel([]byte(args[1]), []byte(args[2]))
}

func del(conn *client.Conn, args []string) {
	conn.Del([]byte(args[1]))
}

func main() {
	flag.Parse()
	conn := client.MustConn(fmt.Sprintf("%v:%v", *ip, *port))
	if len(flag.Args()) == 0 {
		show(conn)
	} else {
		for spec, fun := range actions {
			if spec.cmd == flag.Args()[0] {
				matchingParts := true
				for index, reg := range spec.args {
					if !reg.MatchString(flag.Args()[index+1]) {
						matchingParts = false
						break
					}
				}
				if matchingParts {
					fun(conn, flag.Args())
					return
				}
			}
		}
		fmt.Println("No command given?")
	}
}
