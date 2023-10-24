package page

type PageRequest struct {
	PageNum  int // Page number which want to be accessed, start from 1
	PageSize int // Indicates item count in a single page
}

type PageResponse struct {
	PrevPage  int `json:"prevPage"`  // If previous page not exists, value must 0
	NextPage  int `json:"nextPage"`  // If next page not exists, value must 0
	TotalPage int `json:"totalPage"` // Indicates total page based on total item divide by page size
	TotalItem int `json:"totalItem"` // Indicates total item
}
