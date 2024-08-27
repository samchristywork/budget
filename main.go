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

func main() {
}
