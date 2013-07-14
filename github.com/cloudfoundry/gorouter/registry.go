package router

import (
	"encoding/json"
	"fmt"
	mbus "github.com/cloudfoundry/go_cfmessagebus"
	steno "github.com/cloudfoundry/gosteno"
	"math/rand"
	"github.com/cloudfoundry/gorouter/stats"
	"github.com/cloudfoundry/gorouter/util"
	"strings"
	"sync"
	"time"
)

type Uri string
type Uris []Uri

func (u Uri) ToLower() Uri {
	return Uri(strings.ToLower(string(u)))
}

func (ms Uris) Sub(ns Uris) Uris {
	var rs Uris

	for _, m := range ms {
		found := false
		for _, n := range ns {
			if m == n {
				found = true
				break
			}
		}

		if !found {
			rs = append(rs, m)
		}
	}

	return rs
}

func (x Uris) Has(y Uri) bool {
	for _, xb := range x {
		if xb == y {
			return true
		}
	}

	return false
}

func (x Uris) Remove(y Uri) (Uris, bool) {
	for i, xb := range x {
		if xb == y {
			x[i] = x[len(x)-1]
			x = x[:len(x)-1]
			return x, true
		}
	}

	return x, false
}

type BackendId string

type Backend struct {
	sync.Mutex

	*steno.Logger

	BackendId BackendId

	ApplicationId     string
	Host              string
	Port              uint16
	Tags              map[string]string
	PrivateInstanceId string

	U          Uris
	updated_at time.Time
}

func (b *Backend) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.CanonicalAddr())
}

func newBackend(i BackendId, m *registryMessage, l *steno.Logger) *Backend {
	b := &Backend{
		Logger: l,

		BackendId: i,

		ApplicationId:     m.App,
		Host:              m.Host,
		Port:              m.Port,
		Tags:              m.Tags,
		PrivateInstanceId: m.PrivateInstanceId,

		U:          make([]Uri, 0),
		updated_at: time.Now(),
	}

	return b
}

func (b *Backend) CanonicalAddr() string {
	return fmt.Sprintf("%s:%d", b.Host, b.Port)
}

func (b *Backend) ToLogData() interface{} {
	return struct {
		ApplicationId string
		Host          string
		Port          uint16
		Tags          map[string]string
	}{
		b.ApplicationId,
		b.Host,
		b.Port,
		b.Tags,
	}
}

func (b *Backend) register(u Uri) bool {
	if !b.U.Has(u) {
		b.Infof("Register %s (%s)", u, b.BackendId)
		b.U = append(b.U, u)
		return true
	}

	return false
}

func (b *Backend) unregister(u Uri) bool {
	x, ok := b.U.Remove(u)
	if ok {
		b.Infof("Unregister %s (%s)", u, b.BackendId)
		b.U = x
	}

	return ok
}

// This is a transient struct. It doesn't maintain state.
type registryMessage struct {
	Host string            `json:"host"`
	Port uint16            `json:"port"`
	Uris Uris              `json:"uris"`
	Tags map[string]string `json:"tags"`
	App  string            `json:"app"`

	PrivateInstanceId string `json:"private_instance_id"`
}

func (m registryMessage) BackendId() (b BackendId, ok bool) {
	if m.Host != "" && m.Port != 0 {
		b = BackendId(fmt.Sprintf("%s:%d", m.Host, m.Port))
		ok = true
	}

	return
}

type Registry struct {
	sync.RWMutex

	*steno.Logger

	*stats.ActiveApps
	*stats.TopApps

	byUri       map[Uri][]*Backend
	byBackendId map[BackendId]*Backend

	staleTracker *util.ListMap

	pruneStaleDropletsInterval time.Duration
	dropletStaleThreshold      time.Duration

	messageBus mbus.CFMessageBus
}

func NewRegistry(c *Config, messageBusClient mbus.CFMessageBus) *Registry {
	r := &Registry{
		messageBus: messageBusClient,
	}

	r.Logger = steno.NewLogger("router.registry")

	r.ActiveApps = stats.NewActiveApps()
	r.TopApps = stats.NewTopApps()

	r.byUri = make(map[Uri][]*Backend)
	r.byBackendId = make(map[BackendId]*Backend)

	r.staleTracker = util.NewListMap()

	r.pruneStaleDropletsInterval = c.PruneStaleDropletsInterval
	r.dropletStaleThreshold = c.DropletStaleThreshold

	return r
}

func (r *Registry) StartPruningCycle() {
	go r.checkAndPrune()
}

func (r *Registry) isStateStale() bool {
	return !r.messageBus.Ping()
}

func (r *Registry) NumUris() int {
	r.RLock()
	defer r.RUnlock()

	return len(r.byUri)
}

func (r *Registry) NumBackends() int {
	r.RLock()
	defer r.RUnlock()

	return len(r.byBackendId)
}

func (r *Registry) registerUri(b *Backend, u Uri) {
	u = u.ToLower()

	ok := b.register(u)
	if ok {
		x := r.byUri[u]
		r.byUri[u] = append(x, b)
	}
}

func (registry *Registry) Register(message *registryMessage) {
	i, ok := message.BackendId()
	if !ok || len(message.Uris) == 0 {
		return
	}

	registry.Lock()
	defer registry.Unlock()

	backend, ok := registry.byBackendId[i]
	if !ok {
		backend = newBackend(i, message, registry.Logger)
		registry.byBackendId[i] = backend
	}

	for _, uri := range message.Uris {
		registry.registerUri(backend, uri)
	}

	backend.updated_at = time.Now()

	registry.staleTracker.PushBack(backend)
}

func (r *Registry) unregisterUri(backend *Backend, uri Uri) {
	uri = uri.ToLower()

	ok := backend.unregister(uri)
	if ok {
		backends := r.byUri[uri]
		for i, b := range backends {
			if b == backend {
				// Remove b from list of backends
				backends[i] = backends[len(backends)-1]
				backends = backends[:len(backends)-1]
				break
			}
		}

		if len(backends) == 0 {
			delete(r.byUri, uri)
		} else {
			r.byUri[uri] = backends
		}
	}

	// Remove backend if it no longer has uris
	if len(backend.U) == 0 {
		delete(r.byBackendId, backend.BackendId)
		r.staleTracker.Delete(backend)
	}
}

func (r *Registry) Unregister(m *registryMessage) {
	i, ok := m.BackendId()
	if !ok {
		return
	}

	r.Lock()
	defer r.Unlock()

	b, ok := r.byBackendId[i]
	if !ok {
		return
	}

	for _, u := range m.Uris {
		r.unregisterUri(b, u)
	}
}

func (registry *Registry) pruneStaleDroplets() {
	for registry.staleTracker.Len() > 0 {
		backend := registry.staleTracker.Front().(*Backend)
		if backend.updated_at.Add(registry.dropletStaleThreshold).After(time.Now()) {
			log.Infof("Droplet is not stale; NOT pruning: %v", backend.BackendId)
			break
		}

		log.Infof("Pruning stale droplet: %v ", backend.BackendId)

		for _, uri := range backend.U {
			registry.unregisterUri(backend, uri)
		}
	}
}

func (registry *Registry) resetTracker() {
	for registry.staleTracker.Len() > 0 {
		registry.staleTracker.Delete(registry.staleTracker.Front().(*Backend))
	}
}

func (registry *Registry) PruneStaleDroplets() {
	if registry.isStateStale() {
		log.Info("State is stale; NOT pruning")
		registry.resetTracker()
		return
	}

	registry.Lock()
	defer registry.Unlock()

	registry.pruneStaleDroplets()
}

func (r *Registry) checkAndPrune() {
	if r.pruneStaleDropletsInterval == 0 {
		return
	}

	tick := time.Tick(r.pruneStaleDropletsInterval)
	for {
		select {
		case <-tick:
			log.Debug("Start to check and prune stale droplets")
			r.PruneStaleDroplets()
		}
	}
}

func (r *Registry) Lookup(host string) (*Backend, bool) {
	r.RLock()
	defer r.RUnlock()

	x, ok := r.byUri[Uri(host).ToLower()]
	if !ok {
		return nil, false
	}

	// Return random backend from slice of backends for the specified uri
	return x[rand.Intn(len(x))], true
}

func (r *Registry) LookupByPrivateInstanceId(host string, p string) (*Backend, bool) {
	r.RLock()
	defer r.RUnlock()

	x, ok := r.byUri[Uri(host).ToLower()]
	if !ok {
		return nil, false
	}

	for _, b := range x {
		if b.PrivateInstanceId == p {
			return b, true
		}
	}

	return nil, false
}

func (r *Registry) CaptureBackendRequest(x *Backend, t time.Time) {
	if x.ApplicationId != "" {
		r.ActiveApps.Mark(x.ApplicationId, t)
		r.TopApps.Mark(x.ApplicationId, t)
	}
}

func (r *Registry) MarshalJSON() ([]byte, error) {
	r.RLock()
	defer r.RUnlock()

	return json.Marshal(r.byUri)
}
