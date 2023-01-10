package category_test

import (
	"ddo/category"
	"testing"
)

func TestToOperation(t *testing.T) {
	t.Parallel()

	want := []string{"validate", "create --what-if", "create"}
	categories := category.Categories()

	for i, c := range categories {
		got, err := c.ToOperation()
		if err != nil {
			t.Errorf("ToOperation() returned error: %v", err)
		}
		if got != want[i] {
			t.Errorf("ToOperation() = %v, want %v", got, want[i])
		}
	}

}

func TestToOperationInvalid(t *testing.T) {
	t.Parallel()

    if _, err := category.Category(999).ToOperation(); err == nil {
		t.Fatal("want error for invalid category, got nil")
	}
}
