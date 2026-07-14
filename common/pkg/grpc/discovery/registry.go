package discovery

type Registry interface {
	Resolve(service string) (string, error)
}
