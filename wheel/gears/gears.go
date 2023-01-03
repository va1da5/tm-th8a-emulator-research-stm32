package gears

import "strconv"

type Gear struct {
	data  []byte
	Mode  byte  // S/H
	GearH int16 // Gear in H-mode
	GearS byte  // Gear in S-Mode 0x04 - C=center, 0x05 - D=down, 0x06 - U=up
}

const (
	MODE_S = 0x00
	MODE_H = 0x80

	S_CENTER = 0x04 // C
	S_DOWN   = 0x05 // D
	S_UP     = 0x06 // U
)

func (gear *Gear) Parse(data []byte) {
	gear.data = data
	gear.Mode = gear.getMode()

	switch gear.Mode {
	case 'S':
		gear.GearS = gear.getGearS()
	case 'H':
		gear.GearH = gear.getGearH()
	}

}

func (gear *Gear) Current() string {
	out := "Type: " + string(gear.Mode) + " -> "

	switch gear.Mode {
	case 'S':
		switch gear.GearS {
		case 'U':
			out += "Up"
		case 'D':
			out += "Down"
		case 'C':
			out += "Center"
		}

	case 'H':
		out += strconv.FormatInt(int64(gear.GearH), 10)
	}

	return out
}

func (gear *Gear) getMode() byte {
	if gear.data[0] == MODE_H {
		return 'H'
	}
	return 'S'
}

func (gear *Gear) getGearH() int16 {
	// Neutral/Gear 0
	// 1-7 gears
	// Reverse/Gear 8
	gearByte := byte(gear.data[3])

	for i := 0; i < 9; i++ {

		if gearByte == (0x80 >> (8 - i)) {
			return int16(i)
		}
	}
	return 0x0
}

func (gear *Gear) getGearS() byte {

	switch gear.data[4] {
	case S_UP:
		return 'U'
	case S_DOWN:
		return 'D'
	default:
		return 'C'
	}
}

var gear Gear

func main() {
	data := []byte{
		0x00, // Shifter mode 0 - S / 0x80 - H
		0x0C, // Unknown
		0x01, // Unknown
		0x00, // Gear in H-mode
		0x05, // Gear in S-Mode 0x04 - center, 0x05 - down, 0x06 - up
		0x80, // Unknown
		0x80, // Unknown
		0x00, // Y cordinate
		0x00, // X cordinate
		0x00, // Unknown
		0x00, // Unknown
		0x00, // Unknown
		0x00, // Unknown
		0x00, // Unknown
	}

	gear.Parse(data)

	println(gear.Current())

	return
}
