package ring

import (
	"crapp/internal/middle"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	ringMU sync.Mutex
)

type Ring struct {
	Device  bluetooth.Device
	adapter *bluetooth.Adapter
	rxChar  bluetooth.DeviceCharacteristic
}

func New(d bluetooth.Device, a *bluetooth.Adapter) Ring {
	r := Ring{
		Device:  d,
		adapter: a,
	}
	r.rxChar = r.getCharacteristic(UART_SERVICE_UUID, UART_RX_CHAR_UUID)
	return r
}

func (r Ring) Listen() {
	chrc := r.getCharacteristic(UART_SERVICE_UUID, UART_TX_CHAR_UUID)
	println("now listening")
	chrc.EnableNotifications(func(buf []byte) {
		id, _ := strconv.Atoi(fmt.Sprintf("%d", buf[0]))
		middle.BackToFront <- middle.Packet{
			Id:   id,
			Data: buf[1:],
		}
	})
}

func (r Ring) SetTime() {
	n := time.Now().UTC()
	data := make([]byte, 7)
	data[0] = byte(byteToBcd(n.Year() % 1000))
	data[1] = byte(byteToBcd(int(n.Month())))
	data[2] = byte(byteToBcd(n.Day()))
	data[3] = byte(byteToBcd(n.Hour()))
	data[4] = byte(byteToBcd(n.Minute()))
	data[5] = byte(byteToBcd(n.Second()))
	data[6] = byte(1)

	time_packet, _ := makePacket(middle.CMD_SET_TIME, data)
	d, err := r.rxChar.WriteWithoutResponse(time_packet)
	if err != nil {
		panic(err.Error())
	}
	println("send", d)
	r.Blink()

}

func (r Ring) getCharacteristic(service, characteristics string) bluetooth.DeviceCharacteristic {
	suid, err := bluetooth.ParseUUID(service)
	if err != nil {
		panic(err.Error())
	}
	srvcs, err := r.Device.DiscoverServices([]bluetooth.UUID{suid})
	if err != nil {
		panic(err.Error())
	}
	println("found services")
	srvc := srvcs[0]
	uid, err := bluetooth.ParseUUID(characteristics)
	if err != nil {
		panic(err.Error())
	}
	chrs, err := srvc.DiscoverCharacteristics([]bluetooth.UUID{uid})
	if err != nil {
		panic(err.Error())
	}
	println("found chars")
	return chrs[0]

}

func (r Ring) Blink() {
	println("sending blink")
	//rxChar := r.getCharacteristic(UART_SERVICE_UUID, UART_RX_CHAR_UUID)
	blink_twice_paket, _ := makePacket(middle.CMD_BLINK_TWICE, nil)
	n, err := r.rxChar.WriteWithoutResponse(blink_twice_paket)
	if err != nil {
		panic(err.Error())
	}
	println("send", n)
}

func (r Ring) Disconnect() {
	r.Device.Disconnect()
}

func (r Ring) Battery() {
	packet, _ := makePacket(middle.CMD_BATTERY, nil)
	r.rxChar.WriteWithoutResponse(packet)
	println("send battery packet")
}

func (r Ring) StartHR() {
	packet, _ := makePacket(middle.CMD_START_HEART_RATE, []byte{0x01, 0x00})
	r.rxChar.WriteWithoutResponse(packet)
}

func (r Ring) ContinueHR() {
	packet, _ := makePacket(middle.CMD_REAL_TIME_HEART_RATE, []byte("3"))
	r.rxChar.WriteWithoutResponse(packet)
}

func (r Ring) StopHR() {
	packet, _ := makePacket(middle.CMD_STOP_HEART_RATE, []byte("\x01\x00\x00"))
	r.rxChar.WriteWithoutResponse(packet)
}

func (r Ring) StartSPO2() {
	packet, _ := makePacket(middle.CMD_START_HEART_RATE, []byte{0x03, 0x25})
	r.rxChar.WriteWithoutResponse(packet)
}

func (r Ring) StopSPO2() {
	packet, _ := makePacket(middle.CMD_START_HEART_RATE, []byte{0x03, 0x00, 0x00})
	r.rxChar.WriteWithoutResponse(packet)
}

func (r Ring) HRLogSettings() {
	packet, _ := makePacket(middle.CMD_HEART_RATE_LOG_SETTINGS, []byte("\x01"))
	r.rxChar.WriteWithoutResponse(packet)
}

func (r Ring) SetHRLogSettings(enabled bool, interval int) {
	e := 2
	if enabled {
		e = 1
	}
	packet, _ := makePacket(middle.CMD_HEART_RATE_LOG_SETTINGS, []byte{2, byte(e), byte(interval)})
	r.rxChar.WriteWithoutResponse(packet)
}

func (r Ring) TodayHRHistory() {
	now := time.Now().UTC()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(midnight))
	packet, _ := makePacket(middle.CMD_READ_HEART_RATE, data)
	r.rxChar.WriteWithoutResponse(packet)
}
