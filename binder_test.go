package echo

import (
	"net/url"
	"testing"
)

func TestBinder_BindQueryParams_EmbeddedStruct(t *testing.T) {
	type CommonParams struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}
	type Level1 struct {
		CommonParams
		Filter string `query:"filter"`
	}
	type SearchRequest struct {
		Level1
		*CommonParams
		Query string `query:"q"`
	}

	b := &DefaultBinder{}
	params := url.Values{
		"q":      []string{"golang"},
		"page":   []string{"3"},
		"limit":  []string{"25"},
		"filter": []string{"active"},
	}
	c := NewContext(params)

	u := new(SearchRequest)
	err := b.BindQueryParams(c, u)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Query != "golang" {
		t.Errorf("expected Query 'golang', got '%s'", u.Query)
	}
	if u.Level1.Page != 3 {
		t.Errorf("expected Level1.Page 3, got %d", u.Level1.Page)
	}
	if u.Level1.Limit != 25 {
		t.Errorf("expected Level1.Limit 25, got %d", u.Level1.Limit)
	}
	if u.Filter != "active" {
		t.Errorf("expected Filter 'active', got '%s'", u.Filter)
	}
	if u.CommonParams == nil {
		t.Fatalf("expected CommonParams pointer to be initialized")
	}
	if u.CommonParams.Page != 3 {
		t.Errorf("expected CommonParams.Page 3, got %d", u.CommonParams.Page)
	}
	if u.CommonParams.Limit != 25 {
		t.Errorf("expected CommonParams.Limit 25, got %d", u.CommonParams.Limit)
	}

	// Test un-initialized pointer when no matching query params
	u2 := new(SearchRequest)
	params2 := url.Values{
		"q": []string{"empty_params"},
	}
	c2 := NewContext(params2)
	err = b.BindQueryParams(c2, u2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u2.CommonParams != nil {
		t.Errorf("expected CommonParams pointer to be nil, got initialized")
	}
}
