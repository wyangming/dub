package model

type ConAuth struct {
	AuthUrl, AuthName string
	HasChild          bool
	Childs            []*ConAuth
	Id                uint
}
