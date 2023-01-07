package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TableItem struct {
	prefix string
	description string
}

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
					fmt.Fprintf(os.Stderr, "Prefix `%v` is not a valid string", v)
					os.Exit(5)
				}
				if k == "prefix" {
					item.prefix = parsedVal
				} else {
					item.description = parsedVal
				}
			case "body": continue
			default:
				fmt.Fprintf(os.Stderr, "`%s` is not a valid parameter for snippet", k)
				os.Exit(6)
			}
		}
		prefixDesc = append(prefixDesc, item)
	}
	
	return prefixDesc, nil
}

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

func writeMDFile(mdTable string) (err error) {
	fd, err := os.OpenFile("snippets_table.md", os.O_WRONLY, 0666)
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
	
	if !strings.HasSuffix(os.Args[1], ".code-snippets") && !strings.HasSuffix(os.Args[1], ".json") {
		fmt.Fprintf(os.Stderr, "Wrong file extension: Given `%s`. Accepted: `code-snippets`\n", os.Args[1])
		os.Exit(1)
	}

	buffer, err := os.ReadFile(os.Args[1])
	checkError(err)

	prefixDesc, err := snippetsParsing(buffer)
	checkError(err)

	markdownTable, err := buildMDTable(prefixDesc)
	checkError(err)

	err = writeMDFile(markdownTable)
	checkError(err)
}