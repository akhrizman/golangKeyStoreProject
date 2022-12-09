package datasource

import (
	"fmt"
	"time"
)

// Key type
type Key string

func (key Key) String() string {
	return fmt.Sprintf("Key(%s)", string(key))
}

// Data Contains the value of the key-value store plus any additional data.
type Data struct {
	owner    string
	value    string
	lastUsed time.Time
	writes   int
	reads    int
}

func NewData(user string, value string) Data {
	return Data{
		owner:    user,
		value:    value,
		lastUsed: time.Now(),
		writes:   0,
		reads:    0,
	}
}
func (d *Data) SetToCurrentTime() {
	d.lastUsed = time.Now()
}
func (d *Data) GetValue() string {
	return d.value
}
func (d *Data) Age() int64 {
	currentTime := time.Now()
	return currentTime.UnixMilli() - d.lastUsed.UnixMilli()
}
