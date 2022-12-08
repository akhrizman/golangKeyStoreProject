package datasource

import "fmt"

// Key type
type Key string

func (key Key) String() string {
	return fmt.Sprintf("Key(%s)", string(key))
}

// Data to contain owner and value
type Data struct {
	owner string `json:"owner"`
	value string `json:"value"`
}

func NewData(user string, value string) Data {
	return Data{owner: user, value: value}
}
func (d *Data) GetOwner() string {
	return d.owner
}
func (d *Data) SetOwner(owner string) {
	d.owner = owner
}
func (d *Data) GetValue() string {
	return d.value
}
func (d *Data) SetValue(value string) {
	d.value = value
}
