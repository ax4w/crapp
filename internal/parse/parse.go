package parse

import (
	"context"
	"crapp/internal/bridge"
	hrHitory "crapp/internal/parse/hrHistory"
	"encoding/binary"
	"fmt"
	"time"
)

func battery(d []byte) {
	bridge.StringValueChan <- fmt.Sprintf("%d%% battery left\n", int(d[0]))
}

func heartRate(d []byte) {
	if int(d[2]) > 0 {
		bridge.IntValueChan <- int(d[2])
	}
}

func hrLogSettings(d []byte) {
	bridge.StringValueChan <- fmt.Sprintf("heart rate log settings: enabled %v interval %d\n",
		int(d[1]) == 1, d[2])
}

func heartRateHistory(h *hrHitory.History) {
	h.AddTimes()
}

func timeResponse(d []byte) {
	// println("time response is:")
	// println("supports spO2", d[3]&2 != 0)
	bridge.StringValueChan <- "Setting the time was successful!"
}

func handle(hrHis *hrHitory.History, v bridge.Packet) *hrHitory.History {
	switch v.Id {
	case bridge.CMD_BATTERY:
		battery(v.Data)
	case bridge.CMD_START_HEART_RATE:
		heartRate(v.Data)
	case bridge.CMD_HEART_RATE_LOG_SETTINGS:
		hrLogSettings(v.Data)
	case bridge.CMD_READ_HEART_RATE:
		sub := int(v.Data[0])
		if sub == 255 {
			panic("error code recieved")
		}
		if sub == 23 && hrHis.IsToday() {
			println("done")
			heartRateHistory(hrHis)
		}
		if sub == 0 {
			hrHis = hrHitory.New(false, int(v.Data[1]), int(v.Data[2]))
		} else if sub == 1 {
			for _, k := range v.Data[5 : len(v.Data)-1] {
				hrHis.Heart_rates =
					append(hrHis.Heart_rates, int(k))
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
	case bridge.CMD_SET_TIME:
		timeResponse(v.Data)
	default:
		println("Packet with ID", v.Id, "is not yet implemented")
	}
	return hrHis

}

func Run(ctx context.Context) {
	var hrHis *hrHitory.History
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-bridge.RawDataChannel:
			hrHis = handle(hrHis, v)
		}
	}
}
