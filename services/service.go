package services

// Service represents one service that is to be checked by a checker
type Service struct {
	Name        string
	CheckerName string
	Config      map[string]interface{}
	Interval    int
}

// GetServices returns all registered services to be checked
func GetServices() []*Service {
	websiteServiceConfig := map[string]interface{}{"URL": "max-heidinger.de", "expectedResponse": 200}
	websiteService := &Service{Name: "Own Website", CheckerName: "HTTPGetChecker", Config: websiteServiceConfig, Interval: 5}

	return []*Service{websiteService}
}
