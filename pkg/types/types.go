package types

type ProcessedText struct {
	BaseLanguage   string      `json:"base_language"`
	TargetLanguage string      `json:"target_language"`
	Paragraphs     []Paragraph `json:"paragraphs"`
}

type Paragraph []Chunk

type Chunk struct {
	Tokens      []Token `json:"tokens"`
	Translation string  `json:"translation"`
}

type Token struct {
	Text  string `json:"text"`
	Lemma string `json:"lemma"`
}
