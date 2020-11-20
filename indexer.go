package main

// Indexer define the methods indexers need to implement.
type Indexer interface {
	Init()
	Index(map[int32]Sentence)
}
