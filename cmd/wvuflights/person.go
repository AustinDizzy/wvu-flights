package main

import (
	"log"

	flights "github.com/austindizzy/wvu-flights/internal/wvuflights"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v3"
)

var (
	personAddCmd = &cli.Command{
		Name:   "add",
		Usage:  "Add a person",
		Action: personAdd,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "type",
				Usage: "the person's type (wvu | lja)",
			},
		},
	}
	personLookupCmd = &cli.Command{
		Name:   "lookup",
		Usage:  "Lookup a person",
		Action: personLookup,
	}
)

func personAdd(c *cli.Context) error {
	name := c.String("name")
	if name == "" {
		prompt := promptui.Prompt{
			Label: "Person Name",
		}
		var err error
		name, err = prompt.Run()
		if err != nil {
			return err
		}
	}

	pType := c.String("type")
	if pType == "" {
		prompt := promptui.Select{
			Label: "Person Type",
			Items: []string{"WVU", "LJ Aviation"},
		}

		var err error
		_, pType, err = prompt.Run()

		if pType == "LJ Aviation" {
			pType = flights.PersonTypeLJA
		} else if pType == "WVU" {
			pType = flights.PersonTypeWVU
		}
		if err != nil {
			return err
		}
	}

	person := flights.Person{
		Name: name,
		Type: pType,
	}

	if err := db.Create(&person).Error; err != nil {
		return err
	}

	log.Printf("person \"%s\" added", person.Name)

	return nil
}

func personLookup(c *cli.Context) error {
	name := c.String("name")
	if name == "" {
		prompt := promptui.Prompt{
			Label: "Person Name",
		}
		var err error
		name, err = prompt.Run()
		if err != nil {
			return err
		}
	}

	people, err := personSearchByName(name)
	if err != nil {
		return err
	}

	if len(people) == 0 {
		log.Printf("no people found with name \"%s\"", name)
		return nil
	}

	for _, person := range people {
		log.Printf("person: %s\ntype: %s", person.Name, person.Type)
	}

	return nil
}

func personSearchByName(name string) ([]flights.Person, error) {
	var people []flights.Person

	if err := db.Where("name LIKE ?", "%"+name+"%").Find(&people).Error; err != nil {
		return nil, err
	}

	return people, nil
}
