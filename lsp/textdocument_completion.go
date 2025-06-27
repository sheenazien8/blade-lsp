package lsp

type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}

type CompletionParams struct {
	TextDocumentPositionParams
	Context CompletionContext `json:"context,omitempty"`
}

type CompletionContext struct {
	TriggerKind      int    `json:"triggerKind"`
	TriggerCharacter string `json:"triggerCharacter,omitempty"`
}

type CompletionResponse struct {
	Response
	Result []CompletionItem `json:"result"`
}

type CompletionItem struct {
	Label            string `json:"label"`
	Kind             int    `json:"kind,omitempty"`
	Detail           string `json:"detail"`
	Documentation    string `json:"documentation"`
	InsertText       string `json:"insertText,omitempty"`
	InsertTextFormat int    `json:"insertTextFormat,omitempty"`
}
