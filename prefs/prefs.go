package foxpop

type Entry struct {
	Name  string
	Value interface{}
}

type Data struct {
	Entries []Entry
}
