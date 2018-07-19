package services

// Service represents one service that is to be checked by a checker
type Service struct {
	Name        string
	CheckerName string
	Config      map[string]interface{}
}
