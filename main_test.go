package main

import "testing"

const buffer = `{
		"let": {
			"prefix": "let",
			"body": [
				"let ${1:varname} : ${2:type} = $0"
			],
			"description": "unmutable var declaration"
		},
		"let_mut": {
			"prefix": "let_mut",
			"body": [
				"let mut ${1:varname} : ${2:type} = $0"
			],
			"description": "mutable var declaration"
		}
	}
	`

const expectedTable = `| prefix | description |
| :----- | :---------- |
| let | unmutable var declaration |
| let_mut | mutable var declaration |
`

func TestSnippetsParsing(t *testing.T) {
	items := []TableItem{
		{"let", "unmutable var declaration"},
		{"let_mut", "mutable var declaration"},
	}

	parsedItems, err := snippetsParsing([]byte(buffer))
	if err != nil {
		t.Errorf("error parsing %q", err.Error())
	}

	piLen := len(parsedItems)
	if len(parsedItems) != 2 {
		t.Errorf("wrong number of items parsed, got %d, wanted 2", piLen)	
	}

	for i := range parsedItems {
		if parsedItems[i].prefix != items[i].prefix || parsedItems[i].description != items[i].description {
			t.Errorf("Wrong item! got%v, wanted %v", parsedItems[i], items[i])
		}
	}
}

func TestBuildMDTable(t *testing.T) {
	items := []TableItem{
		{"let", "unmutable var declaration"},
		{"let_mut", "mutable var declaration"},
	}
	mdTable, err := buildMDTable(items)
	if err != nil {
		t.Errorf("error building MD table: %q", err.Error())
	}

	if mdTable != expectedTable {
		t.Errorf("MD Tables mismatch! Given:\n%q\nWanted:\n%q", mdTable, expectedTable)
	}
}