package httpfinger

type Database struct {
	database []*FingerPrint
}

func (d *Database) Add(s string) error {
	httpFinger, err := NewFingerPrint(s)
	if err != nil {
		return err
	}
	d.database = append(d.database, httpFinger)
	return nil
}

func (d *Database) Search(banner *Banner) []string {
	var products []string
	for _, fingerPrint := range d.database {
		if productName := fingerPrint.Match(banner); productName != "" {
			products = append(products, productName)
		}
	}
	return products
}
