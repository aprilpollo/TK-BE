package gormq

import (
	"fmt"
	"strings"

	"aprilpollo/internal/pkg/query"

	"gorm.io/gorm"
)

// ApplyFilters chains filters and sort onto a *gorm.DB scope without pagination.
// Use this for COUNT queries so total is not capped by LIMIT/OFFSET.
func ApplyFilters(db *gorm.DB, opts query.QueryOptions) *gorm.DB {
	for _, f := range opts.Filters {
		db = applyFilter(db, f)
	}

	if opts.Search != "" && len(opts.SearchFields) > 0 {
		parts := make([]string, 0, len(opts.SearchFields))
		args := make([]interface{}, 0, len(opts.SearchFields))
		for _, field := range opts.SearchFields {
			parts = append(parts, fmt.Sprintf("%s LIKE ?", field))
			args = append(args, "%"+opts.Search+"%")
		}
		db = db.Where("("+strings.Join(parts, " OR ")+")", args...)
	}

	if opts.Sort != "" {
		order := opts.Sort
		if strings.ToUpper(opts.Order) == "DESC" {
			order += " DESC"
		}
		db = db.Order(order)
	}

	return db
}

// ApplyToGorm chains filters, sort, limit, offset onto a *gorm.DB scope.
// The caller is responsible for calling .Find(), .Scan(), etc. afterward.
func ApplyToGorm(db *gorm.DB, opts query.QueryOptions) *gorm.DB {
	db = ApplyFilters(db, opts)

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
