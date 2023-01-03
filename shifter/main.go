package main

import (
	"machine"
	"shifter/gears"
	"strconv"
	"time"
)

var gear = gears.GEAR

const address = uint16(0x01)

var empty = []byte{}

func formatAddress(value uint16) string {
	return "0x" + strconv.FormatInt(int64(value), 16)
}

func send(i2c *machine.I2C, data []byte) {
	err := i2c.Tx(address, data, empty)
	if err != nil {
		println(" [FAILED]")
	} else {
		println(" [OK]")
	}
}

func loopGears(i2c *machine.I2C) {

	for {
		// Sequential mode
		gear.SetMode(gears.MODE_SEQUENTIAL)

		gear.SetGearSequential(gears.SEQUENTIAL_CENTER)
		print("S -> Center")
		send(i2c, gear.GetData())
		println()
		time.Sleep(time.Second * 1)

		gear.SetGearSequential(gears.SEQUENTIAL_DOWN)
		print("S -> Down")
		send(i2c, gear.GetData())
		println()
		time.Sleep(time.Second * 1)

		gear.SetGearSequential(gears.SEQUENTIAL_UP)
		print("S -> Up")
		send(i2c, gear.GetData())
		println()
		time.Sleep(time.Second * 1)

		// -------------------------------------
		// H mode
		gear.SetMode(gears.MODE_H)

		for i := 0; i < 9; i++ {
			gear.SetGearH(int16(i))
			print("H -> ", strconv.FormatInt(int64(i), 10))
			send(i2c, gear.GetData())
			println()
			time.Sleep(time.Second * 1)
		}

		println("----------------------")
		machine.I2C0.Configure(machine.I2CConfig{})
		// time.Sleep(time.Second * 1)

	}
}

func main() {
	println("REBOOT")
	uart := machine.DefaultUART
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{})

	uartReceived := []byte{0x0}

	for {
		time.Sleep(time.Millisecond * 10)
		if uart.Buffered() == 0 {
			continue
		}

		uart.Read(uartReceived)
		command := uartReceived[0]

		switch command {
		case ' ':
			send(i2c, gear.GetData())
		case 's':
			gear.SetMode(gears.MODE_SEQUENTIAL)
			println("Mode: Sequential")

		case 'u':
			gear.SetGearSequential(gears.SEQUENTIAL_UP)
			println("Sequential: Up")

		case 'd':
			gear.SetGearSequential(gears.SEQUENTIAL_DOWN)
			println("Sequential: Down")

		case 'c':
			gear.SetGearSequential(gears.SEQUENTIAL_CENTER)
			println("Sequential: Center")

		case 'h':
			gear.SetMode(gears.MODE_H)
			println("Mode: H")

		case '0':
			gear.SetGearH(int16(0))
			println("H: 0/N")
		case '1':
			gear.SetGearH(int16(1))
			println("H: 1")
		case '2':
			gear.SetGearH(int16(2))
			println("H: 2")
		case '3':
			gear.SetGearH(int16(3))
			println("H: 3")
		case '4':
			gear.SetGearH(int16(4))
			println("H: 4")
		case '5':
			gear.SetGearH(int16(5))
			println("H: 5")
		case '6':
			gear.SetGearH(int16(6))
			println("H: 6")
		case '7':
			gear.SetGearH(int16(7))
			println("H: 7")
		case '8':
			gear.SetGearH(int16(8))
			println("H: 8/R")
		}

	}

}
