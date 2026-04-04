package gormq

import (
	"fmt"
	"strings"

	"aprilpollo/internal/pkg/query"
	"gorm.io/gorm"
)

// ApplyToGorm chains filters, sort, limit, offset onto a *gorm.DB scope.
// The caller is responsible for calling .Find(), .Scan(), etc. afterward.
func ApplyToGorm(db *gorm.DB, opts query.QueryOptions) *gorm.DB {
	for _, f := range opts.Filters {
		db = applyFilter(db, f)
	}

	if opts.Sort != "" {
		order := opts.Sort
		if strings.ToUpper(opts.Order) == "DESC" {
			order += " DESC"
		}
		db = db.Order(order)
	}

	if opts.Limit > 0 {
		db = db.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		db = db.Offset(opts.Offset)
	}

	return db
}

func applyFilter(db *gorm.DB, f query.Filter) *gorm.DB {
	col := f.Field

	switch f.Operator {
	case "IS NULL":
		return db.Where(fmt.Sprintf("%s IS NULL", col))
	case "IS NOT NULL":
		return db.Where(fmt.Sprintf("%s IS NOT NULL", col))
	case "IN":
		return db.Where(fmt.Sprintf("%s IN ?", col), f.Value)
	case "NOT IN":
		return db.Where(fmt.Sprintf("%s NOT IN ?", col), f.Value)
	default:
		return db.Where(fmt.Sprintf("%s %s ?", col, f.Operator), f.Value)
	}
}
