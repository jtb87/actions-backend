package entities

import "testing"

func TestCategoryValidate(t *testing.T) {
	c := Category{Name: "haircut"}
	if got := c.Validate(); got != nil {
		t.Errorf("Category.Validate() = %q, want %v", got, nil)
	}
	c.Name = ""
	if got := c.Validate(); got == nil {
		t.Errorf("Category.Validate() = %q, want %v", got, nil)
	}

	c.Name = "name is way too long to be displayed on a page"
	val := "week"
	c.IntervalType = &val
	if got := c.Validate(); got == nil {
		t.Errorf("Category.Validate() = %q, want %v", got, nil)
	}

}
