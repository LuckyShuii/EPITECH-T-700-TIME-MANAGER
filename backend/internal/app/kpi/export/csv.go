package export

import (
	"encoding/csv"
	"log"
	"os"
)

func ExportCSV(headers []string, rows [][]string, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating CSV file at %s: %v", filePath, err)
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(headers); err != nil {
		return err
	}

	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}
