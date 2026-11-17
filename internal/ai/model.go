package ai

type Input struct {
	Script string
	Output string
}

type Candidate struct {
	Script string
	Weight int
}

type Result struct {
	Candidates []Candidate `json:"candidates"`
}
