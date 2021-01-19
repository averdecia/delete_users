package main

// LookupResponse the answer from middleware
type LookupResponse struct {
	Status string `json:"status,omitempty"`
}

// MWresponse the answer from middleware
type MWresponse struct {
	Data Data `json:"data,omitempty"`
}

// Data ...
type Data struct {
	Response Response `json:"response,omitempty"`
}

// Response ...
type Response struct {
	Exist      bool     `json:"emailExist,omitempty"`
	SocialData []Social `json:"socialdata,omitempty"`
}

// Social ...
type Social struct {
	IDUser         string `json:"ID_USUARIO,omitempty"`
	GamificationID string `json:"ID_USUARIO_SOCIAL,omitempty"`
	RedSocial      string `json:"REDSOCIAL,omitempty"`
}

// Args is saving user input parameters
type Args struct {
	Endpoint   string // the endpoint for lookup server
	InputPath  string // the path to read the element rows
	OutputPath string // the path to save the element errors rows
	GoRoutines int    // the amount of concurrent process
	AuthToken  string // Basic token for server
	Middleware string // the endpoint for middleware server
}
