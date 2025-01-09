package main

type UpdateV struct {
	Cell
	I FlexString `json:"i"`
	T string     `json:"t"`
}

type UpdateRV struct {
	I     FlexString `json:"i"`
	Range struct {
		Column []int `json:"column"`
		Row    []int `json:"row"`
	} `json:"range"`
	T string        `json:"t"`
	V [][]CellValue `json:"v"`
}

type UpdateCG struct {
	T string      `json:"t"`
	I FlexString  `json:"i"`
	K string      `json:"k"`
	V interface{} `json:"v"`
}

type UpdateCommon struct {
	T string      `json:"t"`
	I FlexString  `json:"i"`
	K string      `json:"k"`
	V interface{} `json:"v"`
}

type UpdateCalcChain struct {
	I   FlexString `json:"i"`
	Op  string     `json:"op"`
	Pos int        `json:"pos"`
	T   string     `json:"t"`
	V   string     `json:"v"`
}

type UpdateRowColumn struct {
	I FlexString `json:"i"`
	T string     `json:"t"`
	V struct {
		Index int `json:"index"`
		Len   int `json:"len"`
	} `json:"v"`
	RC string `json:"rc"`
}

type UpdateFilter struct {
	I FlexString   `json:"i"`
	T string       `json:"t"`
	V *FilterValue `json:"v"`
}

type FilterValue struct {
	Filter       []Filter     `json:"filter"`
	FilterSelect FilterSelect `json:"filter_select"`
}

type AddSheet struct {
	I FlexString `json:"i"`
	T string     `json:"t"`
	V *Sheet     `json:"v"`
}

type CopySheet struct {
	I FlexString `json:"i"`
	T string     `json:"t"`
	V struct {
		CopyIndex FlexString `json:"copyindex"`
		Name      string     `json:"name"`
	} `json:"v"`
}

type DeleteSheet struct {
	I string `json:"i"`
	T string `json:"t"`
	V struct {
		DeleteIndex FlexString `json:"deleIndex"`
	} `json:"v"`
}

type RecoverSheet struct {
	I FlexString `json:"i"`
	T string     `json:"t"`
	V struct {
		RecoverIndex FlexString `json:"reIndex"`
	} `json:"v"`
}

type UpdateSheetOrder struct {
	I FlexString     `json:"i"`
	T string         `json:"t"`
	V map[string]int `json:"v"`
}

type ToggleSheet struct {
	I FlexString `json:"i"`
	T string     `json:"t"`
	V string     `json:"v"`
}

type HideOrShowSheet struct {
	I   FlexString `json:"i"`
	T   string     `json:"t"`
	V   int        `json:"v"`
	Op  string     `json:"op"`
	Cur string     `json:"cur"`
}
