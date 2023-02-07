package domainname

// 工厂类，负责复杂领域对象创建，封装细节
type XXXFactory struct {
}

func CreateXXX() *XXXVo {
	return &XXXVo{
		F1: "1",
		F2: "2",
		F3: "3",
	}
}

func FindXXXEntity() *XXXEntity {
	// TODO fild xxx from repository
	return &XXXEntity{}
}
