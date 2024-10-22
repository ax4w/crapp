package middle

var (
	BackToFront         = make(chan Packet, 10)
	FrontToBack         = make(chan Packet, 10)
	HeartRateChannel    = make(chan int)
	WaitChan            = make(chan bool)
	HeartRateLogChannel = make(chan []HRObj)
)
