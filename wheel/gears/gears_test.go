package gears

import (
	"testing"
)

// Test Mode H type
func TestModeH(t *testing.T) {
	data := []byte{
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
	}

	gear.Parse(data)

	if gear.Mode != 'H' {
		t.Errorf("Expected to be H mode")
	}
}

// Test Mode Sequential type
func TestModeS(t *testing.T) {
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

	if gear.Mode != 'S' {
		t.Errorf("Expected to be S mode")
	}
}

func TestGearsH(t *testing.T) {
	data := []byte{
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
	}

	for i := 0; i < 9; i++ {
		// TODO: implement parser for the below
		data[3] = (0x80 >> (8 - i))
		gear.Parse(data)

		if gear.Mode != 'H' {
			t.Errorf("Expected to be H mode")
		}

		if gear.GearH != int16(i) {
			t.Errorf("Expected to be %d gear, got %d", i, gear.GearH)
		}
	}

	gear.Parse(data)

}

func TestGearsS(t *testing.T) {
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

	gearsBytes := []byte{
		0x04, 0x05, 0x06,
	}

	gearNames := []byte{
		'C', 'D', 'U',
	}

	for i := 0; i < len(gearsBytes); i++ {
		data[4] = gearsBytes[i]

		gear.Parse(data)

		if gear.Mode != 'S' {
			t.Errorf("Expected to be H mode")
		}

		if gear.GearS != gearNames[i] {
			t.Errorf("Expected to be %s gear, got %s", gearNames[i], gear.GearS)
		}
	}

	gear.Parse(data)

}
