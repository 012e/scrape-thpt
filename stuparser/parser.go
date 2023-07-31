package stuparser

import (
	"context"
	"net/http"

	"github.com/nfx/go-htmltable"
)

func init() {
	// htmltable.Logger = func(_ context.Context, msg string, fields ...any) {
	// 	log.Info("")
	// }
	htmltable.Logger = func(_ context.Context, _ string, _ ...any) {}
}

func ParseTableFromHtmlByte(b []byte) (Student, error) {
	student, err := htmltable.NewSliceFromString[Student](string(b))
	if err != nil {
		return Student{}, err
	}
	return student[0], nil // there is only one student per table
}

func TableFromResponse(r *http.Response) (Student, error) {
	student, err := htmltable.NewSliceFromResponse[Student](r)
	if err != nil {
		return Student{}, err
	}
	return student[0], nil // there is only one student per table
}
