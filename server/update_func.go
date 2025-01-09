package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var updMap = map[string]func(msg []byte, content *string){
	"v":    updateGrid,
	"rv":   updateGridMulti,
	"cg":   updateGridConfig,
	"all":  updateGridCommon,
	"fc":   updateCalcChain,
	"drc":  updateRowColumn,
	"arc":  updateRowColumn,
	"fsc":  updateFilter,
	"fsr":  updateFilter,
	"sha":  addSheet,
	"shc":  copySheet,
	"shd":  deleteSheet,
	"shre": recoverSheet,
	"shr":  updateSheetOrder,
	"sh":   hideOrShowSheet,
	"shs":  toggleSheet,
}

func updateGrid(reqmsg []byte, content *string) {
	req := new(UpdateV)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		sheetArrIndx = ""
		cellArrIdx   = ""
	)
	gjson.Parse(*content).ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.I) {
			sheetArrIndx = key.String()
			value.Get("celldata").ForEach(func(k2, v2 gjson.Result) bool {
				if FlexInt(v2.Get("r").Int()) == req.R && FlexInt(v2.Get("c").Int()) == req.C {
					cellArrIdx = k2.String()
					return false
				}
				return true
			})
			return false
		}
		return true
	})

	if cellArrIdx == "" {
		cellArrIdx = "-1"
	}
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}
	*content, err = sjson.Set(*content, fmt.Sprintf("%s.celldata.%s", sheetArrIndx, cellArrIdx), &req.Cell)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func updateGridMulti(reqmsg []byte, content *string) {
	req := new(UpdateRV)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(req.Range.Column) < 2 || len(req.Range.Row) < 2 {
		fmt.Println("参数错误")
		return
	}

	cellPosMap := make(map[CellPos]bool)

	var (
		doc          = gjson.Parse(*content)
		sheetArrIndx = ""
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.I) {
			sheetArrIndx = key.String()
			return false
		}
		return true
	})
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}

	// 修改已存在单元格
	doc.Get(sheetArrIndx + ".celldata").ForEach(func(key, value gjson.Result) bool {
		var r, c = int(value.Get("r").Int()), int(value.Get("c").Int())
		if (r >= req.Range.Row[0] && r <= req.Range.Row[1]) && (c >= req.Range.Column[0] && c <= req.Range.Column[1]) {
			cell := Cell{
				C: FlexInt(c),
				R: FlexInt(r),
				V: req.V[r-req.Range.Row[0]][c-req.Range.Column[0]],
			}
			if !IsEmptyCell(cell) {
				*content, err = sjson.Set(*content, fmt.Sprintf("%s.celldata.%s", sheetArrIndx, key.String()), &cell)
				if err != nil {
					fmt.Println(err)
				}
			}
			cellPosMap[CellPos{R: r, C: c}] = true
		}
		return true
	})
	// 插入新单元格
	for r := req.Range.Row[0]; r <= req.Range.Row[1]; r++ {
		for c := req.Range.Column[0]; c <= req.Range.Column[1]; c++ {
			if cellPosMap[CellPos{R: r, C: c}] {
				continue
			}
			cell := Cell{
				C: FlexInt(c),
				R: FlexInt(r),
				V: req.V[r-req.Range.Row[0]][c-req.Range.Column[0]],
			}
			if !IsEmptyCell(cell) {
				*content, err = sjson.Set(*content, fmt.Sprintf("%s.celldata.-1", sheetArrIndx), &cell)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func updateGridConfig(reqmsg []byte, content *string) {
	req := new(UpdateCG)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc          = gjson.Parse(*content)
		sheetArrIndx = ""
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.I) {
			sheetArrIndx = key.String()
			return false
		}
		return true
	})
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}

	*content, err = sjson.Set(*content, fmt.Sprintf("%s.config.%s", sheetArrIndx, req.K), &req.V)
	if err != nil {
		fmt.Println(err)
	}
}

func updateGridCommon(reqmsg []byte, content *string) {
	req := new(UpdateCommon)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc          = gjson.Parse(*content)
		sheetArrIndx = ""
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.I) {
			sheetArrIndx = key.String()
			return false
		}
		return true
	})
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}

	*content, err = sjson.Set(*content, fmt.Sprintf("%s.%s", sheetArrIndx, req.K), &req.V)
	if err != nil {
		fmt.Println(err)
	}
}

func updateCalcChain(reqmsg []byte, content *string) {
	req := new(UpdateCalcChain)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	calcChain := new(CalcChain)
	err = json.Unmarshal([]byte(req.V), calcChain)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc          = gjson.Parse(*content)
		sheetArrIndx = ""
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.I) {
			sheetArrIndx = key.String()
			return false
		}
		return true
	})
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}

	switch req.Op {
	case "add":
		*content, err = sjson.Set(*content, fmt.Sprintf("%s.calcChain.-1", sheetArrIndx), calcChain)
		if err != nil {
			fmt.Println(err)
		}
	case "update":
		*content, err = sjson.Set(*content, fmt.Sprintf("%s.calcChain.%d", sheetArrIndx, req.Pos), calcChain)
		if err != nil {
			fmt.Println(err)
		}
	case "del":
		*content, err = sjson.Delete(*content, fmt.Sprintf("%s.calcChain.%s", sheetArrIndx, req.Pos))
		if err != nil {
			fmt.Println(err) // todo 删除空
		}
	}

	// todo 二次设置？？
}

func updateRowColumn(reqmsg []byte, content *string) {
	req := new(UpdateRowColumn)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc          = gjson.Parse(*content)
		sheetArrIndx = ""
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.I) {
			sheetArrIndx = key.String()
			return false
		}
		return true
	})
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}

	switch req.T {
	case "drc": //删除行或列
		ds := req.V.Index
		de := req.V.Index + req.V.Len
		doc.Get(sheetArrIndx + ".celldata").ForEach(func(key, value gjson.Result) bool {
			var rc = int(value.Get(req.RC).Int())
			if rc >= ds && rc < de {
				*content, err = sjson.Delete(*content, fmt.Sprintf("%s.celldata.%s", sheetArrIndx, key.String()))
				if err != nil {
					fmt.Println(err)
				}
			} else if rc >= de {
				*content, err = sjson.Set(*content, fmt.Sprintf("%s.celldata.%s."+req.RC, sheetArrIndx, key.String()), rc-req.V.Len)
				if err != nil {
					fmt.Println(err)
				}
			}
			return true
		})
	case "arc": //增加行或列
		ds := req.V.Index
		doc.Get(sheetArrIndx + ".celldata").ForEach(func(key, value gjson.Result) bool {
			var rc = int(value.Get(req.RC).Int())
			if rc < ds {
				return true
			}
			*content, err = sjson.Set(*content, fmt.Sprintf("%s.celldata.%s."+req.RC, sheetArrIndx, key.String()), rc+req.V.Len)
			if err != nil {
				fmt.Println(err)
			}
			return true
		})
	}
}

func updateFilter(reqmsg []byte, content *string) {
	req := new(UpdateFilter)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc          = gjson.Parse(*content)
		sheetArrIndx = ""
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.I) {
			sheetArrIndx = key.String()
			return false
		}
		return true
	})
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}

	var err1, err2 error
	if req.V == nil {
		*content, err1 = sjson.Set(*content, sheetArrIndx+".filter", "")
		*content, err2 = sjson.Set(*content, sheetArrIndx+".filter_select", "")
	} else {
		*content, err1 = sjson.Set(*content, sheetArrIndx+".filter", req.V.Filter)
		*content, err2 = sjson.Set(*content, sheetArrIndx+".filter_select", req.V.FilterSelect)
	}
	if err1 != nil || err2 != nil {
		fmt.Println(err1, err2)
	}
}

func addSheet(reqmsg []byte, content *string) {
	req := new(AddSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	*content, err = sjson.Set(*content, "-1", req.V)
	if err != nil {
		fmt.Println(err)
	}
}

func copySheet(reqmsg []byte, content *string) {
	req := new(CopySheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc          = gjson.Parse(*content)
		sheetArrIndx = ""
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.V.CopyIndex) {
			sheetArrIndx = key.String()
			return false
		}
		return true
	})
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}

	newSheet := doc.Get(sheetArrIndx).String()
	newSheet, err = sjson.Set(newSheet, "name", req.V.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	newSheet, err = sjson.Set(newSheet, "index", req.I)
	if err != nil {
		fmt.Println(err)
		return
	}
	*content, err = sjson.Set(*content, "-1", newSheet)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteSheet(reqmsg []byte, content *string) {
	req := new(DeleteSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc          = gjson.Parse(*content)
		sheetArrIndx = ""
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.V.DeleteIndex) {
			sheetArrIndx = key.String()
			return false
		}
		return true
	})
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}

	*content, err = sjson.Set(*content, fmt.Sprintf("%s.deleted", sheetArrIndx), 1)
	if err != nil {
		fmt.Println(err)
	}
}

func recoverSheet(reqmsg []byte, content *string) {
	req := new(RecoverSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc          = gjson.Parse(*content)
		sheetArrIndx = ""
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == string(req.V.RecoverIndex) {
			sheetArrIndx = key.String()
			return false
		}
		return true
	})
	if sheetArrIndx == "" {
		fmt.Println("Sheet 不存在")
		return
	}

	*content, err = sjson.Delete(*content, fmt.Sprintf("%s.deleted", sheetArrIndx))
	if err != nil {
		fmt.Println(err)
	}
}

func updateSheetOrder(reqmsg []byte, content *string) {
	req := new(UpdateSheetOrder)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	for index, order := range req.V {
		var (
			doc          = gjson.Parse(*content)
			sheetArrIndx = ""
		)
		doc.ForEach(func(key, value gjson.Result) bool {
			if value.Get("index").String() == index {
				sheetArrIndx = key.String()
				return false
			}
			return true
		})
		if sheetArrIndx == "" {
			fmt.Println("Sheet 不存在")
			return
		}

		*content, err = sjson.Set(*content, fmt.Sprintf("%s.order", sheetArrIndx), order)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func hideOrShowSheet(reqmsg []byte, content *string) {
	req := new(HideOrShowSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc = gjson.Parse(*content)
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		*content, err = sjson.Set(*content, key.String()+".status", 0)
		if err != nil {
			fmt.Println(err)
		}
		if value.Get("index").String() == string(req.I) {
			*content, err = sjson.Set(*content, key.String()+".status", 1)
			if err != nil {
				fmt.Println(err)
			}
			*content, err = sjson.Set(*content, key.String()+".hide", req.V)
			if err != nil {
				fmt.Println(err)
			}
		}
		return true
	})
	switch req.Op {
	case "hide":
		curKey := ""
		hasStatus1 := false
		doc.ForEach(func(key, value gjson.Result) bool {
			if value.Get("index").String() == string(req.I) {
				*content, err = sjson.Set(*content, key.String()+".hide", req.V)
				if err != nil {
					fmt.Println(err)
				}
				*content, err = sjson.Set(*content, key.String()+".status", 0)
				if err != nil {
					fmt.Println(err)
				}
			}
			if value.Get("index").String() == req.Cur {
				curKey = key.String()
			}
			if value.Get("status").Int() == 1 {
				hasStatus1 = true
			}
			return true
		})
		if curKey == "" && !hasStatus1 {
			*content, err = sjson.Set(*content, curKey+".status", req.V)
			if err != nil {
				fmt.Println(err)
			}
		}
	case "show":
		doc.ForEach(func(key, value gjson.Result) bool {
			*content, err = sjson.Set(*content, key.String()+".status", 0)
			if err != nil {
				fmt.Println(err)
			}
			if value.Get("index").String() == string(req.I) {
				*content, err = sjson.Set(*content, key.String()+".status", 1)
				if err != nil {
					fmt.Println(err)
				}
				*content, err = sjson.Set(*content, key.String()+".hide", req.V)
				if err != nil {
					fmt.Println(err)
				}
			}
			return true
		})
	}

}

func toggleSheet(reqmsg []byte, content *string) {
	req := new(ToggleSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		doc = gjson.Parse(*content)
	)
	doc.ForEach(func(key, value gjson.Result) bool {
		if value.Get("index").String() == req.V {
			*content, err = sjson.Set(*content, key.String()+".status", 1)
		} else {
			*content, err = sjson.Set(*content, key.String()+".status", 0)
		}
		if err != nil {
			fmt.Println(err)
		}
		return true
	})
}

type CellPos struct {
	R int
	C int
}
