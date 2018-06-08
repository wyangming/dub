package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Pxy struct {
}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, re *http.Request) {
	fmt.Println(re.URL.Path)

}

func main() {
	//	http.ListenAndServe(":81", &Pxy{})
	//addRegion()
	str := "/aa/bb"
	fmt.Println(len(strings.Split(str, "/")))
	fmt.Println(strings.Split(str, "/"))
}

func rpcDemo() {
	client, err := rpc.DialHTTP("tcp", ":7020")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("aaa")
	args := "dubing"
	var reply string
	err = client.Call("UserRpc.FirstName", &args, &reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
}

func addRegion() {
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

	dbSql, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
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
	stm, err := dbSql.Prepare("insert into dubregion(regionName,regionCode) values(?,?)")
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
