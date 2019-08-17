package types

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	sc := NewUniqueRandomIDSource(32)
	for i := 0; i < 5; i++ {
		id := sc.Get()
		fmt.Println(id)
	}
}
