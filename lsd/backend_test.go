package main

type TestBackend struct {
	data map[string]Rows
}

func (tb *TestBackend) Init() {}
func (tb *TestBackend) GetData(table string) Rows {
	return tb.data[table]
}
