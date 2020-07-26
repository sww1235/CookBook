# Database Schema Documentation

want to have separate ingredient and inventory tables as they can have different
names or purchase quantities. IE, ingredient would be 1lb flour and inventory
would be 5lb bag flour, potentially with a partial remaining quantity.

TODO: Need to evaluate all VARCHARs for potential changes to text and any
integers for changes to floats.

## recipe

store information about a specific recipe. Multiple versions of a recipe are allowed

| Column Name       | Datatype (mysql) | Datatyle (sqlite) | Description                                                                           |
| ----------------- | ---------------- | ----------------- | ------------------------------------------------------------------------------------- |
| ID                | int (pk)         | INTEGER (pk)      | unique ID for each recipe                                                             |
| Name              | text             | TEXT              | Name of recipe                                                                        |
| Description       | text             | TEXT              | Short Description of recipe.                                                          |
| Comments          | text             | TEXT              | comments on recipe or history of recipe                                               |
| Source            | text             | TEXT              | source of recipe, include URL or other info                                           |
| Author            | text             | TEXT              | name of original creator of specific recipe (if known)                                |
| QuantityMade      | decimal(7,2)     | NUM               | a specific quantity that this recipe makes. Allows for easy doubling or meal planning |
| QuantityMadeUnits | INT (fk)         | INTEGER (fk)      | a foreign key linked to the units table to select a unit of measure for QuantityMade  |


## ingredient

stores ingredients as used in recipes.

| Column Name   | Datatype (mysql) | Datatype (sqlite) | Description                                |
| ------------- | ---------------- | ----------------- | ------------------------------------------ |
| ID            | int (pk)         | INTEGER (pk)      | unique ID for each ingredient              |
| Name          | VARCHAR          | TEXT              | name of ingredient                         |
| InventoryID   | integer (fk)     | INTEGER (fk)      | ingredient to its precursor inventory item |
| Quantity      | decimal(7,2)     | NUM               | quantity of ingredient used in recipe      |
| QuantityUnits | int (fk)         | INTEGER (fk)      | units of ingredient used in recipe         |

## ingredient\_inventory

maps ingredients to inventory items

has composite primary key

| Column Name  | Datatype (mysql) | Datatype (sqlite) | Description                       |
| ------------ | ---------------- | ----------------- | --------------------------------- |
| IngredientID | int (pk, fk)     | INTEGER (pk, fk)  | unique ID for each ingredient     |
| InventoryID  | integer (pk, fk) | INTEGER (pk, fk)  | unique ID for each inventory item |


## ingredient\_recipe

maps ingredients or sub recipes to recipes

| Column Name  | Datatype (mysql) | Datatype (sqlite) | Description                   |
| ------------ | ---------------- | ------------------|------------------------------ |
| RecipeID     | int (pk, fk)     | INTEGER (pk, fk)  | unique ID for each recipe     |
| IngredientID | int (pk, fk)     | INTEGER (pk, fk)  | unique ID for each ingredient |



## step

TODO: Somehow figure out how to reference specific ingredients in a step.
Probably done using encoding in the step instructions field of the IngredientID

| Column Name      | Datatype (mysql) | Datatype (sqlite) | Description                                        |
| ---------------- | ---------------- | ----------------- | -------------------------------------------------- |
| ID               | integer (pk)     | INTEGER (pk)      | unique ID for step                                 |
| Instructions     | text             | TEXT              | instructions for step                              |
| Time             | decimal(7,2)     | NUM               | stores time of step in seconds.                    |
| StepTypeID       | int (fk)         | INTEGER (fk)      | foreign key mapping to step type description table |
| temperature      | int              | INTEGER           | cooking temperature of step                        |
| temperatureUnits | int (fk)         | INTEGER (fk)      | mapping to units table                             |


## stepType

initial step types

-	prep
-	cook
-	wait
-	other

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description            |
| ----------- | ---------------- | ----------------- | ---------------------- |
| ID          | int (pk)         | INTEGER (pk)      | unique ID for stepType |
| Name        | VARCHAR(20)      | TEXT              | name of step type      |

## step\_recipe

maps steps to recipes

has composite primary key

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description                   |
| ----------- | ---------------- | ----------------- | ----------------------------- |
| RecipeID    | int (pk, fk)     | INTEGER (pk, fk)  | unique ID for each ingredient |
| StepID      | int (pk, fk)     | INTEGER (pk, fk)  | unique ID for each step       |

## inventory

| Column Name     | Datatype (mysql) | Datatype (sqlite) | Description                     |
| --------------- | ---------------- | ----------------- | ------------------------------- |
| ID              | int (pk          | INTEGER (pk)      | unique ID for inventory item    |
| EAN             | char(14)         | TEXT              | barcode data                    |
| Name            | text             | TEXT              | short name of inventory item    |
| Description     | text             | TEXT              | description of inventory item   |
| storedQty       | decimal(7,2)     | NUM               | quantity of item in inventory   |
| PackageQty      | decimal(7        | NUM               | quantity of item in package     |
| PackageQtyUnits | int (fk          | INTEGER (fk)      | fk for referencing unit table   |

## units

stores all units with a standardized PK and a human readable description

<unitsofmeasure.org/ucum.html>

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description                                     |
| ----------- | ---------------- | ----------------- | ----------------------------------------------- |
| ID          | int (pk)         | INTEGER (pk)      | unique ID for unit (follow ucum standard above) |
| Name        | text             | TEXT              | print name of unit                              |
| Description | text             | TEXT              | description of unit                             |

## tags


| Column Name | Datatype (mysql) | Datatype (sqlite) | Description        |
| ----------- | ---------------- | ----------------- | ------------------ |
| ID          | int (pk          | INTEGER (pk)      | unique ID for tag  |
| Name        | text             | TEXT              | name of tag        |
| Description | text             | TEXT              | description of tag |

## tag\_recipe

has composite primary key

maps tags to recipes

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description               |
| ----------- | ---------------- | ----------------- | ------------------------- |
| tagID       | int (pk, fk)     | INTEGER (pk, fk)  | unique ID for tag         |
| recipeID    | int (pk, fk)     | INTEGER (pk, fk)  | unique ID for recipe      |



