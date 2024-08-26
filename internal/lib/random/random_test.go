package random

import "testing"


// Runs random strung generatir function twice
// and check for different results.
func TestRandStr(t *testing.T){
	size := 5

	res1 := RandStr(size)

	res2 := RandStr(size)

	if res1 == res2 {
		t.Error("Generated strings are the same, should be different")
	}
}