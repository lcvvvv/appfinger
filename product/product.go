package product

type Product struct {
	ProductName string
	CompanyName string
	Description string
	Category    string
	Subcategory string
	Version     string
	Level       int
}

const (
	OperatingSystem = 0x00000a1
	Protocol        = 0x00000b2
	Service         = 0x00000c3
	Application     = 0x00000d4
	Component       = 0x00000e5
)

func New(productName, companyName, description, category, subcategory, version string, level int) *Product {
	return &Product{
		ProductName: productName,
		CompanyName: companyName,
		Description: description,
		Category:    category,
		Subcategory: subcategory,
		Version:     version,
		Level:       level,
	}
}
