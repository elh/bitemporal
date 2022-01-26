package memory_test

import (
	"sync"
	"testing"

	"github.com/elh/bitempura/memory"
	"github.com/stretchr/testify/require"
)

// This test has no assertions but is meant to trigger data race detector if. When struct fields were unsynchronized
// this failed. Calling all functions is a fast way to suss out conflicts
func TestRace(t *testing.T) {
	db, err := memory.NewDB()
	require.Nil(t, err)

	concurrency := 4
	callCount := 25

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < callCount; i++ {
				_ = db.Set("a", id)
				_, _ = db.Get("a")
				_ = db.Delete("a")
				_, _ = db.List()
				_, _ = db.History("a")
				db.SetNow(t0)
			}
		}(i)
	}

	wg.Wait()
}
