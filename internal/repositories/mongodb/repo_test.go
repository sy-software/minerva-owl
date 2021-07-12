package mongodb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sy-software/minerva-owl/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFilters(t *testing.T) {
	t.Run("Test a single field with default operator filter", func(t *testing.T) {
		filters := []ports.Filter{
			{
				Name:  "myfield",
				Value: "myvalue",
			},
		}

		expect := bson.D{
			bson.E{Key: "myfield", Value: "myvalue"},
		}

		got, err := formatFilters(filters)

		if err != nil {
			t.Errorf("Filters failed to format with error: %v", err)
		}

		if !cmp.Equal(expect, got) {
			t.Errorf("Expected filter: %+v. Got: %+v", expect, got)
		}
	})

	t.Run("Test multiple field with default operator filter", func(t *testing.T) {
		filters := []ports.Filter{
			{
				Name:  "myfield1",
				Value: "myvalue1",
			},
			{
				Name:  "myfield2",
				Value: "myvalue2",
			},
		}

		expect := bson.D{
			bson.E{Key: "myfield1", Value: "myvalue1"},
			bson.E{Key: "myfield2", Value: "myvalue2"},
		}

		got, err := formatFilters(filters)

		if err != nil {
			t.Errorf("Filters failed to format with error: %v", err)
		}

		if !cmp.Equal(expect, got) {
			t.Errorf("Expected filter: %+v. Got: %+v", expect, got)
		}
	})

	t.Run("Test filter with operator in value", func(t *testing.T) {
		filters := []ports.Filter{
			{
				Name: "myfield1",
				Value: ports.Filter{
					Name:  "$in",
					Value: []string{"v1", "v2"},
				},
			},
			{
				Name: "myfield2",
				Value: ports.Filter{
					Name:  "$eq",
					Value: 10,
				},
			},
		}

		expect := bson.D{
			bson.E{Key: "myfield1", Value: bson.E{
				Key:   "$in",
				Value: []string{"v1", "v2"},
			}},
			bson.E{Key: "myfield2", Value: bson.E{
				Key:   "$eq",
				Value: 10,
			}},
		}

		got, err := formatFilters(filters)

		if err != nil {
			t.Errorf("Filters failed to format with error: %v", err)
		}

		if !cmp.Equal(expect, got) {
			t.Errorf("Expected filter: %+v. Got: %+v", expect, got)
		}
	})

	t.Run("Test top level operators", func(t *testing.T) {
		filters := []ports.Filter{
			{
				Name: "$or",
				Value: []ports.Filter{
					{
						Name: "myfield1",
						Value: ports.Filter{
							Name:  "$eq",
							Value: 10,
						},
					},
					{
						Name: "myfield2",
						Value: ports.Filter{
							Name:  "$eq",
							Value: 10,
						},
					},
				},
			},
		}

		expect := bson.D{
			bson.E{Key: "$or", Value: []bson.E{
				{
					Key: "myfield1",
					Value: bson.E{
						Key:   "$eq",
						Value: 10,
					},
				},
				{
					Key: "myfield2",
					Value: bson.E{
						Key:   "$eq",
						Value: 10,
					},
				},
			}},
		}

		got, err := formatFilters(filters)

		if err != nil {
			t.Errorf("Filters failed to format with error: %v", err)
		}

		if !cmp.Equal(expect, got) {
			t.Errorf("Expected filter: %+v. Got: %+v", expect, got)
		}
	})
}
