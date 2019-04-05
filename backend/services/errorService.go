package services

type ErrorService struct{}

func (s *ErrorService) Check(err error) {
	if err != nil {
		panic(err)
	}
}




