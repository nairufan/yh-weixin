package utils

import (
	"encoding/csv"
	"github.com/astaxie/beego/context"
	"fmt"
)

func Write(ctx *context.Context, records [][]string, fileName string) {
	ctx.Output.Header("Content-Type", "text/csv;charset=utf-8")
	ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", fileName))
	writer := csv.NewWriter(ctx.ResponseWriter)
	ctx.ResponseWriter.Write([]byte("\xEF\xBB\xBF"))
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}
	writer.Flush()
}
