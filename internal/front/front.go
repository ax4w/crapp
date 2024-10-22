package front

import (
	hrHitory "crapp/internal/front/hrHistory"
	"crapp/internal/middle"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

func battery(d []byte) {
	fmt.Printf("%d%% battery left\n", int(d[0]))
	middle.WaitChan <- true
}

func heartRate(d []byte) {
	if int(d[2]) > 0 {
		middle.HeartRateChannel <- int(d[2])
	}
}

func hrLogSettings(d []byte) {
	println(d[0], "enabled", int(math.Abs(float64(int(d[1])-2))), "interval", d[2])
	middle.WaitChan <- true
	//middle.HeartRateChannel <- true
}

func heartRateHistory(h *hrHitory.History) {
	h.AddTimes()
	//middle.WaitChan <- true
}

func timeResponse(d []byte) {
	println("time response is:")
	println("supports spO2", d[3]&2 != 0)
	middle.WaitChan <- true
}

func Run() {
	var hrHis *hrHitory.History
	for v := range middle.BackToFront {
		switch v.Id {
		case middle.CMD_BATTERY:
			battery(v.Data)
		case middle.CMD_START_HEART_RATE:
			heartRate(v.Data)
		case middle.CMD_DONE:
			return
		case middle.CMD_HEART_RATE_LOG_SETTINGS:
			hrLogSettings(v.Data)
		case middle.CMD_READ_HEART_RATE:
			sub := int(v.Data[0])
			if sub == 255 {
				panic("ahhh")
			}
			if sub == 23 && hrHis.IsToday() {
				println("done")
				heartRateHistory(hrHis)
			}
			if sub == 0 {
				hrHis = hrHitory.New(false, int(v.Data[1]), int(v.Data[2]))
			} else if sub == 1 {
				for _, k := range v.Data[5 : len(v.Data)-1] {
					hrHis.Heart_rates = append(hrHis.Heart_rates, int(k))
				}
				offset := 1
				ts := binary.LittleEndian.Uint32(v.Data[offset : offset+4])
				hrHis.T = time.Unix(int64(ts), 0)
				println("time is", hrHis.T.String())
				for _, v := range v.Data[5:len(v.Data)] {
					hrHis.Heart_rates = append(hrHis.Heart_rates, int(v))
				}
				//hrHis.Index += 9
			} else {
				for _, k := range v.Data[1:14] {
					hrHis.Heart_rates = append(hrHis.Heart_rates, int(k))
				}
				hrHis.Index += 13
				if sub == hrHis.Size-1 {
					println("done")
					heartRateHistory(hrHis)
				}
			}
		case middle.CMD_SET_TIME:
			timeResponse(v.Data)
		default:

		}
		//println("Got paket type", v.Id, "with data", v.Data[0])
		//middle.FrontToBack <- middle.Packet{Id: middle.CMD_DONE}
	}
}
