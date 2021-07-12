package mongodb

import (
	"fmt"

	"github.com/sy-software/minerva-owl/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
)

func formatFilters(filters []ports.Filter) (bson.D, error) {
	result := bson.D{}

	for _, filter := range filters {
		switch filter.Name {
		case "$or", "$and":

			if values, ok := filter.Value.([]ports.Filter); ok {
				formated := []bson.E{}

				for _, v := range values {
					formated = append(formated, formatFilter(v))
				}

				result = append(result, bson.E{
					Key:   filter.Name,
					Value: formated,
				})
			} else {
				return result, fmt.Errorf("invalid value for operator %q: %+v", filter.Name, filter.Value)
			}
		default:
			result = append(result, formatFilter(filter))
		}
	}

	return result, nil
}

func formatFilter(f ports.Filter) bson.E {
	if value, ok := f.Value.(ports.Filter); ok {
		return bson.E{
			Key:   f.Name,
			Value: formatFilter(value),
		}
	} else {
		return bson.E{
			Key:   f.Name,
			Value: f.Value,
		}
	}
}
