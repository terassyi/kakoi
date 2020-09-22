package infra

type Infrastructure interface {
	WorkDir() string
	Provider() string
}
