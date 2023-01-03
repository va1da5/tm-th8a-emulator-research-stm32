// Package i2c provides a wrapper around original I2C peripheral to
// configure it as a slave for receiving data
package i2c

import (
	"device/stm32"
	"errors"
	"machine"
	"runtime/interrupt"
)

const (
	flagOVR     = 0x00010800
	flagAF      = 0x00010400
	flagARLO    = 0x00010200
	flagBERR    = 0x00010100
	flagTXE     = 0x00010080
	flagRXNE    = 0x00010040
	flagSTOPF   = 0x00010010
	flagADD10   = 0x00010008
	flagBTF     = 0x00010004
	flagADDR    = 0x00010002
	flagSB      = 0x00010001
	flagDUALF   = 0x00100080
	flagGENCALL = 0x00100010
	flagTRA     = 0x00100004
	flagBUSY    = 0x00100002
	flagMSL     = 0x00100001
)

var (
	errI2CWriteTimeout       = errors.New("I2C timeout during write")
	errI2CReadTimeout        = errors.New("I2C timeout during read")
	errI2CBusReadyTimeout    = errors.New("I2C timeout on bus ready")
	errI2CSignalStartTimeout = errors.New("I2C timeout on signal start")
	errI2CSignalReadTimeout  = errors.New("I2C timeout on signal read")
	errI2CSignalStopTimeout  = errors.New("I2C timeout on signal stop")
	errI2CAckExpected        = errors.New("I2C error: expected ACK not NACK")
	errI2CBusError           = errors.New("I2C bus error")
	errI2CBufferEmpty        = errors.New("I2C buffer empty")
)

type I2CWrapper struct {
	Buffer       *machine.RingBuffer
	I2C          *machine.I2C
	Bus          *stm32.I2C_Type
	Interrupt_EV interrupt.Interrupt
	Interrupt_ER interrupt.Interrupt
	Response     []byte
}

var (
	I2C  = &_I2C // I2C contains address of _I2C instance
	_I2C = I2CWrapper{
		Buffer:   machine.NewRingBuffer(),
		Response: []byte{0xde, 0xed, 0xbe, 0xef, 0x00},
	}
)

type I2CSlaveConfig struct {
	Address uint32
}

var SR1 = []string{
	"SB", "ADDR", "BTF", "ADD10", "STOPF", "Res.", "RxNE", "TxE", "BERR", "ARLO", "AF", "OVR", "PEC ERR", // "Res.", "TIME OUT", "SMB ALERT",
}

var SR2 = []string{
	"MSL", "BUSY", "TRA", "Res.", "GEN CALL", "SMBDE FAULT", "SMB HOST", "DUALF",
}

func (self *I2CWrapper) Configure(i2c *machine.I2C, config I2CSlaveConfig) {
	self.I2C = i2c
	self.Bus = i2c.Bus

	slaveAddress := uint32(0xAA)

	if !(config.Address == 0) {
		slaveAddress = config.Address
	}

	// configured slave address
	self.Bus.OAR1.SetBits(slaveAddress << 1)
	self.Bus.OAR2.SetBits(slaveAddress << 1)

	// enable ACK
	self.Bus.CR1.SetBits(stm32.I2C_CR1_ACK)

	// enable POS
	self.Bus.CR1.SetBits(stm32.I2C_CR1_POS)
}

func boolToString(val bool) string {
	if val {
		return "1"
	}
	return "0"
}

func (i2c *I2CWrapper) DebugSR() string {

	output := "[SR1] "
	for i := 0; i < len(SR1); i++ {
		output += "| " + SR1[i] + " " + boolToString(i2c.Bus.SR1.HasBits(1<<i)) + " "
	}
	output += "\r\n[SR2] "
	for i := 0; i < len(SR2); i++ {
		output += "| " + SR2[i] + " " + boolToString(i2c.Bus.SR2.HasBits(1<<i)) + " "
	}

	return output
}

func (i2c *I2CWrapper) hasFlag(flag uint32) bool {
	const mask = 0x0000FFFF
	if uint8(flag>>16) == 1 {
		return i2c.Bus.SR1.HasBits(flag & mask)
	} else {
		return i2c.Bus.SR2.HasBits(flag & mask)
	}
}

func (i2c *I2CWrapper) clearFlagADDR() {
	i2c.Bus.SR1.Get()
	i2c.Bus.SR2.Get()
}

func (i2c *I2CWrapper) handleInterrupt(interrupt.Interrupt) {
	if i2c.hasFlag(flagBUSY) {
		i2c.Bus.CR1.SetBits(stm32.I2C_CR1_ACK)
	}

	// Executed when DATA WRITE command is sent to the slave
	if i2c.hasFlag(flagRXNE) {
		i2c.Receive(byte((i2c.Bus.DR.Get() & 0xFF)))
		i2c.Bus.CR1.SetBits(stm32.I2C_CR1_ACK)
	}

	// Executed when DATA READ command is sent to the slave
	if i2c.hasFlag(flagTXE) {
		i2c.controllerTransmit(i2c.Response)
	}

	if i2c.hasFlag(flagSTOPF) {
		i2c.Bus.SR1.Get()
		i2c.Bus.CR1.ClearBits(stm32.I2C_CR1_STOP_Stop << stm32.I2C_CR1_STOP_Pos)
	}

	// println(i2c.DebugSR())
}

func (i2c *I2CWrapper) SetInterrupt() {
	// https://www.st.com/resource/en/reference_manual/cd00171190-stm32f101xx-stm32f102xx-stm32f103xx-stm32f105xx-and-stm32f107xx-advanced-arm-based-32-bit-mcus-stmicroelectronics.pdf
	// Page: 770
	i2c.Bus.CR2.SetBits(stm32.I2C_CR2_ITEVTEN | stm32.I2C_CR2_ITBUFEN)
	i2c.Interrupt_EV = interrupt.New(stm32.IRQ_I2C1_EV, _I2C.handleInterrupt)
	// Examples: 0xff (lowest priority), 0xc0 (low priority), 0x00 (highest possible
	// priority).
	i2c.Interrupt_EV.SetPriority(0x00)
	i2c.Interrupt_EV.Enable()
}

func (i2c *I2CWrapper) SetErrorInterrupt() {
	// https://www.st.com/resource/en/reference_manual/cd00171190-stm32f101xx-stm32f102xx-stm32f103xx-stm32f105xx-and-stm32f107xx-advanced-arm-based-32-bit-mcus-stmicroelectronics.pdf
	// Page: 771
	i2c.Bus.CR2.SetBits(stm32.I2C_CR2_ITERREN)
	// TODO: implement the below handler
	i2c.Interrupt_EV = interrupt.New(stm32.IRQ_I2C1_ER, func(interrupt.Interrupt) { println("IRQ_I2C1_ER") })
	// Examples: 0xff (lowest priority), 0xc0 (low priority), 0x00 (highest possible
	// priority).
	i2c.Interrupt_EV.SetPriority(0x00)
	i2c.Interrupt_EV.Enable()
}

func (i2c *I2CWrapper) Receive(data byte) {
	i2c.Buffer.Put(data)
}

// Buffered returns the number of bytes currently stored in the RX buffer.
func (i2c *I2CWrapper) Buffered() int {
	return int(i2c.Buffer.Used())
}

// Read from the RX buffer.
func (i2c *I2CWrapper) Read(data []byte) (n int, err error) {
	// check if RX buffer is empty
	size := i2c.Buffered()
	if size == 0 {
		return 0, nil
	}

	// Make sure we do not read more from buffer than the data slice can hold.
	if len(data) < size {
		size = len(data)
	}

	// only read number of bytes used from buffer
	for i := 0; i < size; i++ {
		v, _ := i2c.ReadByte()
		data[i] = v
	}

	return size, nil
}

// ReadByte reads a single byte from the RX buffer.
// If there is no data in the buffer, returns an error.
func (i2c *I2CWrapper) ReadByte() (byte, error) {
	// check if RX buffer is empty
	buf, ok := i2c.Buffer.Get()
	if !ok {
		return 0, errI2CBufferEmpty
	}
	return buf, nil
}

func (i2c *I2CWrapper) controllerTransmit(w []byte) error {

	// disable POS
	i2c.Bus.CR1.ClearBits(stm32.I2C_CR1_POS)

	pos := 0
	rem := len(w)

	// clear ADDR flag
	i2c.clearFlagADDR()

	for rem > 0 {
		// wait for TXE flag set
		if !i2c.waitForFlagOrError(flagTXE, true) {
			return errI2CAckExpected
		}

		// write data to DR
		i2c.Bus.DR.Set(uint32(w[pos]))
		// update counters
		pos++
		rem--

		if i2c.hasFlag(flagBTF) && rem != 0 {
			// write data to DR
			i2c.Bus.DR.Set(uint32(w[pos]))
			// update counters
			pos++
			rem--
		}

		// wait for transfer finished flag BTF set
		if !i2c.waitForFlagOrError(flagBTF, true) {
			return errI2CWriteTimeout
		}
	}

	// generate stop condition
	i2c.Bus.CR1.SetBits(stm32.I2C_CR1_STOP)

	return nil
}

func (i2c *I2CWrapper) waitForFlagOrError(flag uint32, set bool) bool {
	const tryMax = 10000
	hasFlag := false
	for i := 0; !hasFlag && i < tryMax; i++ {
		if hasFlag = i2c.hasFlag(flag) == set; !hasFlag {
			// check for ACK failure
			if i2c.hasFlag(flagAF) {
				// generate stop condition
				i2c.Bus.CR1.SetBits(stm32.I2C_CR1_STOP)
				// clear pending flags
				i2c.clearFlag(flagAF)
				return false
			} else if i2c.hasFlag(flagSTOPF) {
				// clear stop flag
				i2c.clearFlag(flagSTOPF)
				return false
			}
		}
	}
	return hasFlag
}

func (i2c *I2CWrapper) clearFlag(flag uint32) {
	const mask = 0x0000FFFF
	i2c.Bus.SR1.Set(^(flag & mask))
}

func (i2c *I2CWrapper) waitForFlag(flag uint32, set bool) bool {
	const tryMax = 10000
	hasFlag := false
	for i := 0; !hasFlag && i < tryMax; i++ {
		hasFlag = i2c.hasFlag(flag) == set
	}
	return hasFlag
}
