package bench

import (
	"flag"
	"fmt"
	"net/rpc"
	"regexp"
	"strconv"
	"strings"
)

func RunMaster() {
	ip := flag.String("ip", "127.0.0.1", "IP address to find a node at")
	slaves := flag.String("slaves", "", "Comma separated list of slave host:ip pairs")
	port := flag.Int("port", 9191, "Port to find a node at")
	maxKey := flag.Int64("maxKey", 10000, "Biggest key as int64 converted to []byte using common.EncodeInt64")
	prepare := flag.String("prepare", "0-0", "The key range (as int64's) to make sure exists before starging")
	keyRangePattern := regexp.MustCompile("^(\\d+)-(\\d+)$")
	flag.Parse()
	slavehosts := strings.Split(*slaves, ",")
	clients := make([]*rpc.Client, len(slavehosts))
	rps := make([]int64, len(slavehosts))
	var err error
	for index, addr := range slavehosts {
		if clients[index], err = rpc.Dial("tcp", addr); err != nil {
			panic(err)
		}
	}
	command := SpinCommand{
		Addr:   fmt.Sprintf("%v:%v", *ip, *port),
		MaxKey: *maxKey,
	}
	if match := keyRangePattern.FindStringSubmatch(*prepare); match != nil && match[1] != match[2] {
		from, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			panic(err)
		}
		to, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			panic(err)
		}
		calls := make([]*rpc.Call, len(clients))
		for index, client := range clients {
			calls[index] = client.Go("Slave.Prepare", PrepareCommand{
				Addr: command.Addr,
				Range: [2]int64{
					from + (int64(index) * ((to - from) / int64(len(clients)))),
					from + ((int64(index) + 1) * ((to - from) / int64(len(clients)))),
				},
			}, &Nothing{}, nil)
		}
		for _, call := range calls {
			<-call.Done
			if call.Error != nil {
				panic(call.Error)
			}
		}
	}
	var oldSpinRes *SpinResult
	var spinRes SpinResult
	for _, client := range clients {
		if err = client.Call("Slave.Spin", command, &spinRes); err != nil {
			panic(err)
		}
		if oldSpinRes == nil {
			oldSpinRes = &spinRes
		} else {
			if spinRes.Keys != oldSpinRes.Keys || spinRes.Nodes != oldSpinRes.Nodes {
				panic(fmt.Errorf("Last slave had %v nodes and %v keys, now I get %v nodes and %v keys?", oldSpinRes.Nodes, oldSpinRes.Keys, spinRes.Nodes, spinRes.Keys))
			}
		}
	}
	for _, client := range clients {
		if err = client.Call("Slave.Wait", Nothing{}, &Nothing{}); err != nil {
			panic(err)
		}
	}
	for index, client := range clients {
		if err = client.Call("Slave.Current", Nothing{}, &(rps[index])); err != nil {
			panic(err)
		}
	}
	for _, client := range clients {
		if err = client.Call("Slave.Stop", Nothing{}, &Nothing{}); err != nil {
			panic(err)
		}
	}
	sum := int64(0)
	for _, r := range rps {
		sum += r
	}
	fmt.Printf("%v\t%v\t%v\n", spinRes.Nodes, spinRes.Keys, sum)
}
