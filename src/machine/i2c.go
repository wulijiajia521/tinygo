// +build avr nrf

package machine

// TWI_FREQ is the I2C bus speed. Normally either 100 kHz, or 400 kHz for high-speed bus.
const (
	TWI_FREQ_100KHZ = 100000
	TWI_FREQ_400KHZ = 400000
)

// Write writes a slice of data bytes to a peripheral.
// A transmission must already have been started with Start() and must be
// stopped afterwards with Stop().
func (i2c I2C) Write(data []byte) {
	for _, v := range data {
		i2c.WriteByte(v)
	}
}

// Read reads a slice of data bytes from an I2C peripheral
// A transmission must already have been started with Start() and must be
// stopped afterwards with Stop().
func (i2c I2C) Read(data []byte) {
	for i, _ := range data {
		data[i] = i2c.ReadByte()
	}
}

// WriteRegister transmits first the register and then the data to the
// peripheral device.
//
// Many I2C-compatible devices are organized in terms of registers. This method
// is a shortcut to easily write to such registers.
func (i2c I2C) WriteRegister(address uint8, register uint8, data []byte) {
	i2c.Start(address, true) // start for writing
	i2c.WriteByte(register)  // write the specified register
	i2c.Write(data)          // write data to the register
	i2c.Stop()
}

// ReadRegister transmits the register, restarts the connection as a read
// operation, and reads the response.
//
// Many I2C-compatible devices are organized in terms of registers. This method
// is a shortcut to easily read such registers.
func (i2c I2C) ReadRegister(address uint8, register uint8, data []byte) {
	i2c.Start(address, true)  // start for writing
	i2c.WriteByte(register)   // write the specified register
	i2c.Start(address, false) // re-start transmission as a read operation
	i2c.Read(data)            // read data from the register
	i2c.Stop()
}
