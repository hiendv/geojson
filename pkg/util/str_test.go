package util

import (
	"testing"

	"github.com/matryer/is"
)

func TestNormalizeString(t *testing.T) {
	is := is.New(t)

	str := NormalizeString(`
		Cuộc sống vẫn trôi
		Xoay quanh một giấc mộng
		Nâng niu ngày tháng bình yên như vô tận
	`)

	is.Equal(str, `
		Cuoc song van troi
		Xoay quanh mot giac mong
		Nang niu ngay thang binh yen nhu vo tan
	`)
}
