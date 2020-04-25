package tools

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func WriteCsv(bundlePath string, writer *os.File) error {
	tasks, err := FindTasks(bundlePath)
	if err != nil {
		return fmt.Errorf("cannot write CSV: %v", err.Error())
	}
	csvWriter := csv.NewWriter(writer)
	for _, t := range tasks {
		err = csvWriter.Write([]string{t.Name,
			printTime(t.Staring),
			printTime(t.Running),
			printTime(t.Killed),
			printTime(t.Failed),
			t.ID,
			fmt.Sprintf("%v", t.HasLogs),
			t.DirName})
		if err != nil {
			return fmt.Errorf("cannot write to the CSV output: %v", err.Error())
		}
	}
	csvWriter.Flush()
	return nil
}

func printTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.String()
}
