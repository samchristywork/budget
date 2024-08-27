package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var saveFile = "inducted.csv"

type LineItem struct {
	Category string
	Date     string
	Amount   string
	Data     string
}

func serialize(lineItems []LineItem) {
	f, err := os.Create(saveFile)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}
	defer f.Close()

	for _, lineItem := range lineItems {
		f.WriteString(lineItem.Category + "	")
		f.WriteString(lineItem.Date + "	")
		f.WriteString(lineItem.Amount + "	")
		f.WriteString(lineItem.Data + "\n")
	}
}

func getInput(prompt string) string {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func addNewLineItems(lineItems []LineItem, f *os.File) []LineItem {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		exists := false
		for _, lineItem := range lineItems {
			if strings.Contains(line, lineItem.Data) {
				exists = true
				continue
			}
		}
		if exists {
			continue
		}

		fmt.Println(line)
		lineItem := LineItem{}
		lineItem.Category = getInput("Category: ")
		if lineItem.Category == "" {
			lineItem.Category = "Ignored"
			lineItem.Data = line
			lineItems = append(lineItems, lineItem)
		} else {
			lineItem.Date = getInput("Date: ")
			lineItem.Amount = getInput("Amount: ")
			lineItem.Data = line
			lineItems = append(lineItems, lineItem)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}

	return lineItems
}

func main() {
}
