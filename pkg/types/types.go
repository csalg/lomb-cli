package types

type ProcessedText struct {
	BaseLanguage   string      `json:"base_language"`
	TargetLanguage string      `json:"target_language"`
	Paragraphs     []Paragraph `json:"paragraphs"`
}

func (txt ProcessedText) Translation() string {
	sentence := ""
	for _, paragraph := range txt.Paragraphs {
		sentence += paragraph.Translation() + "\n"
	}
	return sentence
}

type Paragraph []Chunk

func (p Paragraph) Translation() string {
	sentence := ""
	for _, chunk := range p {
		sentence += chunk.Translation + " "
	}
	return sentence
}

type Chunk struct {
	Tokens      []Token `json:"tokens"`
	Text        string  `json:"text"`
	Translation string  `json:"translation"`
}

func (ch Chunk) Sentence() string {
	if ch.Text != "" {
		return ch.Text
	}
	sentence := ""
	for _, token := range ch.Tokens {
		sentence += token.Text + " "
	}
	return sentence
}

type Token struct {
	Text  string `json:"text"`
	Lemma string `json:"lemma"`
}
