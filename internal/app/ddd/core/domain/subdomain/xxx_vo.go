package subdomain

// 值对象
import (
	"github.com/go-leo/design-pattern/ddd"
)

type IXXXVo ddd.ValueObject[XXXVo]

type XXXVo struct {
	F1 string
	F2 string
	F3 string
}

func (a XXXVo) SameValueAs(other XXXVo) bool {
	return a.F1 == other.F1 && a.F2 == other.F1 && a.F3 == other.F3
}
