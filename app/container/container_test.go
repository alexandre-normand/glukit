package container_test

import (
	. "github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/model"
	"testing"
	"time"
)

func TestListReversal(t *testing.T) {
	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	current := NewImmutableList(nil, model.GlucoseRead{model.Timestamp{"", ct.Unix()}, 0})

	for i := 0; i < 10; i++ {
		readTime := ct.Add(time.Duration(i+1) * 30 * time.Minute)
		current = NewImmutableList(current, model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, i + 1})
	}

	r, _ := ReverseList(current)

	for previous, cursor := r.Next, r.Next.Next; cursor.Next != nil; previous, cursor = previous.Next, cursor.Next {
		t.Logf("Current is %d and previous is %d", cursor.Value.(model.GlucoseRead).Value, previous.Value.(model.GlucoseRead).Value)
		if cursor.Value.(model.GlucoseRead).Value <= previous.Value.(model.GlucoseRead).Value {
			t.Errorf("TestListReversal test failed: list in incorrect order: %s", r)
		}
	}
}
