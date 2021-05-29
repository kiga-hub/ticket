package excel

import (
	"fmt"
	"path"
	"path/filepath"

    "github.com/tealeg/xlsx"
)

type Cell struct {
	Value	string
}

type Row struct {
	Cells 	[]*Cell
}

type Sheet struct {
	Name		string
	Header 		[]string
	Rows		[]*Row
}

func IsExcel(fullname string) bool {
	_, fileName := filepath.Split(fullname)
	suffix := path.Ext(fileName)
	if suffix == ".xlsx" {
		return true
	}
	return false
}

func ReadExcel(excelFileName string) *Sheet {
    xlFile, err := xlsx.OpenFile(excelFileName)
    if err != nil {
		fmt.Printf("open failed: %s\n", err)
		return nil
	}
	sh := new(Sheet)
    for _, sheet := range xlFile.Sheets {
		sh.Name = sheet.Name
		header_flag := true
        for _, row := range sheet.Rows {
			r := new(Row)
            for _, cell := range row.Cells {
                text := cell.String()
				if header_flag {
					sh.Header = append(sh.Header, text)
					continue
				}
				c := new(Cell)
				c.Value = text
				r.Cells = append(r.Cells, c)
			}
			sh.Rows = append(sh.Rows, r)
			header_flag = false
		}
		break
	}
	return sh
}

func WriteExcel(data *Sheet, excelFileName string) bool {
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	
	var file *xlsx.File
	file = xlsx.NewFile()

    var sheet *xlsx.Sheet
    sheet, err = file.AddSheet(data.Name)
    if err != nil {
		fmt.Printf(err.Error())
		return false
	}

	if len(data.Header) > 0 {
		var header *xlsx.Row
		header = sheet.AddRow()
		header.SetHeightCM(1)
		for _, h := range data.Header {
			cell = header.AddCell()
			cell.Value = h
		}
	}

	for _, r := range data.Rows {
		row = sheet.AddRow()
		row.SetHeightCM(1)
		for _, c := range r.Cells {
			cell = row.AddCell()
			cell.Value = c.Value
		}
	}

	err = file.Save(excelFileName)
    if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	return true
}