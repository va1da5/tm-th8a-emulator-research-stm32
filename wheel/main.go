// Thrustmasters Wheel emulator package for testing emulator implementation.
package main

import (
	"machine"
	"strconv"
	"time"

	"wheel/gears"
	wheel "wheel/i2c"
)

var gear gears.Gear

const (
	slaveAddress uint32 = 0x01 // 7 bit address
)

var i2c = wheel.I2C

func formatAddress(value int64) string {
	return "0x" + strconv.FormatInt(value, 16)
}

func main() {

	println("<Reboot>")

	machine.I2C0.Configure(machine.I2CConfig{})
	println("I2C0 UP")

	i2c.Configure(machine.I2C0, wheel.I2CSlaveConfig{Address: slaveAddress})
	println("I2C0 SLAVE UP")

	i2c.SetInterrupt()
	println("I2C0 INT UP")

	data := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	for {

		if i2c.Buffered() > 13 {
			println()
			_, err := i2c.Read(data)
			if err != nil {
				println("Error reading data")
			}

			gear.Parse(data)
			println(gear.Current())

			// print("Data received:")
			// for i := 0; i < len(data); i++ {
			// 	print(" " + formatAddress(int64(data[i])))
			// }
			println()

		}

		print(".")

		// println(i2c.DebugSR())

		time.Sleep(time.Millisecond * 100)

	}
}
