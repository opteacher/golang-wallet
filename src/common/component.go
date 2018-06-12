package common

type Component interface {
	Create() error
	IsCreate() bool
}