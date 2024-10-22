package middle

import "time"

const (
	CMD_DONE                    = 0
	CMD_BLINK_TWICE             = 16
	CMD_SET_TIME                = 1
	CMD_BATTERY                 = 3
	CMD_REAL_TIME_HEART_RATE    = 30
	CMD_START_HEART_RATE        = 105
	CMD_STOP_HEART_RATE         = 106
	CMD_HEART_RATE_LOG_SETTINGS = 22
	CMD_READ_HEART_RATE         = 21
)

type Packet struct {
	Id   int
	Data []byte
}


type HRObj struct {
	T  time.Time
	HR int
}
