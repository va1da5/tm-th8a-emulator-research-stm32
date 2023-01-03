package gears

type Gear struct {
	data  []byte
	Mode  byte  // S/H
	GearH int16 // Gear in H-mode
	GearS byte  // Gear in S-Mode 0x04 - C=center, 0x05 - D=down, 0x06 - U=up

	gearModes       map[byte]byte
	sequentialGears map[byte]byte
}

const (
	MODE_SEQUENTIAL = 'S'
	MODE_H          = 'H'

	SEQUENTIAL_CENTER = 'C'
	SEQUENTIAL_DOWN   = 'D'
	SEQUENTIAL_UP     = 'U'
)

var GEAR = Gear{
	data: []byte{
		0x80, // Shifter mode 0 - S / 0x80 - H
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
	},

	gearModes: map[byte]byte{
		'S': 0x00, // Sequential
		'H': 0x80, // H type
	},

	sequentialGears: map[byte]byte{
		'C': 0x04, // Center
		'D': 0x05, // Down
		'U': 0x06, // Up
	},
}

func (gear *Gear) SetMode(value byte) {
	gear.data[0] = gear.gearModes[value]
	gear.Mode = value
}

func (gear *Gear) SetGearSequential(value byte) { // C - center, U - up, D - down
	gear.data[4] = gear.sequentialGears[value]
	gear.GearS = value
}

func (gear *Gear) SetGearH(value int16) {
	gear.data[3] = (0x80 >> (8 - value))
	gear.GearH = value
}

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
		out += string(gear.GearS)
	}

	return out
}

func (gear *Gear) getMode() byte {
	if gear.data[0] == gear.gearModes['H'] {
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

	for key, value := range gear.sequentialGears {
		if value == gear.data[4] {
			return key
		}
	}
	return 'C'
}

func (gear *Gear) GetData() []byte {
	return gear.data
}
