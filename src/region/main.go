package main

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type RegionModel struct {
	Pid  int
	Id   int
	Code string
	Name string
}

type DfController struct {
	beego.Controller
}

func (d *DfController) Get() {
	d.findAll()
	d.TplName = "index.html"
}

func (d *DfController) Post() {
	d.TplName = "index.html"
	d.Data["res"] = false
	pid := d.GetInt("pregionid", 0)
	name := d.GetString("region_name", "")
	if len(name) < 1 {
		d.Data["msg"] = "地名为空"
		return
	}
	code := d.GetString("region_code", "")
	if len(code) < 1 {
		d.Data["msg"] = "行政代码为空"
		return
	}

	//添加数据
	_, err := mysqlDb.Exec("insert into region(pRegionId,regionName,regionCode) values(?,?,?);", pid, name, code)
	if err != nil {
		d.Data["msg"] = fmt.Sprintf("添加数据出错：%v	", err)
		return
	}

	d.findAll()
	d.Data["res"] = true
	d.Data["msg"] = "添加成功"
}

func (d *DfController) Json() {
	pid := d.GetInt("pregionid", 0)
	rows, err := mysqlDb.Query("select regionId,regionName from region where regionId=? order by regionId asc", pid)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return
	}

	regions := make([]RegionModel, 0)
	for rows.Next() {
		var (
			rName, rCode string
			rId, rPId    int
		)
		err = rows.Scan(&rId, rPId, rName, rCode)
		if err != nil {
			fmt.Println(err)
		}
		region := RegionModel{
			Pid:  rPId,
			Id:   rId,
			Name: rName,
			Code: rCode,
		}
		regions = append(region, region)
	}

	d.Data["json"] = regions
	d.ServeJSON()
}

func (d *DfController) FindAll() {
	rows, err := mysqlDb.Query("select regionId,pRegionId,regionName,regionCode from region order by regionId asc")
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return
	}

	regions := make([]RegionModel, 0)
	for rows.Next() {
		var (
			rName, rCode string
			rId, rPId    int
		)
		err = rows.Scan(&rId, rPId, rName, rCode)
		if err != nil {
			fmt.Println(err)
		}
		region := &RegionModel{
			Pid:  rPId,
			Id:   rId,
			Name: rName,
			Code: rCode,
		}

		regions = append(region, region)
	}

	d.Data["json"] = regions
	d.ServeJSON()
}

var mysqlDb *sql.DB

func main() {

	//database
	var err error
	mysqlDb, err = sql.Open("mysql", "zhuiju365:zhuiju365@tcp(114.215.99.36:3306)/f273c?charset=utf8&loc=Local")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	err = mysqlDb.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	df := new(DfController{})

	beego.Router("/", df, "Get:Get")
	beego.Router("/", df, "Post:Post")
	beego.Router("/json", df, "*:Json")
	beego.Router("/findAll", df, "*:FindAll")

	beego.Run()
}
