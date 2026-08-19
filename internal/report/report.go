package report

import (
	"encoding/json"
	"fmt"
	"io"
)

type Level string

const (
	Pass Level = "pass"
	Warn Level = "warn"
	Fail Level = "fail"
)

type Finding struct {
	Level   Level  `json:"level"`
	Check   string `json:"check"`
	Message string `json:"message"`
}

type Report struct {
	Repository string    `json:"repository"`
	Findings   []Finding `json:"findings"`
}

func (r Report) HasFailures() bool {
	for _, finding := range r.Findings {
		if finding.Level == Fail {
			return true
		}
	}
	return false
}

func (r Report) Write(w io.Writer, asJSON bool) error {
	if asJSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(r)
	}
	for _, finding := range r.Findings {
		if _, err := fmt.Fprintf(w, "%-4s  %-24s %s\n", finding.Level, finding.Check, finding.Message); err != nil {
			return err
		}
	}
	return nil
}
