package system

type RouterGroup struct {
	BaseRouter
	JwtRouter
	UserRouter
	MenuRouter
	AuthorityRouter
	ApiRouter
}
