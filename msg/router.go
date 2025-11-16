package msg

import (
	"time"

	"github.com/rs/zerolog"
)

type route struct {
	in, out *Queue
	kind    Kind
}

type Router struct {
	inputs []*Queue
	routes []route
	logger zerolog.Logger
}

func NewRouter(log zerolog.Logger) *Router {
	return &Router{
		inputs: make([]*Queue, 0),
		routes: make([]route, 0),
		logger: log,
	}
}

func (r *Router) AddInput(cap int) *Queue {
	q := NewQueue(cap)
	r.inputs = append(r.inputs, q)
	return q
}

func (r *Router) AddOutput(cap int) *Queue {
	return NewQueue(cap)
}

// AddRoute adds a routing rule from an input queue and message kind to an output queue.
// If in is nil, the route applies to all inputs.
func (r *Router) AddRoute(in *Queue, kind Kind, out *Queue) {
	r.routes = append(r.routes, route{
		in:   in,
		kind: kind,
		out:  out,
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
	var handled bool

	for _, input := range r.inputs {
		input.Drain(0, func(msg Message) {
			handled = true

			for _, rt := range r.routes {
				if rt.in != nil && rt.in != input {
					continue
				}
				if rt.kind != msg.Kind {
					continue
				}

				rt.out.TryWrite(msg)
			}
		})
	}

	return handled
}
