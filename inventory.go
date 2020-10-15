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

// ByName implements sort.Interface for []Inventory
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview
type ByNameIv []InventoryItem

func (n ByNameIv) Len() int           { return len(n) }
func (n ByNameIv) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByNameIv) Less(i, j int) bool { return n[i].Name < n[j].Name }

// ByID implements sort.Interface for []Inventory
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByIDIv []InventoryItem

func (d ByIDIv) Len() int           { return len(d) }
func (d ByIDIv) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByIDIv) Less(i, j int) bool { return d[i].ID < d[j].ID }
