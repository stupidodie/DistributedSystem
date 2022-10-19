package main

import "fmt"

type Vector_clock struct {
	vector_size  int
	vector_clock []byte
	clock_id     int
}

func (local_clock Vector_clock) update(coming_clock Vector_clock) Vector_clock {
	for index, clock := range coming_clock.vector_clock[0:coming_clock.vector_size] {
		fmt.Println(index, clock)
		if clock > local_clock.vector_clock[index] {
			local_clock.vector_clock[index] = clock
		}
	}
	local_clock.vector_clock[local_clock.clock_id]++
	return local_clock
}

func main() {
	var local_vector_clock []byte
	local_vector_clock = make([]byte, 10)
	fmt.Println(local_vector_clock)

	var v = Vector_clock{vector_size: 3, vector_clock: []byte{1, 1, 2}, clock_id: 2}
	var p = Vector_clock{vector_size: 3, vector_clock: []byte{0, 3, 0}, clock_id: 1}
	// p.increase(v)
	fmt.Println(v.update(p))

}
