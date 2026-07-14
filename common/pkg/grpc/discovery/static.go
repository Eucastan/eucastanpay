package discovery

import "fmt"

type StaticRegistry struct {
	services map[string]string
}

func NewStaticRegistry(services map[string]string) *StaticRegistry {
	return &StaticRegistry{
		services: services,
	}

}

func (r *StaticRegistry) Resolve(service string) (string, error) {
	addr, ok := r.services[service]
	if !ok {
		return "", fmt.Errorf("service %s not found", service)
	}

	return addr, nil

}
