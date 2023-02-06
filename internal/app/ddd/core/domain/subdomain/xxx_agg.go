package subdomain

// 领域聚合，通常表现为实体和值对象的聚合，需要有聚合根
import "github.com/go-leo/design-pattern/ddd"

type IXXXAgg ddd.Aggregate[XXXEntity, string]

type XXXAgg struct {
}

func (X XXXAgg) Root() ddd.Entity[XXXEntity, string] {
	return XXXEntity{}
}
