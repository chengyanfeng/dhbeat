package util

type Aggregator struct {
}

// 往聚合器里面增加数据
func (this *Aggregator) Add(p ...P) {
	// todo
}

// 聚合器导出聚合后的数据，同时清空
func (this *Aggregator) Dump() []P {
	// todo
	return []P{}
}

// 聚合器容量
func (this *Aggregator) Size() int {
	// todo
	return 0
}
