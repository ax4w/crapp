package ui

import (
	"bufio"
	"context"
	"crapp/internal/bridge"
	"crapp/internal/ring"
	"crapp/internal/ui/tools"
	"os"
	"strings"
	"time"
)

func Show(r ring.Ring) {
	r.Listen()
	reader := bufio.NewReader(os.Stdin)
	defer println("done")
	defer r.Disconnect()
	for {
		println("What do you want to do?")
		println("b) battery\nt)time\nrhr)real time heart rate\nhrl) heart rate log\nahrl) active heart rate log\n" +
			"dhrl) disable heart rate log\nshrl) show heart rate log settings" +
			"\nq) quit")
		s, _ := reader.ReadString('\n')
		s = strings.ToLower(strings.TrimSpace(s))
		switch s {
		case "q":
			return
		case "b":
			r.Battery()
			println(<-bridge.StringValueChan)
		case "t":
			r.SetTime()
			println(<-bridge.StringValueChan)
		case "ahrl":
			r.SetHRLogSettings(true, 5)
		case "dhrl":
			r.SetHRLogSettings(false, 5)
		case "shrl":
			r.HRLogSettings()
			println(<-bridge.StringValueChan)
		case "hrl":
			r.TodayHRHistory()
			readings := <-bridge.HeartRateLogChannel
			tools.ShowGraph(readings)
			println("saved graph to file called hr.svg")
		case "rhr":
			println("scanning...")
			var data []int
			r.StartHR()
			//i think this is the only thing that can time out / be stuck forever
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		loop:
			for i := 0; i < 6; i++ {
				select {
				case <-ctx.Done():
					println("timeout")
					break loop
				case v := <-bridge.IntValueChan:
					data = append(data, v)
				}
			}
			r.StopHR()
			for _, v := range data {
				print(v, " ")
			}
			println()
			cancel()
		default:
			continue
		}
	}
}
