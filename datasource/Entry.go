package datasource

type Entry struct {
	Key   string `json:"key"`
	Owner string `json:"owner"`
}

func NewEntry(key Key, data Data) Entry {
	return Entry{
		Key:   string(key),
		Owner: data.owner,
	}
}
