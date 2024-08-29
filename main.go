package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var saveFile = "inducted.csv"

type LineItem struct {
	Category    string
	Date        string
	Amount      string
	Description string
}

type HistogramRow struct {
	Amount float64
	Number int
}

func serialize(lineItems []LineItem) {
	f, err := os.Create(saveFile)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		os.Exit(1)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	for _, lineItem := range lineItems {
		writer.Write([]string{
			lineItem.Category,
			lineItem.Date,
			lineItem.Amount,
			lineItem.Description,
		})
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

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file: ", err)
		os.Exit(1)
	}

	for _, record := range records {
		date := record[0]
		description := record[1]
		amount := record[4]

		exists := false
		for _, lineItem := range lineItems {
			if lineItem.Description == description &&
				lineItem.Date == date &&
				lineItem.Amount == amount {
				exists = true
				continue
			}
		}
		if exists {
			continue
		}

		fmt.Println(record)
		lineItem := LineItem{}
		lineItem.Category = getInput("Category: ")

		if lineItem.Category == "" {
			lineItem.Category = "Ignored"
		}

		lineItem.Date = date
		lineItem.Amount = amount
		lineItem.Description = description
		lineItems = append(lineItems, lineItem)
	}

	return lineItems
}

func deserialize(saveFile string) []LineItem {
	lineItems := make([]LineItem, 0)

	f, err := os.Open(saveFile)
	if err != nil {
		fmt.Println("File does not exist, creating new file")
		return lineItems
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file: ", err)
		os.Exit(1)
	}

	for _, record := range records {
		lineItem := LineItem{
			Category:    record[0],
			Date:        record[1],
			Amount:      record[2],
			Description: record[3],
		}
		lineItems = append(lineItems, lineItem)
	}

	return lineItems
}

func induct(file string) {
	lineItems := deserialize(saveFile)
	lineItems = addNewLineItems(lineItems, file)
	serialize(lineItems)
}

func parseDate(date string) int64 {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		t, err := time.Parse("01/02/2006", date)
		if err != nil {
			return 0
		}
		return t.Unix()
	}
	return t.Unix()
}

func printSummary(startDate string, endDate string) {
	fmt.Printf("Summary from %s to %s\n", startDate, endDate)

	totalAmount := 0.0
	totalNumber := 0

	lineItems := deserialize(saveFile)

	categoryMap := make(map[string]HistogramRow)
	for _, lineItem := range lineItems {
		if parseDate(lineItem.Date) < parseDate(startDate) || parseDate(lineItem.Date) > parseDate(endDate) {
			continue
		}

		a := strings.Replace(lineItem.Amount, "$", "", -1)
		a = strings.Replace(a, ",", "", -1)
		amount, err := strconv.ParseFloat(a, 64)
		if err != nil {
			fmt.Println("Error parsing amount: ", err)
			os.Exit(1)
		}
		categoryMap[lineItem.Category] = HistogramRow{
			Amount: categoryMap[lineItem.Category].Amount + amount,
			Number: categoryMap[lineItem.Category].Number + 1,
		}
	}

	for category, amount := range categoryMap {
		if category == "Ignored" {
			continue
		}

		totalAmount += amount.Amount
		totalNumber += amount.Number

		fmt.Printf("%s: %.2f (%d Transactions)\n",
			category, amount.Amount, amount.Number)
	}

	fmt.Printf("Total: %.2f (%d Transactions)\n", totalAmount, totalNumber)
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
