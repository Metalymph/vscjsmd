package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type TableItem struct {
	prefix string
	description string
}

//snippetsParsing parses the buffer from read file and returns a slice of (prefix, description)
func snippetsParsing(buffer []byte) ([]TableItem, error) {
	var decSnippets map[string]map[string]interface{}
	if err := json.Unmarshal(buffer, &decSnippets); err != nil {
		return nil, err
	}

	var prefixDesc []TableItem
	for _, snippet := range decSnippets {
		var item TableItem
		for k, v := range snippet {
			switch k {
			case "prefix", "description":
				parsedVal, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("prefix `%v` is not a valid string", v)
				}
				if k == "prefix" {
					item.prefix = parsedVal
				} else {
					item.description = parsedVal
				}
			case "body": continue
			default:
				return nil, fmt.Errorf("`%s` is not a valid parameter for snippet", k)
			}
		}
		prefixDesc = append(prefixDesc, item)
	}
	
	return prefixDesc, nil
}

//buildMDTable format the full snippets info table into a string
func buildMDTable(prefixDesc []TableItem) (mdTable string, err error) {
	var sb strings.Builder
	_, err = sb.WriteString("| prefix | description |\n| :----- | :---------- |\n")
	if err != nil {
		return "", err
	}

	for _, tableItem := range prefixDesc {
		str := fmt.Sprintf("| %s | %s |\n", tableItem.prefix, tableItem.description)
		_, err = sb.WriteString(str)
		if err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}

//writeMDFile creates and writes the full Markdown table in a temporary file
func writeMDFile(mdTable string) (err error) {
	fd, err := os.OpenFile("snippets_table.md", os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.WriteString(mdTable)
	if err != nil {
		return err
	}

	fmt.Println("Markdown `snippets_table.md` correctly saved.")
	return nil
}

//checkError reduces verbosity for error checking
func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main()  {
	argcCount := len(os.Args)
	if argcCount != 2 {
		fmt.Fprintf(os.Stderr, "Wrong number of arguments: Given %d, wanted 1 instead\n", argcCount - 1)
		os.Exit(1)
	}
	
	if !strings.HasSuffix(os.Args[1], ".code-snippets") {
		fmt.Fprintf(os.Stderr, "Wrong file extension: Given `%s`. Accepted: `code-snippets`\n", os.Args[1])
		os.Exit(1)
	}

	buffer, err := os.ReadFile(os.Args[1])
	checkError(err)

	prefixDesc, err := snippetsParsing(buffer)
	checkError(err)

	//sorting by asc. prefix name
	sort.SliceStable(prefixDesc, func(i, j int) bool {
		return prefixDesc[i].prefix < prefixDesc[j].prefix
	})

	markdownTable, err := buildMDTable(prefixDesc)
	checkError(err)

	err = writeMDFile(markdownTable)
	checkError(err)
}