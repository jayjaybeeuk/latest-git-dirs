package output

import (
	"encoding/json"
	"fmt"

	"github.com/jamesburton/latestdirs/internal/model"
)

func PrintTable(results []model.Result) {
	fmt.Printf("%-20s %s\n", "Modified", "Directory")

	for _, r := range results {
		fmt.Printf(
			"%-20s %s\n",
			r.Timestamp.Format("2006-01-02 15:04:05"),
			r.Path,
		)
	}
}

func PrintJSON(results []model.Result) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
