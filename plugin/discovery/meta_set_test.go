package discovery

import (
	"fmt"
	"testing"
)

func TestPluginMetaSetManipulation(t *testing.T) {
	metas := []PluginMeta{
		{
			Name:    "foo",
			Version: "1.0.0",
			Path:    "test-foo",
		},
		{
			Name:    "bar",
			Version: "2.0.0",
			Path:    "test-bar",
		},
		{
			Name:    "baz",
			Version: "2.0.0",
			Path:    "test-bar",
		},
	}
	s := make(PluginMetaSet)

	if count := s.Count(); count != 0 {
		t.Fatalf("set has Count %d before any items added", count)
	}

	// Can we add metas?
	for _, p := range metas {
		s.Add(p)
		if !s.Has(p) {
			t.Fatalf("%q not in set after adding it", p.Name)
		}
	}

	if got, want := s.Count(), len(metas); got != want {
		t.Fatalf("set has Count %d after all items added; want %d", got, want)
	}

	// Can we still retrieve earlier ones after we added later ones?
	for _, p := range metas {
		if !s.Has(p) {
			t.Fatalf("%q not in set after all adds", p.Name)
		}
	}

	// Can we remove metas?
	for _, p := range metas {
		s.Remove(p)
		if s.Has(p) {
			t.Fatalf("%q still in set after removing it", p.Name)
		}
	}

	if count := s.Count(); count != 0 {
		t.Fatalf("set has Count %d after all items removed", count)
	}
}

func TestPluginMetaSetValidateVersions(t *testing.T) {
	metas := []PluginMeta{
		{
			Name:    "foo",
			Version: "1.0.0",
			Path:    "test-foo",
		},
		{
			Name:    "bar",
			Version: "0.0.1",
			Path:    "test-bar",
		},
		{
			Name:    "baz",
			Version: "bananas",
			Path:    "test-bar",
		},
	}
	s := make(PluginMetaSet)

	for _, p := range metas {
		s.Add(p)
	}

	valid, invalid := s.ValidateVersions()
	if count := valid.Count(); count != 2 {
		t.Errorf("valid set has %d metas; want 2", count)
	}
	if count := invalid.Count(); count != 1 {
		t.Errorf("valid set has %d metas; want 1", count)
	}

	if !valid.Has(metas[0]) {
		t.Errorf("'foo' not in valid set")
	}
	if !valid.Has(metas[1]) {
		t.Errorf("'bar' not in valid set")
	}
	if !invalid.Has(metas[2]) {
		t.Errorf("'baz' not in invalid set")
	}

	if invalid.Has(metas[0]) {
		t.Errorf("'foo' in invalid set")
	}
	if invalid.Has(metas[1]) {
		t.Errorf("'bar' in invalid set")
	}
	if valid.Has(metas[2]) {
		t.Errorf("'baz' in valid set")
	}

}

func TestPluginMetaSetWithName(t *testing.T) {
	tests := []struct {
		metas     []PluginMeta
		name      string
		wantCount int
	}{
		{
			[]PluginMeta{},
			"foo",
			0,
		},
		{
			[]PluginMeta{
				{
					Name:    "foo",
					Version: "0.0.1",
					Path:    "foo",
				},
			},
			"foo",
			1,
		},
		{
			[]PluginMeta{
				{
					Name:    "foo",
					Version: "0.0.1",
					Path:    "foo",
				},
			},
			"bar",
			0,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test%02d", i), func(t *testing.T) {
			s := make(PluginMetaSet)
			for _, p := range test.metas {
				s.Add(p)
			}
			filtered := s.WithName(test.name)
			if gotCount := filtered.Count(); gotCount != test.wantCount {
				t.Errorf("got count %d in %#v; want %d", gotCount, filtered, test.wantCount)
			}
		})
	}
}
