package redmine

type NamedItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type IssueReference struct {
	Id int `json:"id"`
}

type CustomField struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
