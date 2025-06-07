package cmd

import (
	"encoding/csv"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "To delete a tasks with the ids.",
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
		to_delete := make(map[string]bool)
		for _, i := range args {
			to_delete[i] = true
		}
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				println("Error reading file:", err.Error())
				return
			}
			if !to_delete[record[0]] {
				records = append(records, record)
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
	rootCmd.AddCommand(deleteCommand)
}
