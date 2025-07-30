package main

import (
	"fmt"
	"os"
	"time"
	"log"
	"encoding/json"
	flag "github.com/spf13/pflag"//pflag потому что флаг при передаче не флага начинает по другому парсить аргументы(что то связанное с позиционированием), кароче не работает когда progname add(действие) --тут уже флаги
)


func main(){

	expenses := make([]Expense, 0) // слайс трат
	args := os.Args

	if len(args) == 1{
		//без аргументов выводим приветсвие
		greeting()
		return
	} else{

	readTratyFromFile(&expenses) //слайс трат передаем по указателю иначе изменения будут только в копии в вызываемой фукнции
	
	action := args[1]

	switch action{
	case "list", "l":
		//вывод всех трат

		if len(args) > 2{
			log.Fatal("wrong input")
		}
		printExpenses(expenses)

	case "add", "a":
		//добавление траты

		desc := flag.StringP("description", "d", "", "a description for expense")//парсим флаги
		amount := flag.IntP("amount", "a", 0, "an amount of money")
		flag.Parse()
		//проверяем что параметры верно введены
		if *desc == "" || *amount == 0{
			log.Fatal("wrong input parametrs")
		} 

		t := time.Now()
		date := t.Format("2006-01-02 15:04:05") // засекли время создания покупочки

		exp := Expense{Id: len(expenses)+1, Date: date, Descr: *desc, Amount: *amount} //создаем структуру Expense
		
		expenses = append(expenses, exp) //добавляем к слайсу трат новую 
		
		writeToFileJson(expenses)//пишем слайс в json file

		fmt.Printf("Expense added succesfully (ID: %d)", len(expenses))

	case "summary", "s":
		//суммарное колво потраченных денжат
		if len(args) > 2{
			log.Fatal("wrong input")
		}
		var sum int
		sum = getSum(expenses)
		fmt.Printf("Total expenses: $%d", sum)

	case "delete", "d":
		id := flag.IntP("id", "i", -1, "an id for delete expense")
		flag.Parse()
		*id--
		if *id < 0 || *id >= len(expenses){
			log.Fatal("wrong input parametrs")
		}
		expenses = rmExpense(expenses, *id)

		writeToFileJson(expenses)

		fmt.Printf("Expense deleted succesfully\n")
	}


	}

}

type Expense struct{
	Id	int			`json:"id"`
	Date	string	`json:"date"`
	Descr	string	`json:"description"`
	Amount	int		`json:"amount"`
}

//вывод трат
func printExpenses(exps []Expense){
	if len(exps) == 0{
		fmt.Println("There is no expenses")
		return
	}
	fmt.Println("ID\tDate\t\t\tDescription\tAmount")
	for _, v := range exps{
		fmt.Printf("%d\t%s\t%s\t\t$%d\n", v.Id, v.Date, v.Descr, v.Amount)
	}
}
//подсчет суммы трат
func getSum(exps []Expense) (int) {
	var s int
	for i := 0; i < len(exps); i++{
		s += exps[i].Amount
	}
	return s
}

//удаление траты
func rmExpense(exps []Expense, id int) ([]Expense){
	newExps := append(exps[:id], exps[id+1:]...)
	for i := 0; i < len(newExps); i++{
		newExps[i].Id = i + 1 //обновляем id
	}
	return newExps
}

//работа с json
func readTratyFromFile(exps *[]Expense){
	bytesFromFile, _ := os.ReadFile("json/expenses.json")//считали байты
	err := json.Unmarshal(bytesFromFile, &exps)//распаковали байта в слайс трат
	if err != nil && len(*exps) != 0{log.Fatal(err)}
}

func writeToFileJson(exps []Expense){
	file, err := os.Create("json/expenses.json")
	if err != nil{ log.Fatal(err)}
	defer file.Close()
	tratyJson, err := json.Marshal(exps)//переводим слайс в байты, которые пишем в файл
	if err != nil { log.Fatal(err)}
	file.Write(tratyJson)
}

//приветсвие
func greeting(){
	fmt.Println("\t\tExpense Tracker")
	fmt.Println("Use: tracker [command] [flags])")
	fmt.Println("command: add --description|-d \"something\" --amount|-a \"some\"")
	fmt.Println("\t list|l //list of expenses")
	fmt.Println("\t summary|s // summary expenses")
	fmt.Println("\t delete|d --id|-i 2 //delete some expense")
}

