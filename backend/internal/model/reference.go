package model

type Status struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type Company struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WorkPosition struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type University struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type Faq struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type EnumValue struct {
	Value string `json:"value"`
}
