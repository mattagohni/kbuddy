package internal

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		in                     string
		wantedGroupVersionKind string
	}{
		"yaml1": {in: `---
apiVersion: apps/v1
kind: Deployment`, wantedGroupVersionKind: "apps/v1, Kind=Deployment"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if obj, _, _ := Parse(tc.in); obj.GetObjectKind().GroupVersionKind().String() != tc.wantedGroupVersionKind {
				t.Errorf("Parse(%s) = '%s', want '%s'", tc.in, obj.GetObjectKind().GroupVersionKind().String(), tc.wantedGroupVersionKind)
			}
		})
	}
}
