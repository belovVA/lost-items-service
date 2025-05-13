package service

type AnnService struct {
	annRepository AnnRepository
}

func NewAnnService(repo AnnRepository) *AnnService {
	return &AnnService{
		annRepository: repo,
	}
}
