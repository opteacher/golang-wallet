package managers

type Component interface {
}

const (
	SERVICE = iota
)

var components = map[int]Component {
	SERVICE: new(Service),
}

func GetComponent(name int) (Component, error) {
	return components[name], nil
}