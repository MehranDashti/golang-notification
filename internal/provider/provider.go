package provider

type Message struct {
	To       string
	Title    string
	Body     string
	Metadata map[string]string
}

type Result struct {
	ProviderID string
	Provider   string
}
