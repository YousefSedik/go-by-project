package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
	"time"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to the list",
	Run: func(cmd *cobra.Command, args []string) {
		lastID := 1
		var file *os.File
		var writer *csv.Writer

		if _, err := os.Stat("data.csv"); err == nil {
			// Open file for read/write + appending
			file, _ = os.OpenFile("data.csv", os.O_APPEND|os.O_RDWR, 0644)
			defer file.Close()

			reader := csv.NewReader(file)
			file.Seek(0, 0)

			for {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				if len(record) == 0 || record[0] == "ID" {
					continue
				}
				parsedID, _ := strconv.ParseInt(record[0], 10, 32)
				lastID = int(parsedID) + 1
			}

			writer = csv.NewWriter(file)

		} else if errors.Is(err, os.ErrNotExist) {
			file, _ = os.Create("data.csv")
			defer file.Close()
			writer = csv.NewWriter(file)
			writer.Write([]string{"ID", "Description", "CreatedAt", "IsComplete"})
			fmt.Println("Creating new file: data.csv")
		} else {
			fmt.Printf("Error checking file: %v\n", err)
			return
		}
		// Add new tasks
		for _, value := range args {
			writer.Write([]string{
				strconv.Itoa(lastID),
				value,
				time.Now().Format(time.RFC3339),
				"false",
			})
			fmt.Printf("Added '%s' to the todo list.\n", value)
			lastID++
		}

		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
