package main

import (
	"context"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/mattn/go-sqlite3"
)

type ProgrammingLanguage struct {
	Id            int    `orm:"auto"`
	Name          string `orm:"size(100)"`
	Paradigm      string `orm:"size(100)"`
	FirstAppeared int
}

func (l *ProgrammingLanguage) FirstAppearedStr() string {
	return strconv.Itoa(l.FirstAppeared)
}

func (l *ProgrammingLanguage) IdStr() string {
	return strconv.Itoa(l.Id)
}

func init() {
	orm.RegisterModel(new(ProgrammingLanguage))
	orm.RegisterDataBase("default", "sqlite3", "data.db")
}

type LanguageController struct {
	web.Controller
}

func (c *LanguageController) List() {
	var languages []ProgrammingLanguage
	o.QueryTable("programming_language").All(&languages)
	LanguageList(languages).Render(context.Background(), c.Ctx.ResponseWriter)
}

func (c *LanguageController) AddForm() {
	AddLanguageForm().Render(context.Background(), c.Ctx.ResponseWriter)
}

func (c *LanguageController) Add() {
	var lang ProgrammingLanguage
	if err := c.ParseForm(&lang); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.WriteString("Invalid input")
		return
	}
	if _, err := o.Insert(&lang); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to add language")
		return
	}
	c.Redirect("/", 200)
}

func (c *LanguageController) Details() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.WriteString("Invalid ID")
		return
	}
	lang := ProgrammingLanguage{Id: id}
	if err := o.Read(&lang); err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Ctx.WriteString("Language not found")
		return
	}
	LanguageDetails(lang).Render(context.Background(), c.Ctx.ResponseWriter)
}

func (c *LanguageController) Delete() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.WriteString("Invalid ID")
		return
	}
	if _, err := o.Delete(&ProgrammingLanguage{Id: id}); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to delete language")
		return
	}
	c.Redirect("/", 200)
}
var o orm.Ormer

func main() {
	o = orm.NewOrm()
	orm.RunSyncdb("default", false, true)
	web.Router("/", &LanguageController{}, "get:List")
	web.Router("/languages/add", &LanguageController{}, "get:AddForm;post:Add")
	web.Router("/languages/:id", &LanguageController{}, "get:Details")
	web.Router("/languages/delete/:id", &LanguageController{}, "get:Delete")
	web.Run()
}
