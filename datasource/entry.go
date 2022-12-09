package datasource

type Entry struct {
	Key    string `json:"key"`
	Owner  string `json:"owner"`
	Writes int    `json:"writes"`
	Reads  int    `json:"reads"`
	Age    int64  `json:"age"`
}

func NewEntry(key Key, data Data) Entry {
	return Entry{
		Key:    string(key),
		Owner:  data.owner,
		Writes: data.writes,
		Reads:  data.reads,
		Age:    data.Age(),
	}
}
