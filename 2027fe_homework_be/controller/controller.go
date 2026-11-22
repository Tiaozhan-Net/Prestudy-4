package controller

type Controller struct {
	Image
	Auth
	User
	Article
	Comment
	Conversation
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}
