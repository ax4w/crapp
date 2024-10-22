package internal

import (
	"bufio"
	"crapp/internal/middle"
	"crapp/internal/ring"
	"fmt"
	"os"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"tinygo.org/x/bluetooth"
)

func saveSelection(s string) {
	if err := os.WriteFile(".ring", []byte(s), os.ModePerm); err != nil {
		println(err.Error())
	}
}

func Scan() ring.Ring {
	println(bluetooth.CharacteristicUUIDRestingHeartRate.String())
	adapter := bluetooth.DefaultAdapter
	filter := GetPaired()
	must("enable BLE stack", adapter.Enable())
	println("scanning...")
	deviceChan := make(chan bluetooth.ScanResult)
	selectionChan := make(chan string)
	go func() {
		adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			///println(device.LocalName())
			if len(device.LocalName()) > 0 && strings.HasPrefix(device.LocalName(), "R02") {
				if len(filter) > 0 && device.Address.String() == filter {
					println("connecting to saved device: ", device.LocalName())
					adapter.StopScan()
					deviceChan <- device
				} else {
					deviceChan <- device
					println("waiting for selection")
					selection := <-selectionChan
					println("got selection", selection)
					if selection == "y" {
						println("stopping scan")
						adapter.StopScan()
						println("stopped scan")
					}
				}
			}
		})
	}()
	var d bluetooth.ScanResult
	if len(filter) == 0 {
		println("No paired device found")
		reader := bufio.NewReader(os.Stdin)
		for {
			d = <-deviceChan
			println("Do you want to connec to", d.LocalName(), "(", d.Address.String(), ")? y/n")
			s, err := reader.ReadString('\n')
			if err != nil {
				panic(err.Error())
			}
			s = strings.TrimSpace(strings.ToLower(s))
			selectionChan <- s
			if s == "y" {
				saveSelection(d.Address.String())
				break
			}
		}
	} else {
		d = <-deviceChan
	}
	println("connecting to", d.Address.String())
	r, _ := adapter.Connect(d.Address, bluetooth.ConnectionParams{})
	return ring.New(r, adapter)
}

func makeLabels(data []middle.HRObj) []string {
	labels := make([]string, len(data))
	for i, v := range data {
		labels[i] = fmt.Sprintf("%d bpm", v.HR)
	}
	return labels
}

func showGraph(data []middle.HRObj) {
	if len(data) == 0 {
		return
	}

	p := plot.New()
	p.Title.Text = "HR Vis"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "HR"
	p.X.Tick.Marker = plot.TimeTicks{
		Format: "15:04:05", // Change this format as needed
		Time:   func(t float64) time.Time { return time.Unix(int64(t), 0) },
	}

	var pts plotter.XYs
	//s := data[0].T
	for i, v := range data {
		pts = append(pts, plotter.XY{
			X: float64(v.T.Unix()),
			Y: float64(v.HR),
		})
		println("point", i, "x", pts[i].X, "y", pts[i].Y)
	}
	err := plotutil.AddLinePoints(p, "Data", pts)
	labels, _ := plotter.NewLabels(plotter.XYLabels{
		XYs:    pts,
		Labels: makeLabels(data),
	})
	labels.Offset.X = vg.Points(0)
	labels.Offset.Y = vg.Points(5)
	p.Add(labels)
	if err != nil {
		println(err.Error())
		return
	}
	if err := p.Save(30*vg.Inch, 10*vg.Inch, "hr.svg"); err != nil {
		panic(err)
	}
}

func Run(r ring.Ring) {
	r.Listen()
	reader := bufio.NewReader(os.Stdin)
	defer println("done")
	defer r.Disconnect()
	for {
		println("What do you want to do?")
		println("b) battery t)time rhr)real time heart rate hrl) heart rate log ahrl) active heart rate log\n" +
			"dhrl) disable heart rate log shrl) show heart rate log settings" +
			"\n q) quit")
		s, _ := reader.ReadString('\n')
		s = strings.ToLower(strings.TrimSpace(s))
		switch s {
		case "q":
			return
		case "b":
			r.Battery()
		case "t":
			r.SetTime()
		case "ahrl":
			r.SetHRLogSettings(true, 5)
		case "dhrl":
			r.SetHRLogSettings(false, 5)
		case "shrl":
			r.HRLogSettings()
		case "hrl":
			r.TodayHRHistory()
			readings := <-middle.HeartRateLogChannel
			showGraph(readings)
			println("saved graph to file")
			continue
		case "rhr":
			println("scanning...")
			var data []int
			r.StartHR()
			for i := 0; i < 6; i++ {
				data = append(data, <-middle.HeartRateChannel)
			}
			r.StopHR()
			for _, v := range data {
				print(v, " ")
			}
			println()
			continue //cmd is over, we dont need to wait for feedback
		default:
			continue
		}

		<-middle.WaitChan
	}
	/*t := time.Tick(time.Minute)
	for {
		r.StartHR()
		for i := 0; i < 6; i++ {
			data = append(data, <-middle.HeartRateChannel)
		}
		r.StopHR()
		for _, v := range data {
			print(v, " ")
		}
		println()
		<-t
	}*/
	//r.HRLogSettings()
	//<-middle.HeartRateChannel
	//r.SetHRLogSettings(false, 15)
	//<-middle.HeartRateChannel
	//r.TodayHRHistory()
	//select {}
	//r.Disconnect()
	/*for {
	r.StartHR()
	time.Sleep(2 * time.Second)
	}*/
}

func GetPaired() string {
	if _, err := os.Stat(".ring"); err != nil {
		return ""
	}
	bytes, err := os.ReadFile(".ring")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(bytes))
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
