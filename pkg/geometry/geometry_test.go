package geometry

import (
	"fmt"
	"testing"
)

func TestCircleRectangle_Overlaps(t *testing.T) {
	tests := []struct {
		circle *Circle
		rect   *Rectangle
		result bool
	}{
		{NewCircle(3, 5, 2), NewRectangle(8, 2, 4, 2), false},
		{NewCircle(3, 5, 2), NewRectangle(8, 7, 4, 2), false},
		{NewCircle(3, 5, 2), NewRectangle(-1, 7, 4, 2), false},
		{NewCircle(3, 5, 2), NewRectangle(-1, 5, 4, 2), false},
		{NewCircle(3, 5, 2), NewRectangle(6, 5, 4, 2), true},
		{NewCircle(3, 5, 2), NewRectangle(6, 7, 4, 2), true},
		{NewCircle(3, 5, 2), NewRectangle(0, 7, 4, 2), true},
		{NewCircle(3, 5, 2), NewRectangle(0, 7, 4, 2), true},
		{NewCircle(-1, 5, 2), NewRectangle(0, 7, 4, 2), true},
		{NewCircle(-1, 5, 2), NewRectangle(1, 5, 2, 6), true},
	}
	for _, test := range tests {
		testName := fmt.Sprintf(
			"Circle(%.1f,%.1f)/Rectangle(%.1f,%.1f):%v",
			test.circle.center.X,
			test.circle.center.Y,
			test.rect.center.X,
			test.rect.center.Y,
			test.result,
		)

		t.Run(testName, func(t *testing.T) {
			if test.circle.Overlaps(test.rect) != test.result {
				t.Fatal("circle is expected to overlap rectangle")
			}

			if test.rect.Overlaps(test.circle) != test.result {
				t.Fatal("rectangle is expected to overlap circle")
			}
		})
	}
}

func TestCircleCircle_Overlaps(t *testing.T) {
	tests := []struct {
		c1     *Circle
		c2     *Circle
		result bool
	}{
		{NewCircle(3, 5, 2), NewCircle(3, -1, 2), false},
		{NewCircle(3, 5, 2), NewCircle(3, 0, 2), false},
		{NewCircle(3, 5, 2), NewCircle(3, 1, 2), false},
		{NewCircle(3, 5, 2), NewCircle(3, 2, 2), true},
		{NewCircle(3, 5, 2), NewCircle(3, 3, 2), true},
		{NewCircle(3, 5, 2), NewCircle(3, 4, 2), true},
		{NewCircle(3, 5, 2), NewCircle(3, 5, 2), true},
		{NewCircle(3, 5, 2), NewCircle(3, 7, 2), true},
		{NewCircle(3, 5, 2), NewCircle(3, 8, 2), true},
		{NewCircle(3, 5, 2), NewCircle(3, 9, 2), false},
		{NewCircle(3, 5, 2), NewCircle(2, 5, 2), true},
		{NewCircle(3, 5, 2), NewCircle(1, 5, 2), true},
		{NewCircle(3, 5, 2), NewCircle(-1, 5, 2), false},
		{NewCircle(3, 5, 2), NewCircle(4, 5, 2), true},
		{NewCircle(3, 5, 2), NewCircle(5, 5, 2), true},
		{NewCircle(3, 5, 2), NewCircle(7, 5, 2), false},
		{NewCircle(3, 5, 2), NewCircle(4, 6, 2), true},
		{NewCircle(3, 5, 2), NewCircle(5, 7, 2), true},
		{NewCircle(3, 5, 2), NewCircle(7, 9, 2), false},
	}
	for _, test := range tests {
		testName := fmt.Sprintf(
			"Circle(%.1f,%.1f)/Circle(%.1f,%.1f):%v",
			test.c1.center.X,
			test.c1.center.Y,
			test.c2.center.X,
			test.c2.center.Y,
			test.result,
		)

		t.Run(testName, func(t *testing.T) {
			if test.c1.Overlaps(test.c2) != test.result {
				t.Fatal("circle is expected to overlap rectangle")
			}

			if test.c1.Overlaps(test.c2) != test.result {
				t.Fatal("rectangle is expected to overlap circle")
			}
		})
	}
}
