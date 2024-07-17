package ag

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/012e/scrape-thptscore/models"
	"github.com/nfx/go-htmltable"
)

type rawResponse struct {
	ID             int    `json:"id" header:"Số báo danh"`
	Name           string `json:"name" header:"Tên"`
	UnparsedResult string `json:"unparsedResult" header:"Điểm chi tiết"`
	BirthdayStr    string `json:"birthdayStr" header:"Ngày sinh"`
	Gender         string `json:"gender" header:"Giới tính"`
	CMND           string `json:"cmnd" header:"Số CMND"`
}

type score struct {
	Toan  *float32
	Ly    *float32
	Hoa   *float32
	Van   *float32
	Anh   *float32
	Sinh  *float32
	Su    *float32
	Dia   *float32
	GDCD  *float32
	Trung *float32
	Nhat  *float32
	Phap  *float32
	KHTN  *float32
	KHXH  *float32
}

const (
	anh   = "Tiếng Anh"
	toan  = "Toán"
	ly    = "Vật lí"
	hoa   = "Hóa học"
	van   = "Ngữ văn"
	sinh  = "Sinh học"
	su    = "Lịch sử"
	dia   = "Địa lí"
	gdcd  = "GDCD"
	nhat  = "Tiếng Nhật"
	phap  = "Tiếng Pháp"
	trung = "Tiếng Trung"
	khxh  = "KHXH"
	khtn  = "KHTN"
)

func parseScore(s string) score {
	scores := make(map[string]*float32)
	scores[anh] = nil
	scores[toan] = nil
	scores[ly] = nil
	scores[hoa] = nil
	scores[van] = nil
	scores[sinh] = nil
	scores[su] = nil
	scores[dia] = nil
	scores[gdcd] = nil
	scores[nhat] = nil
	scores[phap] = nil
	scores[trung] = nil
	scores[khxh] = nil
	scores[khtn] = nil

	r := regexp.MustCompile(`(Toán|Ngữ văn|Lịch sử|Địa lí|Vật lí|Lịch sử|Hóa học|Sinh học|Tiếng Anh|Tiếng Trung|Tiếng Nhật|Tiếng Pháp|KHTN|KHXH|GDCD):\s*(\d+\.\d+)`)
	result := r.FindAllStringSubmatch(s, -1)
	for _, match := range result {

		parsed, err := strconv.ParseFloat(match[2], 32)
		if err != nil {
			panic(fmt.Sprintf("invalid score string: %v", err))
		}
		parsed32 := float32(parsed)

		scores[match[1]] = &parsed32
	}

	return score{
		Toan:  scores[toan],
		Ly:    scores[ly],
		Hoa:   scores[hoa],
		Van:   scores[van],
		Anh:   scores[anh],
		Sinh:  scores[sinh],
		Su:    scores[su],
		Dia:   scores[dia],
		GDCD:  scores[gdcd],
		Trung: scores[trung],
		Nhat:  scores[nhat],
		Phap:  scores[phap],
		KHTN:  scores[khtn],
		KHXH:  scores[khxh],
	}
}

func (r rawResponse) toStudent() (*models.Student, error) {
	score := parseScore(r.UnparsedResult)
	return &models.Student{
		SBD:      r.ID,
		Name:     r.Name,
		Gender:   r.Gender,
		Birthday: r.BirthdayStr,
		CMND:     r.CMND,
		Toan:     score.Toan,
		Ly:       score.Ly,
		Hoa:      score.Hoa,
		Van:      score.Van,
		Anh:      score.Anh,
		Sinh:     score.Sinh,
		Su:       score.Su,
		Dia:      score.Dia,
		GDCD:     score.GDCD,
		Trung:    score.Trung,
		Nhat:     score.Nhat,
		Phap:     score.Phap,
		KHTN:     score.KHTN,
		KHXH:     score.KHXH,
	}, nil
}

func ParseTableFromHtmlByte(b []byte) (*models.Student, error) {
	response, err := htmltable.NewSliceFromString[rawResponse](string(b))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html from %s: %v", string(b), err)
	}
	return response[0].toStudent()
}

func TableFromResponse(r *http.Response) (*models.Student, error) {
	response, err := htmltable.NewSliceFromResponse[rawResponse](r)
	if err != nil {
		return nil, err
	}
	return response[0].toStudent()
}
