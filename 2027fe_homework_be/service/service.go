package service

type Service struct {
	ImageSrv
	AuthSrv
	UserSrv
	ArticleSrv
	CommentSrv
	ConversationSrv
}

func New() *Service {
	service := &Service{}
	return service
}
