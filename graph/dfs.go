package graph

import (
	"container/list"
)

type DFSElement struct {
	Color int
	D, F  int
	P     *DFSElement
	V     interface{}
}

func (e *DFSElement) Init(v interface{}) *DFSElement {
	e.V = v
	e.Color = WHITE
	e.D = 0
	e.F = 0
	e.P = nil
	return e
}

//path compression, inspiration from disjoint-set
func (e *DFSElement) FindRoot() *DFSElement {
	_e := e
	for _e.P != nil {
		if _e.P.P != nil {
			_e.P = _e.P.P
		}
		_e = _e.P
	}
	return _e

}

func NewDFSElement(v interface{}) *DFSElement {
	return new(DFSElement).Init(v)
}

func DFS(g Graph, sorter func([]interface{})) (dfsGraph map[string]Graph) {
	dfsGraph = make(map[string]Graph)
	dfsGraph["dfsForest"] = CreateGraphByType(g)
	dfsGraph["dfsBackEdges"] = CreateGraphByType(g)
	dfsGraph["dfsForwardEdges"] = CreateGraphByType(g)
	dfsGraph["dfsCrossEdges"] = CreateGraphByType(g)

	timer := 0
	stack := list.New()
	//to keep vertices order
	elements := new(linkedMap).init()
	//init
	vertices := g.AllVertices()
	if sorter != nil {
		sorter(vertices)
	}
	for _, v := range vertices {
		elements.add(v, NewDFSElement(v))
	}
	for v := elements.frontKey(); v != nil; v = elements.nextKey(v) {
		if elements.get(v).(*DFSElement).Color == WHITE {
			//push root vertex to stack
			stack.PushBack(elements.get(v).(*DFSElement))
			for stack.Len() != 0 {
				e := stack.Back().Value.(*DFSElement);
				if e.Color == BLACK {
					//if is black, it must be deeper in stack, and has be visited through deeper path
					stack.Remove(stack.Back())
				} else {
					//white or gray
					if e.Color == WHITE {
						e.Color = GRAY
						timer ++
						e.D = timer
						//add all children has not been visited to stack
						for _, c := range g.AllConnectedVertices(e.V) {
							if elements.get(c).(*DFSElement).Color == WHITE {
								// parent in deeper path always override that in shallower
								if elements.get(c).(*DFSElement).P != nil {
									//if it is shallower in the dfs tree and it is white, it's parent will be override,and it's a forward edge
									dfsGraph["dfsForwardEdges"].AddEdge(Edge{elements.get(c).(*DFSElement).P, elements.get(c).(*DFSElement)})
								}
								elements.get(c).(*DFSElement).P = e
								stack.PushBack(elements.get(c).(*DFSElement))
							} else if elements.get(c).(*DFSElement).Color == GRAY {
								// if color is already gray, it's a flip edge
								dfsGraph["dfsBackEdges"].AddEdge(Edge{e,elements.get(c).(*DFSElement)})
							} else {
								// if color is already black, it's a cross edge,d(e) > d(elements.get(c).(*DFSElement))
								dfsGraph["dfsCrossEdges"].AddEdge(Edge{e,elements.get(c).(*DFSElement)})
							}
						}
					}
					if e == stack.Back().Value.(*DFSElement) {
						// if the stack did not grow, it is end-point vertex, finish visit process and pop stack
						e.Color = BLACK
						timer ++
						e.F = timer
						stack.Remove(stack.Back())
						for i := range dfsGraph {
							dfsGraph[i].AddVertex(e)
						}
						//tree edge definition. First time visit
						if e.P != nil {
							dfsGraph["dfsForest"].AddEdge(Edge{e.P, e})
						}
					}
					// else if the stack grew, update pointer to the top of stack and visit it
				}
			}
		}
	}
	return
}