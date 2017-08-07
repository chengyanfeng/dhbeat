package util

type Aggregator struct {
 	Ju []P
	Lm map[string]float64
}
// 往聚合器里面增加数据

func (this *Aggregator) Add(p ...P) {

	if IsEmpty(this.Ju) {
		this.Ju = []P{}
	}
	if this.Lm == nil{
		this.Lm = make(map[string]float64)
	}

	// todo

	for i := 0; i < len(p); i++ {
		pl := p[i]
		key := ToString(pl["time_local"]) + "|" + ToString(pl["spid"]) + "|" + ToString(pl["pid"]) + "|" + ToString(pl["dhbeat_hostname"])
		value := pl["bw"].(float64)
		if _, ok := this.Lm[key]; ok {
			this.Lm[key] += value
			for i := 0; i < len(this.Ju); i++ {
				if ToString(this.Ju[i]["time_local"]) == ToString(pl["time_local"]) && ToString(this.Ju[i]["spid"]) == ToString(pl["spid"]) && ToString(this.Ju[i]["pid"]) == ToString(pl["pid"]) {
					t := this.Ju[i]["bw"].(float64) + pl["bw"].(float64)
					this.Ju[i]["bw"] = t
				}
			}
		} else {
			this.Lm[key] = value
			this.Ju = append(this.Ju, pl)
		}
	}
//	fmt.Println(len(this.Lm))

		//if len(this.Ju) == 0 {
		//	this.Ju = append(this.Ju, pl)
		//} else {
		//	for i := 0; i < len(this.Ju); i++ {
		//				jkey := ToString(this.Ju[i]["time_local"]) + "|" + ToString(this.Ju[i]["spid"]) + "|" + ToString(this.Ju[i]["pid"]) + "|" + ToString(this.Ju[i]["dhbeat_hostname"])
		//			if jkey == key {
		//	//	if ToString(this.Ju[i]["time_local"]) == ToString(pl["time_local"]) && ToString(this.Ju[i]["spid"]) == ToString(pl["spid"]) && ToString(this.Ju[i]["pid"]) == ToString(pl["pid"]) {
		//			t := this.Ju[i]["bw"].(float64) + pl["bw"].(float64)
		//			this.Ju[i]["bw"] = t
		//			break
		//		} else if ToString(this.Ju[i]["time_local"]) != ToString(pl["time_local"]) || ToString(this.Ju[i]["spid"]) != ToString(pl["spid"]) || ToString(this.Ju[i]["pid"]) != ToString(pl["pid"]) {
		//			this.Ju = append(this.Ju, pl)
		//			break
		//		}
		//	}
		//}
	}



// 聚合器导出聚合后的数据，同时清空
func (this *Aggregator) Dump()[]P  {
	// todo

	return this.Ju
}
//func (this *Aggregator) Dump()map[string]float64 {
//		// todo
//
//		return this.Lm
//	}

// 聚合器容量
func (this *Aggregator) Size() int {
	// todo
	return len(this.Ju)
}
