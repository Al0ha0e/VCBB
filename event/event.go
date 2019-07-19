package event

type EventSystem struct {
	registerTable map[string][]chan string
}

func (this *EventSystem) Register(event string, ch chan string) {
	this.registerTable[event] = append(this.registerTable[event], ch)
}

func (this *EventSystem) Emit(event string, msg string) {
	for _, ch := range this.registerTable[event] {
		ch <- msg
	}
}
