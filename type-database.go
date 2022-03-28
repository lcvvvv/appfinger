package main

import (
	"errors"
	"httpfinger/product"
	"strings"
)

type Database struct {
	database []*FingerPrint
}

var Data *Database

func (d *Database) Add(s string) error {
	httpFinger, err := NewFingerPrint(s)
	if err != nil {
		return err
	}
	d.database = append(d.database, httpFinger)
	return nil
}

func (d *Database) Search(banner *Banner) []*product.Product {
	var products []*product.Product
	for _, fingerPrint := range d.database {
		if p := fingerPrint.Match(banner); p != nil {
			products = append(products, p)
		}
	}
	return products
}

func makeDatabase(source string) error {
	Data = &Database{database: []*FingerPrint{}}

	for _, line := range strings.Split(source, "\n") {
		err := Data.Add(line)
		if err != nil {
			return errors.New(err.Error() + line)
		}
	}
	return nil
}
