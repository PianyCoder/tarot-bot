package genapi_cli

type GenAPIRequest struct {
	Messages         []GenAPIMessage `json:"messages"`
	IsSync           bool            `json:"is_sync"`
	Stream           bool            `json:"stream"`
	N                int             `json:"n"`
	FrequencyPenalty float64         `json:"frequency_penalty"`
	MaxTokens        int             `json:"max_tokens"`
	PresencePenalty  float64         `json:"presence_penalty"`
	Temperature      float64         `json:"temperature"`
	TopP             float64         `json:"top_p"`
	ResponseFormat   string          `json:"response_format"`
}

type GenAPIMessage struct {
	Role    string          `json:"role"`
	Content []GenAPIContent `json:"content"`
}

type GenAPIContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type GenAPIResponse struct {
	RequestID int            `json:"request_id"`
	Model     string         `json:"model"`
	Cost      float64        `json:"cost"`
	Response  []GenAPIOutput `json:"response"`
}

type GenAPIOutput struct {
	Index   int `json:"index"`
	Message struct {
		Role    string      `json:"role"`
		Content string      `json:"content"`
		Refusal interface{} `json:"refusal"`
	} `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}
