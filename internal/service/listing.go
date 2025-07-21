package service

type ListingService struct {
}

func NewListingService() *ListingService {
	return &ListingService{}
}

func (s *ListingService) CreateListing() error {
	return nil
}

func (s *ListingService) GetListings() error {
	return nil
}
