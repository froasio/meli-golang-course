package category

import (
	"testing"
)

var totalPagesTest = []struct {
	total    uint
	expected uint
}{
	{0, 0},
	{199, 1},
	{200, 1},
	{250, 2},
}

func TestTotalPages(t *testing.T) {

	cat := New()

	for _, tt := range totalPagesTest {
		actual := cat.getTotalPages(tt.total)
		if actual != tt.expected {
			t.Fail()
		}
	}

}
