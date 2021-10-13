package main

type Grok struct {
	Field              string            `json:"field"`
	PatternDefinitions map[string]string `json:"pattern_definitions"`
	Patterns           []string          `json:"patterns"`
	IgnoreMissing      bool              `json:"ignore_missing"`
}

type Processor struct {
	Grok Grok `json:"grok"`
}

type Set struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type OnFailure struct {
	Set Set `json:"set"`
}

type Pipeline struct {
	Description string      `json:"description"`
	Processors  []Processor `json:"processors"`
	OnFailure   []OnFailure `json:"on_failure"`
}

type Doc struct {
	Message string `json:"message"`
}

type SimulateDoc struct {
	Index  string            `json:"_index"`
	ID     string            `json:"_id"`
	Source map[string]interface{} `json:"_source"`
}