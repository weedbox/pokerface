package psae

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_MemoryWaitingRoom_Pending(t *testing.T) {

	done := make(chan struct{})

	rto := NewTestRuntimeOptions()
	rto.WaitingRoomMatched = func(p PSAE, players []*Player) {
		assert.Equal(t, 8, len(players))
		done <- struct{}{}
	}

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
		WithWaitingRoom(NewMemoryWaitingRoom(time.Second*5)),
	)
	defer p.Close()

	// Prepare players
	for i := 0; i < 8; i++ {
		player := NewTestPlayer()
		err := p.EnterWaitingRoom(player)
		assert.Nil(t, err)
	}

	<-done
}
