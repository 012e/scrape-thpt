package vietnamnet

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed example-response.html
var exampleHTML string

func TestParse(t *testing.T) {
	student, err := parseTableFromHTML(exampleHTML)
	assert.NoError(t, err)

	assert.Nil(t, student.KHTN)
	assert.Nil(t, student.KHXH)
	assert.Nil(t, student.GDCD)
	assert.Nil(t, student.Phap)
	assert.Nil(t, student.Trung)
	assert.Nil(t, student.Nhat)

	expToan := float32(7.80)
	expLy := float32(6.50)
	expHoa := float32(8.75)
	expSinh := float32(7.75)
	expVan := float32(7.50)
	expAnh := float32(7.60)
	assert.NotNil(t, student)
	assert.NotNil(t, student.Toan)
	assert.NotNil(t, student.Ly)
	assert.NotNil(t, student.Hoa)
	assert.NotNil(t, student.Sinh)
	assert.NotNil(t, student.Van)
	assert.NotNil(t, student.Anh)
	assert.Equal(t, expToan, *student.Toan)
	assert.Equal(t, expLy, *student.Ly)
	assert.Equal(t, expHoa, *student.Hoa)
	assert.Equal(t, expSinh, *student.Sinh)
	assert.Equal(t, expVan, *student.Van)
	assert.Equal(t, expAnh, *student.Anh)
}
