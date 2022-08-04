package skiplist

import (
	"testing"
)

func TestInit(t *testing.T) {
	var cmp *CmpInstanceInt
	sl, _ := New(cmp)
	t.Log(sl)
}
