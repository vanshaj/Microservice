package add

import "testing"

func TestAdd(t *testing.T) {
	testCases := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{name: "Two positive values",
			a:        1,
			b:        2,
			expected: 4},
		{name: "One positive one negative",
			a:        1,
			b:        -1,
			expected: 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := add(tc.a, tc.b)
			if actual != tc.expected {
				t.Fatalf("Expected %d, got %d", tc.expected, actual)
			}
		})
	}
}
