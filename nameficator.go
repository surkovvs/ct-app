package ctapp

import "fmt"

var nameficator groupSequentialNamer

type groupSequentialNamer struct {
	counter int
}

func (r *groupSequentialNamer) getNextGroupName() string {
	r.counter++
	return fmt.Sprintf("group_%d", r.counter)
}
