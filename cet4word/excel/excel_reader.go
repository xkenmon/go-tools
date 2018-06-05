package excel

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

func Reader(filename string, lineNum int) ([]string, error) {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	lineStr := strconv.Itoa(lineNum)
	word := xlsx.GetCellValue("Table 1", "A"+lineStr)
	pronu := xlsx.GetCellValue("Table 1", "B"+lineStr)
	mean := xlsx.GetCellValue("Table 1", "C"+lineStr)
	return []string{word, pronu, mean}, nil
}
