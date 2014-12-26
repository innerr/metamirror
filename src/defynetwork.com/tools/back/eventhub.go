package tools

func (p *EventHub) Register(name string, handler EventHandler) {
	handlers, ok := p.handlers[name]
	if !ok {
		handlers = []EventHandler{handler}
	} else {
		handlers = append(handlers, handler)
	}
	p.handlers[name] = handlers
}

func (p *EventHub) Event(name string, val interface{}) {
	handlers, ok := p.handlers[name]
	if !ok {
		if p.def != nil {
			p.def(name, val)
		}
		return
	}
	for _, it := range handlers {
		it(name, val)
	}
}

func NewEventHub(def EventHandler) *EventHub {
	return &EventHub{def, make(map[string][]EventHandler)}
}

type EventHub struct {
	def EventHandler
	handlers map[string][]EventHandler
}

type EventHandler func(name string, val interface{})
