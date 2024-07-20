package ag

import (
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

//go:embed example-response.html
var exampleHTML string

//go:embed example-empty-response.html
var emptyExampleHTML string

func TestParseScoreSimple(t *testing.T) {
	score := "Toán:   4.40   Ngữ văn:   7.25   Lịch sử:   7.25   Địa lí:   5.50 "
	parsed := parseScore(score)
	assert.NotNil(t, parsed.Toan, "Toan can't be nil")
	assert.NotNil(t, parsed.Van, "Van can't be nil")
	assert.NotNil(t, parsed.Su, "Su can't be nil")
	assert.NotNil(t, parsed.Dia, "Dia can't be nil")
	assert.Equal(t, float32(4.40), *parsed.Toan)
	assert.Equal(t, float32(7.25), *parsed.Van)
	assert.Equal(t, float32(7.25), *parsed.Su)
	assert.Equal(t, float32(5.50), *parsed.Dia)
	assert.Nil(t, parsed.Anh)
	assert.Nil(t, parsed.Ly)
	assert.Nil(t, parsed.Hoa)
	assert.Nil(t, parsed.Sinh)
	assert.Nil(t, parsed.GDCD)
	assert.Nil(t, parsed.Trung)
	assert.Nil(t, parsed.Nhat)
	assert.Nil(t, parsed.Phap)
	assert.Nil(t, parsed.KHTN)
	assert.Nil(t, parsed.KHXH)
}

func TestParseKHTN(t *testing.T) {
	score := `            Toán: 8.80 Ngữ văn: 8.75 Vật lí: 8.75 Hóa học: 5.75 Sinh học: 6.75
            KHTN: 7.08 Tiếng Anh: 9.80`
	parsed := parseScore(score)
	assert.Nil(t, parsed.KHXH)
	assert.Nil(t, parsed.GDCD)
	assert.Nil(t, parsed.Trung)
	assert.Nil(t, parsed.Nhat)
	assert.Nil(t, parsed.Phap)

	assert.NotNil(t, parsed.Toan)
	assert.NotNil(t, parsed.Van)
	assert.NotNil(t, parsed.Ly)
	assert.NotNil(t, parsed.Hoa)
	assert.NotNil(t, parsed.Sinh)
	assert.NotNil(t, parsed.KHTN)
	assert.NotNil(t, parsed.Anh)
	assert.Equal(t, float32(8.80), *parsed.Toan)
	assert.Equal(t, float32(8.75), *parsed.Van)
	assert.Equal(t, float32(8.75), *parsed.Ly)
	assert.Equal(t, float32(5.75), *parsed.Hoa)
	assert.Equal(t, float32(6.75), *parsed.Sinh)
	assert.Equal(t, float32(7.08), *parsed.KHTN)
	assert.Equal(t, float32(9.80), *parsed.Anh)
}

func TestParseHTMLTable(t *testing.T) {
	parsed, err := parseTableFromHtmlByte([]byte(exampleHTML))
	assert.NoError(t, err)
	assert.NotNil(t, parsed)

	assert.Equal(t, parsed.Name, "TRỊNH NGỌC ÁI")
	assert.Equal(t, parsed.ID, 51000001)
	assert.Equal(t, parsed.Birthday, "15/11/2005")

	assert.Nil(t, parsed.Ly)
	assert.Nil(t, parsed.Sinh)
	assert.Nil(t, parsed.Hoa)
	assert.Nil(t, parsed.GDCD)
	assert.Nil(t, parsed.Trung)
	assert.Nil(t, parsed.Nhat)
	assert.Nil(t, parsed.Phap)
	assert.Nil(t, parsed.KHTN)
	assert.Nil(t, parsed.KHXH)

	assert.NotNil(t, parsed.Toan)
	assert.NotNil(t, parsed.Van)
	assert.NotNil(t, parsed.Su)
	assert.NotNil(t, parsed.Dia)
	assert.Equal(t, *parsed.Toan, float32(4.40))
	assert.Equal(t, *parsed.Van, float32(7.25))
	assert.Equal(t, *parsed.Su, float32(7.25))
	assert.Equal(t, *parsed.Dia, float32(5.50))
}
