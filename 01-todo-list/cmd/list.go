package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"text/tabwriter"
)

var all bool = false

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list the todo list",
	Run: func(cmd *cobra.Command, args []string) {
		var file *os.File
		if _, err := os.Stat("data.csv"); err == nil {
			file, _ = os.OpenFile("data.csv", os.O_RDONLY, 0644)
			defer file.Close()
			reader := csv.NewReader(file)
			records, _ := reader.ReadAll()
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
			if all {
				fmt.Fprintln(w, "ID\tTask\tCreated\tDone")
			} else {
				fmt.Fprintln(w, "ID\tTask\tCreated")
			}
			for _, record := range records[1:] {
				id, _ := strconv.ParseUint(record[0], 10, 32)
				task := record[1]
				created := timesince(record[2])
				if all {
					fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", id, task, created, record[3])
				} else if record[3] == "false" {
					fmt.Fprintf(w, "%d\t%s\t%s\n", id, task, created)
				}
			}
			w.Flush()
		}
	},
}

func init() {
	listCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "list completed and uncompleted tasks")
	rootCmd.AddCommand(listCmd)
}
