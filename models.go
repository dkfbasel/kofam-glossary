package main

type GlossaryItem struct {
	Url         string
	Name        string
	English     string
	Description string
	Source      string
	ContentFull string
}

// define a glossary page (for each letter in the alphabet)
type GlossaryPage struct {
	Letter string
	Url    string
}
