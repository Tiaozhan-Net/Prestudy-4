package model

import (
	"fe_exam/common"

	"gorm.io/gorm"
)

func Paginate(pager common.PagerForm) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pager.Page <= 0 {
			pager.Page = 1
		}

		switch {
		case pager.Size > 100:
			pager.Size = 100
		case pager.Size <= 0:
			pager.Size = 10
		}

		offset := (pager.Page - 1) * pager.Size
		return db.Offset(offset).Limit(pager.Size)
	}
}
