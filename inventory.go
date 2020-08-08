package main

import "fmt"

type InventoryItem struct {
	ID           int     // id in database
	EAN          string  // UPC/EAN barcode
	Name         string  // name of inventory item
	Description  string  // description of inventory item
	Quantity     float64 // quantity of inventory item stored
	QuantityUnit Unit    // unit of measure of inventory item

}

func (i *InventoryItem) String() string {

	return fmt.Sprintf("Inventory contains %.2f %s of %s.", i.Quantity, i.QuantityUnit.Name, i.Name)

}
