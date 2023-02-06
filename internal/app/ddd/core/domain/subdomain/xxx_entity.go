package subdomain

import "github.com/go-leo/design-pattern/ddd"

type IXXXEntity ddd.Entity[XXXEntity, string]

type XXXEntity struct {
	id string
}

func (p XXXEntity) SameIdentityAs(other XXXEntity) bool {
	return p.Identity() == other.Identity()
}

func (p XXXEntity) Identity() string {
	return p.id
}
