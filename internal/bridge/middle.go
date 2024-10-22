package bridge

var (
	RawDataChannel  = make(chan Packet, 10)
	StringValueChan = make(chan string) //is used to send string type data back to the UI
	IntValueChan    = make(chan int)    //is used to send int type data back to the UI
	//WaitChan            = make(chan bool)    //is used to signal the frontend, that it can continue
	HeartRateLogChannel = make(chan []HRObj) //to transmit the parsed heart rate log
)
