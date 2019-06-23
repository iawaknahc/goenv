package goenv

import (
	"reflect"
	"testing"
)

func TestSortPointers(t *testing.T) {
	input := []pointer{
		pointer{
			pointerPart{
				Name: "HOST",
			},
		},
		pointer{
			pointerPart{
				Name: "WORKERS",
			},
			pointerPart{
				SliceIndex: 0,
			},
			pointerPart{
				Name: "MYNAME",
			},
		},
		pointer{
			pointerPart{
				Name: "WORKERS",
			},
			pointerPart{
				SliceIndex: 1,
			},
			pointerPart{
				Name: "MYNAME",
			},
		},
		pointer{
			pointerPart{
				Name: "COORDINATES",
			},
			pointerPart{
				SliceIndex: 0,
			},
			pointerPart{
				SliceIndex: 0,
			},
		},
		pointer{
			pointerPart{
				Name: "COORDINATES",
			},
			pointerPart{
				SliceIndex: 9,
			},
			pointerPart{
				SliceIndex: 9,
			},
		},
	}

	sortPointers(input)

	actual := make([]string, len(input))
	for i, ptr := range input {
		actual[i] = ptr.String()
	}

	expected := []string{
		"WORKERS_1_MYNAME",
		"WORKERS_0_MYNAME",
		"HOST",
		"COORDINATES_9_9",
		"COORDINATES_0_0",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%q != %q", actual, expected)
	}
}
