package main

// Indexer define the methods indexers need to implement.
type Indexer interface {
	Init()
	Index(map[string]Sentence)
}
