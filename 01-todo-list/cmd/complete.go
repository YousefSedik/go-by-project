package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "To mark a task as completed",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.OpenFile("data.csv", os.O_RDWR, 0644)
		if err != nil {
			if os.IsNotExist(err) {
				println("File does not exist. Please add tasks first.")
				return
			}
			println("Error opening file:", err.Error())
			return
		}
		records := make([][]string, 0)
		reader := csv.NewReader(file)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				println("Error reading file:", err.Error())
				return
			}
			records = append(records, record)
		}
		file.Seek(0, 0)
		defer file.Close()
		// mark tasks as completed
		for _, i := range args {
			is_found := false
			for _, record := range records {
				if record[0] == i {
					fmt.Printf("Marking Id: %s as completed. \n", i)
					record[3] = "true"
					is_found = true
					break
				}
			}
			if !is_found {
				fmt.Println("Unable to find Id: ", i)
			}
		}
		file.Close()
		file, _ = os.Create("data.csv")
		defer file.Close()
		writer := csv.NewWriter(file)
		writer.WriteAll(records)
		writer.Flush()
		file.Close()
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
