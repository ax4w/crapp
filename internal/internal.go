package internal

import (
	"bufio"
	"crapp/internal/ring"
	"os"
	"strings"

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
