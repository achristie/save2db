package csv

import (
	"encoding/csv"
	"io"
)

type AssessmentsEncoder struct {
	w *csv.Writer
}

func NewAssessmentsEncoder(w io.Writer) *AssessmentsEncoder {
	enc := AssessmentsEncoder{w: csv.NewWriter(w)}
	return &enc

}
