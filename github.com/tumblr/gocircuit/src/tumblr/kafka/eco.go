// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafka

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	zookeeper "github.com/petar/gozk"
	zutil "circuit/kit/zookeeper/util"
)

// Ecosystem maintains a real-time record of the current Kafka broker
// ecosystem, as read from Zookeeper
type Ecosystem struct {
	lk        sync.Mutex
	zookeeper *zookeeper.Conn
	watch     *zutil.Watch
	stat      *zookeeper.Stat
	brokers   BrokersView
}

// NewEcosystem returns a new ecosystem instance using the supplied Zookeeper connection
func NewEcosystem(z *zookeeper.Conn) (*Ecosystem, error) {
	return &Ecosystem{
		zookeeper: z,
		watch:     zutil.InstallWatch(z, "/brokers/ids"),
	}, nil
}

// refresh updates the brokers view if necessary
func (e *Ecosystem) refresh() error {
	e.lk.Lock()
	defer e.lk.Unlock()
	if e.watch == nil {
		return errors.New("already closed")
	}
	children, stat, err := e.watch.Children()
	if err != nil {
		return err
	}
	if e.stat != nil && stat.Mzxid() <= e.stat.Mzxid() {
		return nil
	}
	e.stat = stat
	brokers := make(BrokersView, 0)
	for _, brokerID := range children {
		var err error
		b := &BrokerView{}
		b.BrokerID, err = ParseBrokerID(brokerID)
		if err != nil {
			continue
		}
		data, _, err := e.zookeeper.Get("/brokers/ids/" + brokerID)
		if err != nil {
			continue
		}
		s := strings.SplitN(data, ":", 2)
		if len(s) != 2 {
			continue
		}
		b.Creator, b.HostPort = s[0], s[1]
		brokers = append(brokers, b)
	}
	sort.Sort(brokers)
	e.brokers = brokers
	return nil
}

// Brokers returns the current brokers ecosystem
func (e *Ecosystem) Brokers() (brokers []*BrokerView, err error) {
	if err := e.refresh(); err != nil {
		return nil, err
	}

	e.lk.Lock()
	defer e.lk.Unlock()

	brokers = make([]*BrokerView, len(e.brokers))
	copy(brokers, e.brokers)
	return brokers, nil
}

// ChooseBroker chooses a uniformly random broker from the current broker ecosystem
func (e *Ecosystem) ChooseBroker() (*BrokerView, error) {
	brokers, err := e.Brokers()
	if err != nil {
		return nil, err
	}
	if len(brokers) == 0 {
		return nil, ErrNoBrokers
	}
	return brokers[rand.Intn(len(brokers))], nil
}

// Close closes the ecosystem instance and disconnects from Zookeeper
func (e *Ecosystem) Close() error {
	e.lk.Lock()
	defer e.lk.Unlock()

	if e.zookeeper == nil {
		return errors.New("kafka eco already closed")
	}
	e.zookeeper = nil

	err := e.watch.Close()
	e.watch = nil
	return err
}

// BrokerID represents a broker ID :)
type BrokerID int32

// String returns the textual representation of the broker ID, as used in the
// Zookeeper node name for this ID
func (x BrokerID) String() string {
	return strconv.Itoa(int(x))
}

// ParseBrokerID parses a broker ID from the string s, expecting the format produced by BrokerID.String
func ParseBrokerID(s string) (BrokerID, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return BrokerID(i), err
}

// BrokerView represents an active Kafka broker in the ecosystem
type BrokerView struct {
	BrokerID
	Creator  string
	HostPort string
}

// String returns a textual representation of the broker view
func (x *BrokerView) String() string {
	return fmt.Sprintf("broker(%d, %s, %s)", x.BrokerID, x.Creator, x.HostPort)
}

// BrokersView is a slice of BrokerView in ascending order of ID
type BrokersView []*BrokerView

// Len returns the length of the BrokerView slice
func (t BrokersView) Len() int {
	return len(t)
}

// Less compares the i-th and j-th elements of the BrokerView slice
func (t BrokersView) Less(i, j int) bool {
	return t[i].BrokerID < t[j].BrokerID
}

// Swap swaps the i-th and j-th elements of the BrokerView slice
func (t BrokersView) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
