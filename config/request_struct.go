package config

type Request struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Params  map[string]string `json:"params,omitempty"`
	Auth    Auth              `json:"auth"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    Body              `json:"body"`
}

type Auth struct {
	Bearer string `json:"bearer,omitempty"`
	Basic  string `json:"basic,omitempty"`
}

type Body struct {
	StringJson string `json:"json,omitempty"`
	Text       string `json:"text,omitempty"`
	Raw        string `json:"raw,omitempty"`
}
