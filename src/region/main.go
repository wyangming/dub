package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("xzgh.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	defer f.Close()
	br := bufio.NewReader(f)

	regions := make([][]string, 0)

	for {
		data, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
			break
		}
		str := string(data)
		str = strings.Trim(str, " ")
		if len(str) < 1 {
			break
		}

		regions = append(regions, strings.Split(str, "-"))
	}

	dbSql, err := sql.Open("mysql", "zhuiju365:zhuiju365@tcp(114.215.99.36:3306)/f273c?charset=utf8&loc=Local")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	err = dbSql.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	count := 0
	length := len(regions)
	stm, err := dbSql.Prepare("insert into region(regionName,regionCode) values(?,?)")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	for count < length {

		stm.Exec(regions[count][0], regions[count][1])
		fmt.Println(regions[count])

		count++
	}
	dbSql.Close()
}
