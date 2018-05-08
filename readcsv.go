package main

import (
	"bufio"
	"encoding/csv"
	// "encoding/json"
	// "fmt"
	"io"
	"log"
	"os"
	// "reflect"
	"strconv"
	"strings"
	"time"
)

type TunaikuStock struct {
	Date     time.Time `json:"date"`
	Open     int       `json:"open"`
	High     int       `json:"high"`
	Low      int       `json:"low"`
	Close    int       `json:"close"`
	AdjClose int       `json:"adj_close"`
	Volume   int       `json:"volume"`
}

type BuySale struct {
	Buy  int
	Sell int
}

var boughtDayList []int
var soldDayList []int
var closePriceList []int

func main() {
	log.SetFlags(log.Lshortfile)
	fileName := "data.csv"
	dataList := GetFromCsv(fileName)
	// fmt.Println(dataList)
	BuySellStock(dataList)
	for i := 0; i < len(soldDayList); i++ {
		log.Printf("buy on day %v worth %v, sale on day %v worth %v", boughtDayList[i], dataList[boughtDayList[i+1]].Open, soldDayList[i], dataList[soldDayList[i]+1].Close)
	}
}

func GetIntFromString(num string) int {
	var numInt int
	stringNum := strings.Split(num, ".")
	if len(stringNum) > 1 {
		numInt, _ = strconv.Atoi(stringNum[0])
	} else {
		numInt, _ = strconv.Atoi(num)
	}
	return numInt
}

func GetFromCsv(fileName string) []TunaikuStock {
	layoutDate := "2006-01-02"
	csvFile, _ := os.Open(fileName)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var tunaikuStock []TunaikuStock
	for {

		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		dateStock, _ := time.Parse(layoutDate, line[0])
		openStock := GetIntFromString(line[1])
		highStock := GetIntFromString(line[2])
		lowStock := GetIntFromString(line[3])
		closeStock := GetIntFromString(line[4])
		adjCloseStock := GetIntFromString(line[5])
		volumeStock := GetIntFromString(line[6])
		tunaikuStock = append(tunaikuStock, TunaikuStock{
			Date:     dateStock,
			Open:     openStock,
			High:     highStock,
			Low:      lowStock,
			Close:    closeStock,
			AdjClose: adjCloseStock,
			Volume:   volumeStock,
		})
	}
	// fmt.Println(tunaikuStock)
	// tunaikuStockJson, _ := json.Marshal(tunaikuStock)
	// fmt.Println(string(tunaikuStockJson))

	return tunaikuStock
}

func BuySellStock(dataLists []TunaikuStock) ([]int, []int) {
	var dataList []TunaikuStock
	var nextSellPrice int
	dataList = dataLists[1:len(dataLists)]
	// var checkPoint bool
	bools := true
	// maxClosePrice := FindMaxInSlice(dataList)
	_ = StoreBoughtDayStock(0)
	i := 0
	lenData := len(dataList) - 1
	boughtPrice := dataList[0].Open
	for bools {
		if i < len(dataList) {
			if i <= lenData {
				if i == lenData {
					nextSellPrice = dataList[i].Close
				} else {
					nextSellPrice = dataList[i+1].Close
				}
				checkClosePrice := CheckSellTheFuture(dataList[i].Close, nextSellPrice)
				checkOpenPrice := CheckPrice(dataList[i].Close, boughtPrice)
				if checkClosePrice && checkOpenPrice {
					currentBoughtDay := StoreSoldDayStock(i)
					_ = StoreBoughtDayStock(currentBoughtDay)
					boughtPrice = dataList[currentBoughtDay].Open

				}
				if i <= len(dataList) {
					i += 1
				} else if i == (lenData) {
					break
				}
			}
		} else {
			break
		}
	}
	return nil, nil
}

func CheckSellTheFuture(todaySellPrice, tomorrowSellPrice int) bool {
	var meetPrice bool
	if todaySellPrice > tomorrowSellPrice {
		meetPrice = true
	} else {
		meetPrice = false
	}
	return meetPrice
}

func StoreBoughtDayStock(boughtDay int) int {
	boughtDayList = append(boughtDayList, boughtDay)
	return boughtDay + 1
}

func StoreSoldDayStock(soldDay int) int {
	soldDayList = append(soldDayList, soldDay)
	return soldDay + 1
}

func FindMaxInSlice(array []TunaikuStock) int {
	// var arrayOfInt []int
	for i := range array {
		closePriceList = append(closePriceList, array[i].Close)
	}
	var max int = closePriceList[0]
	// var min int = array[0]
	for _, value := range closePriceList {
		if max < value {
			max = value
		}
	}
	return max
}

func CheckPrice(soldPrice, boughtPrice int) bool {
	if soldPrice > boughtPrice {
		return true
	} else {
		return false
	}
}
