package domain

type PageinationOptions struct {
	Page int
	Size int
}

func NewPaginationOptions() PageinationOptions {
	return PageinationOptions{
		Page: 1,
		Size: 100,
	}
}
