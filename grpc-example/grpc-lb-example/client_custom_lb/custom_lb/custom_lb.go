package custom_lb

import (
	"log"
	"math/rand"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const Name = "custom_robin"

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &randomPickerBuilder{}, base.Config{HealthCheck: false})
}

func init() {
	log.Println("CustomLB: Registering balancer builder with name:", Name)
	balancer.Register(newBuilder())
}

type randomPickerBuilder struct{}

func (pb *randomPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	log.Printf("CustomLB (randomPickerBuilder): Build called. Have %d Ready SubConns.", len(info.ReadySCs))

	if len(info.ReadySCs) == 0 {
		log.Println("CustomLB (randomPickerBuilder): No ready SubConns. Returning ErrNoSubConnAvailable.")
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	// Create a slice of ready SubConns
	var scs []balancer.SubConn
	for sc := range info.ReadySCs { // info.ReadySCs is a map[balancer.SubConn]balancer.SubConnInfo
		scs = append(scs, sc)
	}

	return &randomPicker{
		subConns: scs,
		// mu:       sync.Mutex{}, // Mutex might be needed if picker state is modified after creation
	}
}

// randomPicker is a Picker that randomly selects a SubConn from the available list.
type randomPicker struct {
	subConns []balancer.SubConn
	mu       sync.Mutex // To protect access to subConns if it could be modified concurrently
	// For this simple picker where subConns is set at creation, it might not be strictly necessary
	// for Pick, but good practice if state can change.
}

// Pick implements the Picker interface.
// It selects a SubConn from the list of ready SubConns.
func (p *randomPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.subConns) == 0 {
		log.Println("CustomLB (randomPicker): Pick - No SubConns available. Returning ErrNoSubConnAvailable.")
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	// Select a SubConn at random.
	idx := rand.Intn(len(p.subConns))
	sc := p.subConns[idx]
	log.Printf("CustomLB (randomPicker): Picked SubConn at index %d.", idx)

	// Return the picked SubConn.
	// The Done func in PickResult can be used for per-RPC stats or cleanup, nil if not needed.
	return balancer.PickResult{SubConn: sc, Done: nil}, nil
}
