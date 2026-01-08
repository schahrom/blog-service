package models

import (
	"errors"
	"strconv"
)

type PageDto struct {
	Limit  int
	Offset int
	Sort   string
}

func NewPageDto(limit string, offset string, sort string) (PageDto, error) {
	var pagedDto PageDto
	if limit == "" {
		pagedDto.Limit = 10
	} else {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return PageDto{}, err
		}
		pagedDto.Limit = limitInt
	}
	if offset == "" {
		pagedDto.Offset = 0
	} else {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return PageDto{}, err
		}
		pagedDto.Offset = offsetInt
	}
	if sort == "" {
		pagedDto.Sort = "asc"
	} else {
		if sort != "asc" && sort != "desc" {
			return PageDto{}, errors.New("sort must be either 'asc' or 'desc'")
		}
		pagedDto.Sort = sort
	}
	return pagedDto, nil
}
