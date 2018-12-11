package shader

import (
	"fmt"
	"testing"
)

func TestDefineList(t *testing.T) {
	dl := NewDefineList()
	dl.Add(HAS_NORMAL)
	dl.Add(HAS_COORD_0)
	dl.Add(HAS_NORMAL)
	fmt.Println(dl.Condition(NewDefineList(HAS_NORMAL, HAS_COORD_0, HAS_BASECOLORMAP)))
	fmt.Println(dl)
}
