package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"math"
	"time"
)

const (
	DB_USER            = "postgres"
	DB_PASSWORD        = "postgres"
	DB_NAME            = "tunaiku"
	LimitPrime         = 100 // this is limit of prine number function
	EvenOddLimit       = 100 // this is limit of even odd function
	NominalAmountStart = 2000000
	NominalAmountEnd   = 15000000
)

type PrimeNumber struct {
	Numbers    int
	NameNumber string
}
type NominalAmount struct {
	Nominal          float64
	AdditionalNumber int
	Total            float64
}
type EvenOddNumber struct {
	NumberEven int
	NumberOdd  int
	Total      int
	NameTotal  string
}

func main() {

	primeNumberFunc()
	EvenOddNumberFunc()
	NominalAmountFunc()
	defer elapsed("execution time")()
	time.Sleep(time.Second * 2)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//function to get prime number
func primeNumberFunc() {
	var primListNumber []PrimeNumber
	listPrime := countPrimes()
	for i := 0; i < len(listPrime); i++ {
		var primeNumber PrimeNumber
		primeNumber.Numbers = listPrime[i]
		primeNumber.NameNumber = Convert(listPrime[i])
		primListNumber = append(primListNumber, primeNumber)
	}
	// columName := fmt.Sprintf("(%v,%v", "numbers", "name_number")
	columnName := []string{"numbers", "name_number"}
	// t := []int{1, 2, 3, 4}
	s := make([]interface{}, len(primListNumber))
	for i, v := range primListNumber {
		s[i] = v
	}
	InsertDB(columnName, s)
}

//function to count nominal amount
func NominalAmountFunc() {
	var NominalAmountList []NominalAmount
	columnName := []string{"nominal", "additional_number", "total"}
	endNominal := NominalAmountEnd
	beginNominal := float64(NominalAmountStart)
	for beginNominal <= float64(endNominal) {
		var nominalAmount NominalAmount
		// for i := NominalAmountStart; NominalAmountStart <= endNominal; i += total {
		fmt.Println(beginNominal, "<=", endNominal)
		interest, total := getInterestAndTotal(beginNominal)
		nominalAmount.Nominal = beginNominal
		nominalAmount.AdditionalNumber = interest
		nominalAmount.Total = total
		NominalAmountList = append(NominalAmountList, nominalAmount)
		fmt.Println(interest, total)
		beginNominal = total
	}
	fmt.Println(NominalAmountList)
	s := make([]interface{}, len(NominalAmountList))
	for i, v := range NominalAmountList {
		s[i] = v
	}
	InsertDB(columnName, s)

}

func getInterestAndTotal(nominal float64) (int, float64) {
	interest := 0.1 * nominal
	total := nominal + interest
	return int(interest), total
}

//function odd even number
func EvenOddNumberFunc() {
	var EvenOddNumberList []EvenOddNumber
	var EvenOddNumber EvenOddNumber
	var EvenNumber, OddNumber []int
	columnName := []string{"number_even", "number_odds", "total", "name_total"}
	for i := 0; i < EvenOddLimit; i++ {
		if i%2 == 0 {
			EvenNumber = append(EvenNumber, i)
		} else {
			OddNumber = append(OddNumber, i)
		}
		// i += 2
	}
	for i := 0; i < len(EvenNumber); i++ {
		EvenOddNumber.NumberEven = EvenNumber[i]
		EvenOddNumber.NumberOdd = OddNumber[i]
		EvenOddNumber.Total = EvenNumber[i] + OddNumber[i]
		EvenOddNumber.NameTotal = Convert(EvenOddNumber.Total)
		EvenOddNumberList = append(EvenOddNumberList, EvenOddNumber)
	}

	s := make([]interface{}, len(EvenOddNumberList))
	for i, v := range EvenOddNumberList {
		s[i] = v
	}
	InsertDB(columnName, s)

}

func InsertDB(ColumnName []string, val []interface{}) {
	var sqlStr string
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, _ := sql.Open("postgres", dbinfo)
	// checkErr(err)
	defer db.Close()

	fmt.Println("# Inserting values")
	// query := fmt.Sprintf("INSERT INTO %v(%v) VALUES,", DB_NAME, ColumnName)
	// sqlStr := fmt.Sprintf("%v", query)
	fmt.Println(val)
	byteData, _ := json.Marshal(val)
	switch len(ColumnName) {
	case 2:
		sqlStr = fmt.Sprintf("INSERT INTO prime_number(%v, %v) VALUES ", ColumnName[0], ColumnName[1])
		var listPrimNumber []PrimeNumber
		json.Unmarshal(byteData, &listPrimNumber)
		for i := range listPrimNumber {
			if i == 0 {
				sqlStr += fmt.Sprintf("(%v,'%v')", listPrimNumber[i].Numbers, listPrimNumber[i].NameNumber)
			} else {
				sqlStr += fmt.Sprintf(",(%v,'%v')", listPrimNumber[i].Numbers, listPrimNumber[i].NameNumber)
			}
		}
	case 3:
		sqlStr = fmt.Sprintf("INSERT INTO nominal_amount(%v, %v, %v) VALUES ", ColumnName[0], ColumnName[1], ColumnName[2])
		var listNominalAmount []NominalAmount
		json.Unmarshal(byteData, &listNominalAmount)
		for i := range listNominalAmount {
			if i == 0 {
				sqlStr += fmt.Sprintf("(%v,'%v',%v)", listNominalAmount[i].Nominal, listNominalAmount[i].AdditionalNumber, listNominalAmount[i].Total)
			} else {
				sqlStr += fmt.Sprintf(",(%v,'%v',%v)", listNominalAmount[i].Nominal, listNominalAmount[i].AdditionalNumber, listNominalAmount[i].Total)
			}
		}
	case 4:
		var listOddEvenNumber []EvenOddNumber
		json.Unmarshal(byteData, &listOddEvenNumber)
		sqlStr = fmt.Sprintf("INSERT INTO even_odd_number(%v, %v, %v,%v) VALUES ", ColumnName[0], ColumnName[1], ColumnName[2], ColumnName[3])
		for i := range listOddEvenNumber {
			if i == 0 {
				sqlStr += fmt.Sprintf("(%v,'%v',%v,'%v')", listOddEvenNumber[i].NumberEven, listOddEvenNumber[i].NumberOdd, listOddEvenNumber[i].Total, listOddEvenNumber[i].NameTotal)
			} else {
				sqlStr += fmt.Sprintf(",(%v,'%v',%v,'%v')", listOddEvenNumber[i].NumberEven, listOddEvenNumber[i].NumberOdd, listOddEvenNumber[i].Total, listOddEvenNumber[i].NameTotal)
			}
		}
	}
	fmt.Println(sqlStr)
	_, err := db.Exec(sqlStr)
	// fmt.Println(res)
	fmt.Println(err)
}

func countPrimes() []int {
	var x, y, n int
	nsqrt := math.Sqrt(LimitPrime)

	is_prime := [LimitPrime]bool{}

	for x = 1; float64(x) <= nsqrt; x++ {
		for y = 1; float64(y) <= nsqrt; y++ {
			n = 4*(x*x) + y*y
			if n <= LimitPrime && (n%12 == 1 || n%12 == 5) {
				is_prime[n] = !is_prime[n]
			}
			n = 3*(x*x) + y*y
			if n <= LimitPrime && n%12 == 7 {
				is_prime[n] = !is_prime[n]
			}
			n = 3*(x*x) - y*y
			if x > y && n <= LimitPrime && n%12 == 11 {
				is_prime[n] = !is_prime[n]
			}
		}
	}

	for n = 5; float64(n) <= nsqrt; n++ {
		if is_prime[n] {
			for y = n * n; y < LimitPrime; y += n * n {
				is_prime[y] = false
			}
		}
	}

	is_prime[2] = true
	is_prime[3] = true

	primes := make([]int, 0, 1270606)
	for x = 0; x < len(is_prime)-1; x++ {
		if is_prime[x] {
			primes = append(primes, x)
		}
	}
	return primes
}

// func descint

//start

// how many digit's groups to process
const groupsNumber int = 4

var _smallNumbers = []string{
	"zero", "one", "two", "three", "four",
	"five", "six", "seven", "eight", "nine",
	"ten", "eleven", "twelve", "thirteen", "fourteen",
	"fifteen", "sixteen", "seventeen", "eighteen", "nineteen",
}
var _tens = []string{
	"", "", "twenty", "thirty", "forty", "fifty",
	"sixty", "seventy", "eighty", "ninety",
}
var _scaleNumbers = []string{
	"", "thousand", "million", "billion",
}

type digitGroup int

// Convert converts number into the words representation.
func Convert(number int) string {
	return convert(number, false)
}

// ConvertAnd converts number into the words representation
// with " and " added between number groups.
func ConvertAnd(number int) string {
	return convert(number, true)
}

func convert(number int, useAnd bool) string {
	// Zero rule
	if number == 0 {
		return _smallNumbers[0]
	}

	// Divide into three-digits group
	var groups [groupsNumber]digitGroup
	positive := math.Abs(float64(number))

	// Form three-digit groups
	for i := 0; i < groupsNumber; i++ {
		groups[i] = digitGroup(math.Mod(positive, 1000))
		positive /= 1000
	}

	var textGroup [groupsNumber]string
	for i := 0; i < groupsNumber; i++ {
		textGroup[i] = digitGroup2Text(groups[i], useAnd)
	}
	combined := textGroup[0]
	and := useAnd && (groups[0] > 0 && groups[0] < 100)

	for i := 1; i < groupsNumber; i++ {
		if groups[i] != 0 {
			prefix := textGroup[i] + " " + _scaleNumbers[i]

			if len(combined) != 0 {
				prefix += separator(and)
			}

			and = false

			combined = prefix + combined
		}
	}

	if number < 0 {
		combined = "minus " + combined
	}

	return combined
}

func intMod(x, y int) int {
	return int(math.Mod(float64(x), float64(y)))
}

func digitGroup2Text(group digitGroup, useAnd bool) (ret string) {
	hundreds := group / 100
	tensUnits := intMod(int(group), 100)

	if hundreds != 0 {
		ret += _smallNumbers[hundreds] + " hundred"

		if tensUnits != 0 {
			ret += separator(useAnd)
		}
	}

	tens := tensUnits / 10
	units := intMod(tensUnits, 10)

	if tens >= 2 {
		ret += _tens[tens]

		if units != 0 {
			ret += "-" + _smallNumbers[units]
		}
	} else if tensUnits != 0 {
		ret += _smallNumbers[tensUnits]
	}

	return
}

// separator returns proper separator string between
// number groups.
func separator(useAnd bool) string {
	if useAnd {
		return " and "
	}
	return " "
}

//end

// getexecutiontime
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v after 2s\n", what, time.Since(start))
	}
}
