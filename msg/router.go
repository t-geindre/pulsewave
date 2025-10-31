package msg

import (
	"time"

	"github.com/rs/zerolog"
)

// Todo queue id maybe useless since we can probably route user pointers to queues directly
type route struct {
	in, out int
	src     Source
	kind    Kind
}

type Router struct {
	inputs  []*Queue
	outputs []*Queue
	routes  []route
	logger  zerolog.Logger
}

func NewRouter(log zerolog.Logger) *Router {
	return &Router{
		inputs:  make([]*Queue, 0),
		outputs: make([]*Queue, 0),
		routes:  make([]route, 0),
		logger:  log.With().Str("component", "router").Logger(),
	}
}

func (r *Router) AddInput(cap int) (*Queue, int) {
	r.inputs = append(r.inputs, NewQueue(cap))
	idx := len(r.inputs) - 1

	return r.inputs[idx], idx
}

func (r *Router) AddOutput(cap int) (*Queue, int) {
	r.outputs = append(r.outputs, NewQueue(cap))
	idx := len(r.outputs) - 1

	return r.outputs[idx], idx
}

func (r *Router) AddRoute(in, out int, src Source, kind Kind) {
	r.routes = append(r.routes, route{
		in:   in,
		out:  out,
		src:  src,
		kind: kind,
	})
}

// Route processes messages from inputs to outputs based on defined routes.
// blocking call; should be run in its own goroutine.
func (r *Router) Route() {
	for {
		if !r.doRoute() {
			time.Sleep(time.Millisecond)
		}
	}
}

func (r *Router) doRoute() bool {
	var msg Message
	var handled bool

	for inK, input := range r.inputs {
		if !input.TryRead(&msg) {
			continue
		}

		handled = true

		for _, rt := range r.routes {
			if rt.in != inK {
				continue
			}
			if rt.src != msg.Source {
				continue
			}
			if rt.kind != msg.Kind {
				continue
			}

			output := r.outputs[rt.out]
			if !output.TryWrite(msg) {
				r.logger.Warn().Int("out", rt.out).Msg("message dropped")
			}
		}
	}

	return handled
}
