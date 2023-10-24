package page

func (p *PageRequest) GetOffset() int {
	return (p.PageNum - 1) * p.PageSize
}

func (p *PageRequest) GetLimit() int {
	return p.PageSize
}

// NewPageRequest return PageRequest with corresponded value, or default value if args is 0
func NewPageRequest(pageNum int, pageSize int) *PageRequest {

	return &PageRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}

func NewPageResponse(page *PageRequest, totalItem int) *PageResponse {
	var totalPage int
	if page.PageSize > 0 {
		if totalItem%page.PageSize > 0 {
			totalPage = totalItem/page.PageSize + 1
		} else {
			totalPage = totalItem / page.PageSize
		}
	}

	var nextPage int
	if page.PageNum < totalPage {
		nextPage = page.PageNum + 1
	} else {
		nextPage = 0
	}

	return &PageResponse{
		PrevPage:  page.PageNum - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
		TotalItem: totalItem,
	}
}

func (p *PageRequest) IsValid() bool {
	if p.PageNum < 1 || p.PageSize < 1 {
		return false
	}
	return true
}
