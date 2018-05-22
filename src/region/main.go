package main

import (
	"os"
	"fmt"
	"bufio"
	"io"
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
	fmt.Println(regions)
}
