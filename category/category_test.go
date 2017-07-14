package category

import (
	"testing"
)

var totalPagesTest = []struct {
	total    uint
	pageSize uint
	expected uint
}{
	{0, 200, 0},
	{199, 200, 1},
	{200, 200, 1},
	{250, 200, 2},
}

func TestTotalPages(t *testing.T) {

	cat := New()

	for _, tt := range totalPagesTest {
		actual := cat.getTotalPages(tt.total, tt.pageSize)
		if actual != tt.expected {
			t.Fail()
		}
	}

}
