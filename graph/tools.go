package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func toIntWithDefault(i *int, defaultVaue int) int {
	if i == nil {
		return defaultVaue
	}

	return *i
}

func IsFieldRequested(ctx context.Context, field string) bool {
	for _, f := range getPreloads(ctx) {
		if f == field {
			return true
		}
	}
	return false
}

// Note: got this from here: https://github.com/99designs/gqlgen/blob/7dd971c871c0b0159ad26c9bf3095a8ba3780402/docs/content/reference/field-collection.md
func getPreloads(ctx context.Context) []string {
	return getNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func getNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string) (preloads []string) {
	for _, column := range fields {
		prefixColumn := getPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(preloads, getNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...)
	}
	return
}

func getPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
