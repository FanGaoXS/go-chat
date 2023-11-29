package hub

import "fangaoxs.com/go-chat/environment"

type Hub interface{}

func New(env environment.Env) (Hub, error) {
	return &hub{}, nil
}

type hub struct{}
