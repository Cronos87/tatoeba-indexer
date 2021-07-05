package main

// Sentence describes the fields to index.
type Sentence struct {
	ID                  int32           `json:"id"`
	Language            string          `json:"language"`
	Content             string          `json:"content"`
	Username            string          `json:"username"`
	AddedAt             string          `json:"added_at"`
	UpdatedAt           string          `json:"updated_at"`
	DirectRelations     []int32         `json:"direct_translations"`
	IndirectRelations   []int32         `json:"indirect_translations"`
	TranslatedLanguages []string        `json:"translated_languages"`
	AudioUsername       string          `json:"audio_username,omitempty"`
	Transcriptions      []Transcription `json:"transcriptions,omitempty"`
}

// Transcription describes the fields used to simplify
// reading languages like Chinese, Cantonese or Japanese.
type Transcription struct {
	ScriptName    string `json:"script_name"`
	Username      string `json:"username"`
	Transcription string `json:"transcription"`
}

// Indexer define the methods indexers need to implement.
type Indexer interface {
	Init()
	Index(map[string]Sentence)
}
