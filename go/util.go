package main

// GetValidPagination pagination
func GetValidPagination(total, offset, limit int) (startIndex, endIndex int) {
	// no pagination
	if limit == 0 {
		return 0, total
	}

	// out of range
	if limit < 0 || offset < 0 || offset > total {
		return 0, 0
	}

	startIndex = offset
	endIndex = startIndex + limit

	if endIndex > total {
		endIndex = total
	}

	return startIndex, endIndex
}