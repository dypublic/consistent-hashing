package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"hash/crc32"
)

type hashNode struct {
	code uint32
	node string
}

type HashingRing struct {
	hashVnodes map[uint32]string
	vnode      map[string]string
	sortedNode []*hashNode
}

func NewHashingRing() *HashingRing {
	h := &HashingRing{}
	h.hashVnodes = make(map[uint32]string)
	h.vnode = make(map[string]string)
	return h
}

func (h *HashingRing) resortVnode() {
	var vnodes []*hashNode
	for hash, vnode := range h.hashVnodes {
		vnodes = append(vnodes, &hashNode{code: hash, node: vnode})
	}
	sort.Slice(vnodes, func(i, j int) bool {
		return vnodes[i].code < vnodes[j].code
	})
	h.sortedNode = vnodes
}

func (h *HashingRing) Put(node string, cap int) error {
	if _, ok := h.vnode[node]; ok == true {
		return errors.New("already have")
	}
	var vnames []string
	for i:=0; i <= cap; i++ {
		vname := node+strconv.FormatUint(uint64(i),10)
		vnames = append(vnames, vname)
		h.vnode[vname] = node
		hash_num := crc32.ChecksumIEEE([]byte(vname))
		h.hashVnodes[hash_num] = vname
	}
	h.resortVnode()
	return nil

}

func (h *HashingRing) String() string {
	res := ""
	for _, n := range h.sortedNode{
		res = res + fmt.Sprintf("(%x, %s) ", n.code, n.node)
	}
	return res
}

func search(hash uint32, list []*hashNode) *hashNode {
	start_pos := 0
	mid_pos := len(list) /2
	end_pos := len(list)
	if hash > list[end_pos].code || hash < list[start_pos].code {
		return list[start_pos]
	}
	for {
		if hash == list[mid_pos].code {
			return list[mid_pos]
		}else if hash > list[mid_pos].code{
			start_pos = mid_pos
			mid_pos = (end_pos - start_pos)/2 + start_pos
		} else {
			end_pos = mid_pos
			mid_pos = (end_pos - start_pos)/2 + start_pos
		}
		if start_pos +1 == end_pos || start_pos == end_pos {
			return list[end_pos]
		}

	}
}

func (h HashingRing) target(name string) string {
	hash_num := crc32.ChecksumIEEE([]byte(name))
	node := search(hash_num, h.sortedNode)
	return node.node
}

func main() {
	h := NewHashingRing()
	nodes := []string{"apple", "orange", "banana"}
	for _, n := range nodes {
		h.Put(n, 5)
	}
	node := h.target("pear")
	fmt.Println(node)
	fmt.Println(h)
}