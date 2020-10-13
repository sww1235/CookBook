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
type ByName []Inventory

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j] }

// ByID implements sort.Interface for []Inventory
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByID []Inventory

func (d ByID) Len() int           { return len(d) }
func (d ByID) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByID) Less(i, j int) bool { return d[i].ID < d[j].ID }
