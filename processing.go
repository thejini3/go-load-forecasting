package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	hel "github.com/thejini3/go-helper"
	"gonum.org/v1/plot/plotter"
)

func dailyData1() {

	var csvRows []xDailyData

	parseCsv(csvDailyData, &csvRows)

	var count float64 = 0
	var plotData plotter.XYs

	for _, row := range csvRows {
		plotData = append(plotData, plotter.XY{
			X: count,
			Y: row.Load / 1000,
		})
		count++
	}

	build("build/temp.png", []interface{}{
		"daily data", plotData,
	})
}

func sldcToDailyData() {

	folder := "SLDC_Data"
	outputTxt := "Date,Hour,Min,Load\n"

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".csv" {
			return nil
		}

		var date = strings.ReplaceAll(info.Name(), ".csv", "")
		var rows []xTimeValue
		parseCsv(path, &rows)

		for _, row := range rows {

			if len(row.Time) == 0 || len(row.Value) == 0 {
				continue
			}

			for hour := 0; hour <= 23; hour++ {
				// row.Time 00:00
				// hourMin [01, 00]
				hourMin := strings.Split(row.Time, ":")
				// minute string : 00
				if len(hourMin) != 2 {
					hel.Pl(hourMin)
					panic("Invalid hourMin")
				}
				minStr := hourMin[1]
				cond := timify(hour) == hourMin[0] && (minStr == "00" || minStr == "30")
				if cond {
					var min int
					if minStr == "00" {
						min = 0
					} else if minStr == "30" {
						min = 30
					} else {
						panic("Unknown minute: " + err.Error())
					}
					// outputTxt += fmt.Sprintf("%s,%d,%s\n", date, hour, row.Value)
					outputTxt += fmt.Sprintf("%s,%d,%d,%s\n", date, hour, min, row.Value)
				}
			}
		}

		return err
	})

	if err != nil {
		panic(err)
	}

	if err = hel.StrToFile(folder+"/processed-by-30-min.csv", outputTxt); err != nil {
		panic(err)
	}
}