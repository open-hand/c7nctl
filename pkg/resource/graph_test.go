package resource

import "testing"

var g Graph

func initGraph() {

	rls1 := &Release{Name: "1"}
	rls2 := &Release{Name: "2"}

	rls3 := &Release{Name: "3"}
	rls4 := &Release{Name: "4"}

	g.AddVertex(rls1)
	g.AddVertex(rls2)
	g.AddVertex(rls3)
	g.AddVertex(rls4)

	g.AddEdges(rls1, rls2)
	g.AddEdges(rls1, rls4)
	g.AddEdges(rls2, rls3)
	g.AddEdges(rls3, rls4)
	g.AddVertex(rls2)
}
func TestGraph_TopoSortByKahn(t *testing.T) {
	initGraph()
	queues := g.TopoSortByKahn()

	for queues.IsEmpty() {
		rls := queues.Dequeue()
		t.Log(rls.Name)
	}
}
