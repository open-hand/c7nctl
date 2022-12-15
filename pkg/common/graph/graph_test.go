package graph

import (
	"github.com/choerodon/c7nctl/pkg/resource"
	"testing"
)

/*
2 -> 1 -> 3 -> 4 -> 5
*/
func buildGraphA() (g Graph) {
	rls1 := &resource.Release{Name: "1"}
	rls2 := &resource.Release{Name: "2"}
	rls3 := &resource.Release{Name: "3"}
	rls4 := &resource.Release{Name: "4"}
	rls5 := &resource.Release{Name: "5"}

	g.AddVertex(rls1)
	g.AddVertex(rls2)
	g.AddVertex(rls3)
	g.AddVertex(rls4)
	g.AddVertex(rls5)

	g.AddEdges(rls2, rls1)
	g.AddEdges(rls1, rls3)
	g.AddEdges(rls2, rls3)
	g.AddEdges(rls3, rls4)
	g.AddEdges(rls3, rls5)
	g.AddEdges(rls5, rls4)
	return g
}

/*
1 -> 2 -> 3 -> 4
*/
func buildGraphB() (g Graph) {
	rls1 := &resource.Release{Name: "1"}
	rls2 := &resource.Release{Name: "2"}

	rls3 := &resource.Release{Name: "3"}
	rls4 := &resource.Release{Name: "4"}

	g.AddVertex(rls1)
	g.AddVertex(rls2)
	g.AddVertex(rls3)
	g.AddVertex(rls4)

	g.AddEdges(rls1, rls2)
	g.AddEdges(rls1, rls4)
	g.AddEdges(rls2, rls3)
	g.AddEdges(rls3, rls4)
	return g
}

func TestGraph_TopoSortByKahn(t *testing.T) {

	GraphTest := []struct {
		Graph
		result []string
	}{
		{
			buildGraphA(),
			[]string{"2", "1", "3", "5", "4"},
		},
		{
			buildGraphB(),
			[]string{"1", "2", "3", "4"},
		},
	}

	for _, tt := range GraphTest {
		q := tt.TopoSortByKahn()
		for i := 0; i < len(tt.result); i++ {
			rls := q.Dequeue()

			if rls != nil && tt.result[i] != rls.Name {
				t.Errorf("Graph error sorting: release %s", rls.Name)
			}
		}
	}
}

func TestNewReleaseGraph(t *testing.T) {
	rls1 := resource.Release{
		Name:         "1",
		Requirements: nil,
	}
	rls2 := resource.Release{
		Name:         "2",
		Requirements: []string{"3"},
	}
	rls3 := resource.Release{
		Name:         "3",
		Requirements: []string{"2"},
	}
	rls4 := resource.Release{
		Name:         "4",
		Requirements: []string{"1", "2"},
	}
	rls5 := resource.Release{
		Name:         "5",
		Requirements: []string{"3"},
	}

	q1 := []*resource.Release{&rls1, &rls2, &rls3}
	q2 := []*resource.Release{&rls1, &rls2, &rls3, &rls4, &rls5}
	ReleaseGraphTest := []struct {
		*Graph
		result []string
	}{
		{
			NewReleaseGraph(q1),
			[]string{"1", "3", "2"},
		},
		{
			NewReleaseGraph(q2),
			[]string{"1", "3", "5", "2", "4"},
		},
	}

	for _, tt := range ReleaseGraphTest {
		q := tt.TopoSortByKahn()
		for i := 0; i < len(tt.result); i++ {
			rls := q.Dequeue()

			if rls != nil && tt.result[i] != rls.Name {
				t.Errorf("Graph error sorting: release %s", rls.Name)
			}
		}
	}
}
