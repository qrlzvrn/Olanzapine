package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/olekukonko/tablewriter"
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
				Usage: "add a new task",
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
				Usage: "complete the task",
				Action: func(c *cli.Context) error {
					ID := c.Args().First()
					if ID == "" {
						err := fmt.Errorf("task ID cannot be empty")
						return err
					}

					intID, err := strconv.Atoi(ID)
					if err != nil {
						return err
					}

					if err := changeTask(intID, "complete", ""); err != nil {
						return err
					}

					fmt.Printf("task with an id=%v has completed\n", intID)
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "delete the task",
				Action: func(c *cli.Context) error {
					ID := c.Args().First()
					if ID == "" {
						err := fmt.Errorf("task ID cannot be empty")
						return err
					}

					intID, err := strconv.Atoi(ID)
					if err != nil {
						return err
					}

					if err := changeTask(intID, "delete", ""); err != nil {
						return err
					}

					fmt.Printf("task with an id=%v has removed\n", intID)
					return nil
				},
			},
			{
				Name:  "init",
				Usage: "create the table for storing tasks",
				Action: func(c *cli.Context) error {
					initTableTasks()

					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "list the tasks",
				Action: func(c *cli.Context) error {
					if c.NArg() > 0 {
						category := c.Args().First()
						if err := listTasks(category); err != nil {
							return err
						}
					} else {
						if err := listTasks("allCategory"); err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name:    "reDead",
				Aliases: []string{"rd"},
				Usage:   "reassign task deadline",
				Action: func(c *cli.Context) error {
					ID := c.Args().First()
					intID, err := strconv.Atoi(ID)
					if err != nil {
						return err
					}
					newDeadline := c.Args().Get(1)

					changeTask(intID, "reDead", newDeadline)

					fmt.Printf("deadline of task with an ID %v changed on %s\n", intID, newDeadline)

					return nil
				},
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

func listTasks(category string) error {

	db, err := sqlx.Open("postgres", "dbname=olanza sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	data := [][]string{}

	if category == "allCategory" {
		rows, err := db.Queryx("SELECT id, content, complete, category, deadline-now()::date FROM tasks")
		if err != nil {
			return err
		}

		for rows.Next() {
			var id int
			var content string
			var complete bool
			var category string
			var deadline string
			rows.Scan(&id, &content, &complete, &category, &deadline)

			data = append(data, []string{strconv.Itoa(id), content, strconv.FormatBool(complete), category, deadline})
			defer rows.Close()
		}

	} else {
		rows, err := db.Queryx("SELECT id, content, complete, category, deadline-now()::date FROM tasks WHERE category=$1", category)
		if err != nil {
			return err
		}

		for rows.Next() {
			var id int
			var content string
			var complete bool
			var category string
			var deadline string
			rows.Scan(&id, &content, &complete, &category, &deadline)

			data = append(data, []string{strconv.Itoa(id), content, strconv.FormatBool(complete), category, deadline})
			defer rows.Close()
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "content", "complete", "category", "deadline"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgGreenColor})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()

	return nil
}

func initTableTasks() error {
	db, err := sqlx.Open("postgres", "dbname=olanza sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	db.MustExec(table)

	return nil
}

func changeTask(ID int, action string, value string) error {
	db, err := sqlx.Open("postgres", "dbname=olanza sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	if action == "complete" {
		db.MustExec("UPDATE tasks SET complete=true WHERE id = $1", ID)

		return nil
	}

	if action == "delete" {
		db.MustExec("DELETE FROM tasks WHERE id=$1", ID)

		return nil
	}

	if action == "reDead" {
		db.MustExec("UPDATE tasks set deadline=$1 WHERE id=$2", value, ID)

		return nil
	}

	return nil
}
