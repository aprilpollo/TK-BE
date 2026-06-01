package query

import (
	"strconv"
	"strings"
)

const defaultLimit = 10

// Parse builds a QueryOptions from Fiber query params (map[string]string via c.Queries())
func Parse(params map[string]string) (QueryOptions, error) {
	// Convert map[string]string → map[string][]string for ParseFilters
	multi := make(map[string][]string, len(params))
	for k, v := range params {
		multi[k] = []string{v}
	}

	opts := QueryOptions{
		Filters: ParseFilters(multi),
		Limit:   defaultLimit,
		Order:   "ASC",
	}

	if v, ok := params["_q"]; ok {
		opts.Search = strings.TrimSpace(v)
	}
	if v, ok := params["_search_fields"]; ok {
		for _, field := range strings.Split(v, ",") {
			f := strings.TrimSpace(field)
			if f != "" && isValidIdentifier(f) {
				opts.SearchFields = append(opts.SearchFields, f)
			}
		}
	}

	if v, ok := params["_sort"]; ok && isValidIdentifier(v) {
		opts.Sort = v
	}
	if v := strings.ToUpper(params["_order"]); v == "DESC" {
		opts.Order = "DESC"
	}
	if v, ok := params["_limit"]; ok {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			opts.Limit = n
		}
	}
	if v, ok := params["_offset"]; ok {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			opts.Offset = n
		}
	}
	// _page overrides _offset
	if v, ok := params["_page"]; ok {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 {
			opts.Offset = (n - 1) * opts.Limit
		}
	}

	return opts, nil
}
