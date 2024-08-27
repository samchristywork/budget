package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var saveFile = "inducted.csv"

type LineItem struct {
	Category string
	Date     string
	Amount   string
	Data     string
}

type HistogramRow struct {
	Amount float64
	Number int
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

func addNewLineItems(lineItems []LineItem, file string) []LineItem {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer f.Close()

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

		if strings.Contains(line, "	") {
			line = strings.ReplaceAll(line, "	", " ")
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

func deserialize(saveFile string) []LineItem {
	f, err := os.Open(saveFile)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return []LineItem{}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineItems := []LineItem{}
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "	")
		lineItem := LineItem{}
		lineItem.Category = parts[0]
		lineItem.Date = parts[1]
		lineItem.Amount = parts[2]
		lineItem.Data = parts[3]
		lineItems = append(lineItems, lineItem)
	}
	return lineItems
}

func induct(file string) {
	lineItems := deserialize(saveFile)
	lineItems = addNewLineItems(lineItems, file)
	serialize(lineItems)
}

func printSummary(startDate string, endDate string) {
	fmt.Printf("Summary from %s to %s\n", startDate, endDate)

	lineItems := deserialize(saveFile)

	categoryMap := make(map[string]HistogramRow)
	for _, lineItem := range lineItems {
		if lineItem.Date < startDate || lineItem.Date > endDate {
			continue
		}

		amount, _ := strconv.ParseFloat(lineItem.Amount, 64)
		categoryMap[lineItem.Category] = HistogramRow{
			Amount: categoryMap[lineItem.Category].Amount + amount,
			Number: categoryMap[lineItem.Category].Number + 1,
		}
	}

	for category, amount := range categoryMap {
		if category == "Ignored" {
			continue
		}

		fmt.Printf("%s: %.2f (%d Transactions)\n",
			category, amount.Amount, amount.Number)
	}
}

func main() {
	end := time.Now().Format("2006-01-02")
	start := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	inductFile := flag.String("induct", "", "The file to induct")
	summaryFlag := flag.Bool("summary", false, "Print a summary of all transactions")
	endDate := flag.String("end", end, "The end date for the summary")
	startDate := flag.String("start", start, "The start date for the summary")

	flag.Parse()

	if *inductFile != "" {
		induct(*inductFile)
	} else if *summaryFlag {
		printSummary(*startDate, *endDate)
	} else {
		flag.Usage()
	}
}
