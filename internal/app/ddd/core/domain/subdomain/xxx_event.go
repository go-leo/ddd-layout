package subdomain

// 领域事件
import "github.com/go-leo/design-pattern/ddd"

type IXXXEvent ddd.DomainEvent[XXXEvent]

type XXXEvent struct {
	ID string
}

func (h XXXEvent) SameEventAs(other XXXEvent) bool {
	return h.ID == other.ID
}
