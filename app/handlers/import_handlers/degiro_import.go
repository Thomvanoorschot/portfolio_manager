package import_handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/enums"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	BuyOrSellRegex = "(Koop|Verkoop)\\s(\\d*)\\s@\\s(\\d*,?\\d*)\\s[A-Z]*"
)

var (
	SymbolMap = map[string]string{}
)

type Commission struct {
	AmountInCents int64
	ExternalId    string
}
type Commissions []*Commission

type YahooSearch struct {
	Quotes []struct {
		Symbol string `json:"symbol"`
	} `json:"quotes"`
}

func DegiroImport(server *server.Webserver, ctx *gin.Context) {
	fileHeader, _ := ctx.FormFile("file")
	portfolioId := ctx.Request.Form.Get("portfolioId")
	file, _ := fileHeader.Open()
	reader := csv.NewReader(file)
	firstLineProcessed := false

	var convertedTransactions entities.Transactions
	var commissions Commissions
	portfolioUuid := uuid.MustParse(portfolioId)
	var previouslyConvertedTransactions entities.Transactions
	if portfolioId != "" {
		previouslyConvertedTransactions = server.UnitOfWork.TransactionRepository.GetByPortfolioId(portfolioUuid)
	}
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
		convertTransaction(line, &convertedTransactions, &commissions, previouslyConvertedTransactions, portfolioUuid)
	}
	if len(convertedTransactions) == 0 {
		return
	}
	setCommissions(&convertedTransactions, &commissions)
	if portfolioId == "" {
		server.UnitOfWork.PortfolioRepository.Create(&entities.Portfolio{
			Title:        "Main portfolio",
			Transactions: convertedTransactions,
			EntityBase:   entities.EntityBase{},
		})
	} else {
		server.UnitOfWork.TransactionRepository.AddToPortfolio(convertedTransactions)
	}
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

func convertTransaction(line []string, transactions *entities.Transactions,
	commissions *Commissions,
	previouslyConvertedTransactions entities.Transactions,
	portfolioId uuid.UUID) {

	uniqueHash := computeHashForList(line)
	for _, previousConvertedTransaction := range previouslyConvertedTransactions {
		if uniqueHash == previousConvertedTransaction.UniqueHash {
			return
		}
	}

	isBuyOrSellTransaction, _ := regexp.MatchString(BuyOrSellRegex, line[5])
	if isBuyOrSellTransaction {
		convertBuyOrSellTransaction(line, transactions, portfolioId)
		return
	}

	isCommissionTransaction := line[5] == "DEGIRO Transactiekosten en/of kosten van derden"
	if isCommissionTransaction {
		convertCommission(line, commissions)
		return
	}

	isDeposit := line[5] == "Reservation iDEAL / Sofort Deposit" || line[5] == "iDEAL storting"
	if isDeposit {
		convertDeposit(line, transactions, portfolioId)
		return
	}

	isWithdrawal := line[5] == "Processed Flatex Withdrawal"
	if isWithdrawal {
		convertWithdrawal(line, transactions, portfolioId)
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

func convertDeposit(line []string,
	transactions *entities.Transactions,
	portfolioId uuid.UUID) {
	transaction := &entities.Transaction{}
	convertGeneralTransactionInfo(line, transaction, portfolioId)
	cost, _ := strconv.ParseFloat(line[8], 64)
	if cost < 0 {
		return
	}
	transaction.Product = "CASH"
	transaction.Amount = 1
	transaction.PriceInCents = int64(cost * 100)
	transaction.TransactionType = enums.Deposit
	*transactions = append(*transactions, transaction)
}
func convertWithdrawal(line []string,
	transactions *entities.Transactions,
	portfolioId uuid.UUID) {
	transaction := &entities.Transaction{}
	convertGeneralTransactionInfo(line, transaction, portfolioId)
	cost, _ := strconv.ParseFloat(line[8], 64)
	if cost < 0 {
		return
	}
	transaction.Product = "CASH"
	transaction.Amount = 1
	transaction.PriceInCents = int64(cost*100) * -1
	transaction.TransactionType = enums.Withdrawal
	*transactions = append(*transactions, transaction)
}

func convertBuyOrSellTransaction(line []string,
	transactions *entities.Transactions,
	portfolioId uuid.UUID) {
	transaction := &entities.Transaction{}
	convertGeneralTransactionInfo(line, transaction, portfolioId)
	r, _ := regexp.Compile(BuyOrSellRegex)
	parsedDescription := r.FindStringSubmatch(line[5])

	transactionType, _ := entities.ConvertToTransactionType(parsedDescription[1])
	amount, _ := strconv.ParseFloat(parsedDescription[2], 64)
	if transactionType == enums.Sell {
		amount = amount * -1
	}
	parsedDescription[3] = strings.Replace(parsedDescription[3], ",", ".", -1)
	pricePerUnit, _ := strconv.ParseFloat(parsedDescription[3], 64)
	transaction.Amount = amount
	transaction.PriceInCents = int64(pricePerUnit * 100)
	transaction.TransactionType = transactionType
	transaction.Symbol = getSymbol(transaction.ISIN)
	*transactions = append(*transactions, transaction)
}
func convertGeneralTransactionInfo(line []string,
	transaction *entities.Transaction,
	portfolioId uuid.UUID) {
	uniqueHash := computeHashForList(line)
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
	transaction.UniqueHash = uniqueHash
	transaction.PortfolioID = portfolioId
}

func getSymbol(searchTerm string) string {
	symbol, found := SymbolMap[searchTerm]
	if found {
		return symbol
	}
	url := fmt.Sprintf("https://query2.finance.yahoo.com/v1/finance/search?q=%s&lang=en-US&region=US&quotesCount=1&newsCount=0&listsCount=0&enableFuzzyQuery=false&quotesQueryId=tss_match_phrase_query&multiQuoteQueryId=multi_quote_single_token_query&newsQueryId=news_cie_vespa&enableCb=true&enableNavLinks=false&enableEnhancedTrivialQuery=false&enableResearchReports=false&enableCulturalAssets=false&enableLogoUrl=false&researchReportsCount=0", searchTerm)
	r, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)
	searchResult := &YahooSearch{}
	_ = json.NewDecoder(r.Body).Decode(searchResult)

	if len(searchResult.Quotes) == 0 {
		SymbolMap[searchTerm] = ""
		return ""
	}
	SymbolMap[searchTerm] = searchResult.Quotes[0].Symbol
	return searchResult.Quotes[0].Symbol
}

func computeHashForList(list []string) string {
	var buffer bytes.Buffer
	for i := range list {
		buffer.WriteString(list[i])
		buffer.WriteString("0")
	}
	b := sha256.Sum256([]byte(buffer.String()))
	s := hex.EncodeToString(b[:])
	return s
}
