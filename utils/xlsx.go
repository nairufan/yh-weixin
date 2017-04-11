package utils

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/astaxie/beego/context"
)

func WriteXlsx(ctx *context.Context, records [][]string, fileName string) {
	ctx.Output.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=utf-8")
	ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", fileName))

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		panic(err)
	}
	for _, record := range records {
		row := sheet.AddRow()
		for _, val := range record {
			cell := row.AddCell()
			cell.SetString(val)
		}
	}
	sheet.SetColWidth(0, 0, 10)
	sheet.SetColWidth(1, 1, 10)
	sheet.SetColWidth(2, 2, 15)
	sheet.SetColWidth(3, 3, 30)
	sheet.SetColWidth(4, 4, 30)
	sheet.SetColWidth(6, 6, 30)
	sheet.SetColWidth(7, 7, 30)
	err = file.Write(ctx.ResponseWriter)
	if err != nil {
		panic(err)
	}

}
