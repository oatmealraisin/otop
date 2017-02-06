package otop

type Entry map[string]string

// TODO: Find a way to make this *Entry
func (e Entry) String() string {
	result := ""
	for _, v := range e {
		result += " " + v
	}
	return result
}

// Each entry should have a unique id
type Entries []Entry

// TODO: Stub
func (e *Entries) String() string {
	return ""
}

func (e *Entries) Sort(key string) {
}
