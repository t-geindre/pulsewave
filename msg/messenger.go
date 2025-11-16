package msg

type Handler interface {
	HandleMessage(msg Message)
}

// Messenger handles incoming and outgoing message queues.
// Designed to allow multiple components to share the same queues.
type Messenger struct {
	in, out   *Queue
	handlers  []Handler
	drainRate int
}

// NewMessenger creates a new Messenger
// drainRate is the maximum number of messages to process per Process call.
func NewMessenger(in, out *Queue, drainRate int) *Messenger {
	return &Messenger{
		in:        in,
		out:       out,
		drainRate: drainRate,
	}
}

// Process drains incoming messages and dispatches them to registered handlers.
func (m *Messenger) Process() {
	m.in.Drain(m.drainRate, func(msg Message) {
		for _, h := range m.handlers {
			h.HandleMessage(msg)
		}
	})
}

// RegisterHandler registers a Handler function for a specific message kind.
func (m *Messenger) RegisterHandler(h Handler) {
	m.handlers = append(m.handlers, h)
}

// SendMessage sends a message to the outgoing queue.
func (m *Messenger) SendMessage(msg Message) {
	m.out.TryWrite(msg)
}
