package main

import (
	"container/heap"
	"fmt"
)

type Queue struct {
	items []Vertex
	m     map[Vertex]int
	pr    map[Vertex]int
}

type sg struct {
	ids   map[string]Vertex
	names map[Vertex]string
	edges map[Vertex]map[Vertex]int
}

type Graph interface {
	Vertices() []Vertex
	Neighbors(v Vertex) []Vertex
	Weight(u, v Vertex) int
}

type Vertex int

func (q *Queue) Len() int           { return len(q.items) }
func (q *Queue) Less(i, j int) bool { return q.pr[q.items[i]] < q.pr[q.items[j]] }
func (q *Queue) Swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
	q.m[q.items[i]] = i
	q.m[q.items[j]] = j
}
func (q *Queue) Push(x interface{}) {
	n := len(q.items)
	item := x.(Vertex)
	q.m[item] = n
	q.items = append(q.items, item)
}
func (q *Queue) Pop() interface{} {
	old := q.items
	n := len(old)
	item := old[n-1]
	q.m[item] = -1
	q.items = old[0 : n-1]
	return item
}

func (q *Queue) update(item Vertex, priority int) {
	q.pr[item] = priority
	heap.Fix(q, q.m[item])
}
func (q *Queue) addWithPriority(item Vertex, priority int) {
	heap.Push(q, item)
	q.update(item, priority)
}

const (
	Infinity      = int(^uint(0) >> 1)
	Uninitialized = -1
)

func ShortestPath(g Graph, source Vertex) (dist map[Vertex]int, prev map[Vertex]Vertex) {
	dist = make(map[Vertex]int)
	prev = make(map[Vertex]Vertex)
	sid := source
	dist[sid] = 0
	q := &Queue{[]Vertex{}, make(map[Vertex]int), make(map[Vertex]int)}
	for _, v := range g.Vertices() {
		if v != sid {
			dist[v] = Infinity
		}
		prev[v] = Uninitialized
		q.addWithPriority(v, dist[v])
	}
	for len(q.items) != 0 {
		u := heap.Pop(q).(Vertex)
		for _, v := range g.Neighbors(u) {
			alt := dist[u] + g.Weight(u, v)
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
				q.update(v, alt)
			}
		}
	}
	return dist, prev
}

func newgraph(ids map[string]Vertex) sg {
	g := sg{ids: ids}
	g.names = make(map[Vertex]string)
	for k, v := range ids {
		g.names[v] = k
	}
	g.edges = make(map[Vertex]map[Vertex]int)
	return g
}
func (g sg) edge(u, v string, w int) {
	if _, ok := g.edges[g.ids[u]]; !ok {
		g.edges[g.ids[u]] = make(map[Vertex]int)
	}
	g.edges[g.ids[u]][g.ids[v]] = w
}
func (g sg) path(v Vertex, prev map[Vertex]Vertex) (s string) {
	s = g.names[v]
	for prev[v] >= 0 {
		v = prev[v]
		s = g.names[v] + "->" + s
	}
	return s
}
func (g sg) Vertices() (vs []Vertex) {
	for _, v := range g.ids {
		vs = append(vs, v)
	}
	return vs
}
func (g sg) Neighbors(u Vertex) (vs []Vertex) {
	for v := range g.edges[u] {
		vs = append(vs, v)
	}
	return vs
}
func (g sg) Weight(u, v Vertex) int { return g.edges[u][v] }

func main() {
	graph := newgraph(map[string]Vertex{
		"aarau":    1,
		"basel":    2,
		"chur":     3,
		"dietikon": 4,
		"zürich":   5,
		"freiburg": 6,
	})
	graph.edge("aarau", "berlin", 7)
	graph.edge("aarau", "chur", 9)
	graph.edge("aarau", "freiburg", 14)
	graph.edge("basel", "chur", 10)
	graph.edge("basel", "dietikon", 15)
	graph.edge("chur", "dietikon", 11)
	graph.edge("chur", "freiburg", 2)
	graph.edge("basel", "zürich", 6)
	graph.edge("zürich", "freiburg", 9)

	dist, prev := ShortestPath(graph, graph.ids["aarau"])
	fmt.Printf("Distance to %s: %d, Path: %s\n", "zürich", dist[graph.ids["zürich"]], graph.path(graph.ids["zürich"], prev))
	fmt.Printf("Distance to %s: %d, Path: %s\n", "freiburg", dist[graph.ids["freiburg"]], graph.path(graph.ids["freiburg"], prev))
}
