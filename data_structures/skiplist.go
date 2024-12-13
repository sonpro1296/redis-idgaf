package data_structures

// Skiplist is used in redis sorted set structure, as in ZSET, ZADD, ZREM,...

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

type Node struct {
	key     string
	score   float64
	forward []*Node
}

type Skiplist struct {
	level   int
	p       float64
	head    *Node
	nodeMap map[string]*Node
}

func NewNode(key string, score float64, level int) *Node {
	return &Node{
		key:     key,
		score:   score,
		forward: make([]*Node, level),
	}
}

func NewSkiplist(level int, p float64) *Skiplist {
	headNewNode := make([]*Node, level)
	newNode := NewNode("None key", math.MaxFloat64*(-1), level)
	for i := 0; i < level; i++ {
		headNewNode[i] = newNode
	}
	node := &Node{"", 0, headNewNode}

	return &Skiplist{
		level:   level,
		p:       p,
		head:    node,
		nodeMap: make(map[string]*Node),
	}
}

func (s *Skiplist) Display() {
	for i := 0; i < s.level; i++ {
		fmt.Printf("Level %d: head ", i)
		node := s.head.forward[i]
		for (node.forward[i] != nil) && (node.forward[i].key != "None key") {
			fmt.Printf(" -> %s: %f ", node.forward[i].key, node.forward[i].score)
			node = node.forward[i]
		}
		fmt.Printf(" -> nil\n")

	}
	// fmt.Println(s.head.forward[0].forward[0])
	// fmt.Print(len(s.nodeMap))
}

func (s *Skiplist) Add(key string, score float64) error {
	if _, ok := s.nodeMap[key]; ok {
		log.Print("key exists")
		return nil
	}

	update := make([]*Node, s.level)

	// select level
	level := 0
	for i := 0; i < s.level; i++ {
		if r := rand.Float64(); r < s.p {
			level++
		} else {
			break
		}
	}
	startNode := s.head.forward[level]
	i := level

	for i >= 0 {
		for (startNode != nil) && (startNode.forward[i] != nil) && (startNode.forward[i].score < score) {
			startNode = startNode.forward[i]
		}
		update[i] = startNode
		i--

	}

	newNode := NewNode(key, score, s.level)
	for i := 0; i <= level; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	s.nodeMap[key] = newNode

	return nil

}

func (s *Skiplist) Delete(key string) error {
	var deleteNode *Node
	var ok bool

	if deleteNode, ok = s.nodeMap[key]; !ok {
		return fmt.Errorf("key not exist")
	}

	update := make([]*Node, s.level)

	head := s.head.forward[s.level-1]

	for i := s.level - 1; i >= 0; i-- {
		for (head.forward[i] != nil) && (head.forward[i] != deleteNode) && (head.forward[i].score <= deleteNode.score) {
			head = head.forward[i]
		}
		if head.forward[i] == deleteNode {
			update[i] = head
		}
	}
	for i, node := range update {
		if node != nil {
			node.forward[i] = deleteNode.forward[i]
		}
	}

	delete(s.nodeMap, key)
	deleteNode = nil

	// fmt.Println(update)

	return nil

}

func (s *Skiplist) Search(key string) (*Node, error){
	if node, ok := s.nodeMap[key]; ok {
		return node, nil
	} 
	return nil, fmt.Errorf("key not exist")
}
