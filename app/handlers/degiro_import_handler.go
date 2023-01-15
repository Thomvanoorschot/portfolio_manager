package handlers

import (
	"encoding/csv"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const BuyOrSellRegex = "(Koop|Verkoop)\\s(\\d*)\\s@\\s(\\d*,?\\d*)\\s[A-Z]*"

type Commission struct {
	AmountInCents int64
	ExternalId    string
}
type Commissions []*Commission

func DegiroImportHandler(server *server.Webserver, ctx *gin.Context) {
	fileHeader, _ := ctx.FormFile("file")
	file, _ := fileHeader.Open()
	reader := csv.NewReader(file)
	firstLineProcessed := false

	var convertedTransactions entities.Transactions
	var commissions Commissions

	for {
		line, readError := reader.Read()
		if readError == io.EOF {
			break
		} else if readError != nil {
			log.Fatal(readError)
		} else if !firstLineProcessed {
			firstLineProcessed = true
			continue
		}
		convertTransaction(line, &convertedTransactions, &commissions)
	}
	setCommissions(&convertedTransactions, &commissions)
	fmt.Println(len(convertedTransactions))
	server.UnitOfWork.PortfolioRepository.Create(&entities.Portfolio{
		Title:        "Main portfolio",
		Transactions: convertedTransactions,
		EntityBase:   entities.EntityBase{},
	})
}

func setCommissions(transactions *entities.Transactions, commissions *Commissions) {
	for _, commission := range *commissions {
		for _, transaction := range *transactions {
			if transaction.ExternalId == commission.ExternalId {
				transaction.CommissionInCents = commission.AmountInCents
			}
		}
	}
}

func convertTransaction(line []string, transactions *entities.Transactions, commissions *Commissions) {
	isBuyOrSellTransaction, _ := regexp.MatchString(BuyOrSellRegex, line[5])
	if isBuyOrSellTransaction {
		convertBuyOrSellTransaction(line, transactions)
		return
	}

	isCommissionTransaction := line[5] == "DEGIRO Transactiekosten en/of kosten van derden"
	if isCommissionTransaction {
		convertCommission(line, commissions)
		return
	}

	isDeposit := line[5] == "Reservation iDEAL / Sofort Deposit" || line[5] == "iDEAL storting"
	if isDeposit {
		convertDeposit(line, transactions)
		return
	}

	isWithdrawal := line[5] == "Processed Flatex Withdrawal"
	if isWithdrawal {
		convertWithdrawal(line, transactions)
		return
	}
}

func convertCommission(line []string, commissions *Commissions) {
	line[8] = strings.Replace(line[8], ",", ".", -1)
	cost, _ := strconv.ParseFloat(line[8], 64)

	commission := &Commission{
		AmountInCents: int64(cost * 100),
		ExternalId:    line[11],
	}
	*commissions = append(*commissions, commission)
}

func convertDeposit(line []string, transactions *entities.Transactions) {
	transaction := &entities.Transaction{}
	convertGeneralTransactionInfo(line, transaction)
	cost, _ := strconv.ParseFloat(line[8], 64)
	if cost < 0 {
		return
	}
	transaction.Product = "CASH"
	transaction.Amount = 1
	transaction.PriceInCents = int64(cost * 100)
	transaction.TransactionType = entities.Deposit
	*transactions = append(*transactions, transaction)
}
func convertWithdrawal(line []string, transactions *entities.Transactions) {
	transaction := &entities.Transaction{}
	convertGeneralTransactionInfo(line, transaction)
	cost, _ := strconv.ParseFloat(line[8], 64)
	if cost < 0 {
		return
	}
	transaction.Product = "CASH"
	transaction.Amount = 1
	transaction.PriceInCents = int64(cost*100) * -1
	transaction.TransactionType = entities.Withdrawal
	*transactions = append(*transactions, transaction)
}

func convertBuyOrSellTransaction(line []string, transactions *entities.Transactions) {
	transaction := &entities.Transaction{}
	convertGeneralTransactionInfo(line, transaction)
	r, _ := regexp.Compile(BuyOrSellRegex)
	parsedDescription := r.FindStringSubmatch(line[5])

	transactionType, _ := entities.ConvertToTransactionType(parsedDescription[1])
	amount, _ := strconv.ParseFloat(parsedDescription[2], 64)
	if transactionType == entities.Sell {
		amount = amount * -1
	}
	parsedDescription[3] = strings.Replace(parsedDescription[3], ",", ".", -1)
	pricePerUnit, _ := strconv.ParseFloat(parsedDescription[3], 64)
	transaction.Amount = amount
	transaction.PriceInCents = int64(pricePerUnit * 100)
	transaction.TransactionType = transactionType

	*transactions = append(*transactions, transaction)
}
func convertGeneralTransactionInfo(line []string, transaction *entities.Transaction) {
	transactedDate := strings.Split(line[0], "-")
	transactedTime := strings.Split(line[1], ":")
	transactedYear, _ := strconv.Atoi(transactedDate[2])
	transactedMonth, _ := strconv.Atoi(transactedDate[1])
	transactedConvertedMonth := time.Month(transactedMonth)
	transactedDay, _ := strconv.Atoi(transactedDate[0])
	transactedTimeHour, _ := strconv.Atoi(transactedTime[0])
	transactedTimeMinute, _ := strconv.Atoi(transactedTime[1])
	line[8] = strings.Replace(line[8], ",", ".", -1)
	line[10] = strings.Replace(line[10], ",", ".", -1)

	transaction.TransactedAt = time.Date(transactedYear, transactedConvertedMonth, transactedDay, transactedTimeHour, transactedTimeMinute, 0, 0, time.UTC)
	transaction.Product = line[3]
	transaction.ISIN = line[4]
	transaction.Description = line[5]
	transaction.ExternalId = line[11]
	transaction.CurrencyCode = line[9]
}
