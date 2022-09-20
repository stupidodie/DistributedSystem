package main

import (
	"fmt"
	"sync"
)

type Command struct {
	Type string
	recv chan bool
}
type philosopher struct {
	id        int
	left      chan Command
	right     chan Command
	max_count int
	count     int
}
type fork struct {
	id          int
	count       int
	isAvailable bool
	request     chan Command
	max_count   int
}

func (f *fork) run() {

out:
	for {
		cmd := <-f.request
		switch cmd.Type {
		case "set":
			f.count++
			f.isAvailable = true
			if f.count == f.max_count {
				// fmt.Printf("%d is exit \n", f.id)
				break out
			}
		case "get":
			if f.isAvailable {
				f.isAvailable = false
				cmd.recv <- true
			} else {
				cmd.recv <- false
			}
		case "clear":
			// fmt.Printf("%d is clear \n", f.id)
			f.isAvailable = true
		default:
			panic("unknown error")
		}
	}
}
func (p *philosopher) think() {
	fmt.Printf("%d is thinking\n", p.id)
}
func (p *philosopher) eat() {
	recv := make(chan bool)
	defer close(recv)
	for {
		cmd_get := Command{
			Type: "get",
			recv: recv,
		}
		p.left <- cmd_get
		left_fork := <-recv
		if left_fork {
			fmt.Printf("%d is get %d\n", p.id, p.id%count)
			p.right <- cmd_get
			right_fork := <-recv
			if right_fork {
				fmt.Printf("%d is get %d\n", p.id, (p.id+1)%count)
				cmd_set := Command{
					Type: "set",
					recv: nil,
				}
				p.left <- cmd_set
				p.right <- cmd_set
				p.count++
				fmt.Printf("%d is eating \n", p.id)
				if p.count == p.max_count {
					fmt.Printf("%d is exit \n", p.id)
					break
				}
			} else {
				cmd_clear := Command{
					Type: "clear",
					recv: nil,
				}
				p.left <- cmd_clear
			}

		}
	}
}

const count = 5
const philosopher_number = count
const fork_number = count
const max_number = 3
const fork_max_count = 2 * max_number
const philosopher_max_count = max_number

func main() {
	var requests [fork_number]chan Command
	for i := range requests {
		requests[i] = make(chan Command)
	}

	forks := make([]*fork, fork_number)
	for i := 0; i < fork_number; i++ {
		forks[i] = &fork{id: i, request: requests[i], max_count: fork_max_count, count: 0, isAvailable: true}
	}

	philosophers := make([]*philosopher, philosopher_number)
	for i := 0; i < philosopher_number; i++ {
		philosophers[i] = &philosopher{
			id:        i,
			left:      requests[i%fork_number],
			right:     requests[(i+1)%fork_number],
			max_count: philosopher_max_count,
			count:     0,
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < fork_number; i++ {
		j := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			forks[j].run()
		}()
	}
	for i := 0; i < philosopher_number; i++ {
		j := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			philosophers[j].think()
			philosophers[j].eat()
		}()
	}
	wg.Wait()
	for i := range requests {
		close(requests[i])
	}
}
