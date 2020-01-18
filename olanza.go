package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"

	_ "github.com/lib/pq"
)

//Структура, в которую считываются данные введенные пользователем при добавлении новой задачи
type task struct {
	content  string
	category string
	deadline string
}

var table = `
CREATE TABLE tasks (
	id SERIAL,
	content TEXT NOT NULL,
	complete BOOLEAN,
	category VARCHAR,
	deadline DATE
)`

func main() {
	app := &cli.App{
		Name:  "olanza",
		Usage: "Wait a minute",
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "content",
						Aliases:  []string{"C"},
						Usage:    "set task content",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Value:   "Purgatorium",
						Usage:   "set task category",
					},
					&cli.StringFlag{
						Name:    "deadline",
						Aliases: []string{"d"},
						Value:   "NULL",
						Usage:   "set task deadline",
					},
				},
				Action: func(c *cli.Context) error {
					t := task{
						content:  c.String("content"),
						category: c.String("category"),
						deadline: c.String("deadline"),
					}

					if t.content == "" {
						err := fmt.Errorf("task content can't be empty")
						return err
					}

					if err := addTask(t); err != nil {
						return err
					}

					fmt.Println("added new task: ", t.content)

					return nil
				},
			},
			{
				Name:  "complete",
				Usage: "",
				Action: func(c *cli.Context) error {
					ID := c.Args().First()
					if ID == "" {
						err := fmt.Errorf("task ID cannot be empty")
						return err
					}

					IDint, err := strconv.Atoi(ID)
					if err != nil {
						return err
					}

					if err := completeTask(IDint); err != nil {
						return err
					}
					fmt.Printf("task with an id=%v has completed\n", IDint)
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "",
				Action: func(c *cli.Context) error {
					ID := c.Args().First()
					if ID == "" {
						err := fmt.Errorf("task ID cannot be empty")
						return err
					}

					IDint, err := strconv.Atoi(ID)
					if err != nil {
						return err
					}

					if err := deleteTask(IDint); err != nil {
						return err
					}
					fmt.Printf("task with an id=%v has removed\n", IDint)
					return nil
				},
			},
			{
				Name:  "init",
				Usage: "",
				//Action: ,
				//проверяем существование бд
				//Если все в порядке, то создаем таблицу
				//иначе выбрасываем ошибку и сообщаем пользователю, что базы данных не существует и ее нужно создать
				//закрываем соединение с бд

			},
			{
				Name:    "list",
				Aliases: []string{"s"},
				Usage:   "",
				//Action: ,
				//если пользователь вводит еще и категорию, то выводим только подходящие
				//если не введено ничего, то выводим все задачи
				//если получаем ключ --all/-a, то выводим даже выполненные задачи
				//открывем базу; делаем SELECT, считвыаем все в [][]string
				//строим из этих данных таблицу, что бы все аккуратно выводилось в консоль
				//закрываем соединение с бд
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func addTask(t task) error {
	db, err := sqlx.Open("postgres", "dbname=olanza sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.MustBegin()

	if t.deadline == "NULL" {
		tx.MustExec("INSERT INTO tasks (content, category, complete) VALUES ($1, $2, $3)", t.content, t.category, false)
	} else {
		tx.MustExec("INSERT INTO tasks (content, category, deadline, complete) VALUES ($1, $2, $3, $4)", t.content, t.category, t.deadline, false)
	}
	tx.Commit()
	return nil
}

func completeTask(ID int) error {
	db, err := sqlx.Open("postgres", "dbname=olanza sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.MustBegin()
	tx.MustExec("UPDATE tasks SET complete=true WHERE id = $1", ID)
	tx.Commit()

	return nil
}

func deleteTask(ID int) error {
	db, err := sqlx.Open("postgres", "dbname=olanza sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.MustBegin()
	tx.MustExec("DELETE FROM tasks WHERE id=$1", ID)
	tx.Commit()

	return nil
}

func listTasks() error {

	return nil
}

func initTableTasks() error {

	return nil
}
