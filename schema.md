# Database Schema Documentation

want to have separate ingredient and inventory tables as they can have different
names or purchase quantities. IE, ingredient would be 1lb flour and inventory
would be 5lb bag flour, potentially with a partial remaining quantity.

## recipes

store information about a specific recipe. Multiple versions of a recipe are allowed

| Column Name       | Datatype (mysql) | Datatyle (sqlite) | Description                                            |
| ----------------- | ---------------- | ----------------- | ------------------------------------------------------ |
| ID                | int (pk)         | INTEGER (pk)      | unique ID for each recipe                              |
| Name              | text             | TEXT              | Name of recipe                                         |
| Description       | text             | TEXT              | Short Description of recipe.                           |
| Comments          | text             | TEXT              | comments on recipe or history of recipe                |
| Source            | text             | TEXT              | source of recipe, include URL or other info            |
| Author            | text             | TEXT              | name of original creator of specific recipe (if known) |
| Quantity          | decimal(7,2)     | NUM               | a specific quantity that this recipe makes.            |
| QuantityMadeUnits | int (fk)         | INTEGER (fk)      | unit of measure for QuantityMade                       |
| initialVersion    | int (fk)         | INTEGER (fk)      | id of initial version of recipe. null if no revs       |
| version           | int              | INTEGER           | revision number of recipe                              |

## ingredients

stores ingredients as used in recipes.

| Column Name   | Datatype (mysql) | Datatype (sqlite) | Description                                |
| ------------- | ---------------- | ----------------- | ------------------------------------------ |
| ID            | int (pk)         | INTEGER (pk)      | unique ID for each ingredient              |
| Name          | VARCHAR          | TEXT              | name of ingredient                         |
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



## steps

TODO: Somehow figure out how to reference specific ingredients in a step.
Probably done using encoding in the step instructions field of the IngredientID

| Column Name  | Datatype (mysql) | Datatype (sqlite) | Description                                        |
| ------------ | ---------------- | ----------------- | -------------------------------------------------- |
| ID           | integer (pk)     | INTEGER (pk)      | unique ID for step                                 |
| Instructions | text             | TEXT              | instructions for step                              |
| Time         | decimal(7,2)     | NUM               | stores time of step in seconds.                    |
| StepType     | int (fk)         | INTEGER (fk)      | mapping to either stepType table in database, or stepType iota in code |
| temperature  | decimal(8,3)     | NUM               | cooking temperature of step                        |
| tempUnit     | int (fk)         | INTEGER (fk)      | mapping to units table                             |


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

| Column Name  | Datatype (mysql) | Datatype (sqlite) | Description                   |
| ------------ | ---------------- | ----------------- | ----------------------------- |
| ID           | int (pk)         | INTEGER (pk)      | unique ID for inventory item  |
| EAN          | char(14)         | TEXT              | barcode data                  |
| Name         | text             | TEXT              | short name of inventory item  |
| Description  | text             | TEXT              | description of inventory item |
| Quantity     | decimal(7,2)     | NUM               | quantity of item in inventory |
| QuantityUnit | int (fk)         | INTEGER (fk)      | fk for referencing unit table |

## units

stores all units with a standardized PK and a human readable description

<unitsofmeasure.org/ucum.html>

| Column Name   | Datatype (mysql) | Datatype (sqlite) | Description                                   |
| ------------- | ---------------- | ----------------- | --------------------------------------------- |
| ID            | int (pk)         | INTEGER (pk)      | unique ID for unit                            |
| Name          | text             | TEXT              | print name of unit                            |
| Description   | text             | TEXT              | description of unit                           |
| Symbol        | text             | TEXT              | unit symbol                                   |
| isCustom      | bool             | INTEGER           | is unit custom or standard                    |
| refIngredient | int (fk)         | INGEGER (fk)      | fk of ingredient for ingredient specific unit |
| unitType      | int (fk)         | INTEGER (fk)      | base type of unit                             |


## unitType

base unit types

-	time
-	length
-	mass
-	current
-	temperature
-	quantity
-	lum\_intensity

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description          |
| ----------- | ---------------- | ----------------- | -------------------- |
| ID          | int (pk)         | INTEGER (pk)      | unique id            |
| Name        | text             | TEXT              | base type name       |

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


## lastMade

records when a recipe was made. This is marked manually by the chef

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description          |
| ----------- | ---------------- | ----------------- | -------------------- |
| ID          | int (pk)         | INTEGER (pk)      | unique id            |
| recipe      | int (fk)         | INTEGER (fk)      | related recipe       |
| dateMade    | date             | TEXT              | date recipe was made |
| notes       | text             | TEXT              | notes from cooking   |

## equipment

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description        |
| ----------- | ---------------- | ----------------- | ------------------ |
| ID          | int (pk)         | INTEGER (pk)      | unique id          |
| Name        | text             | TEXT              | name of equipment  |
| isOwned     | bool             | NUM               | is equipment owned |

## equipment\_recipe

has composite primary key

maps equipment to recipes

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description             |
| ----------- | ---------------- | ----------------- | ----------------------- |
| equipmentID | int (pk, fk)     | INTEGER (pk, fk)  | unique ID for equipment |
| recipeID    | int (pk, fk)     | INTEGER (pk, fk)  | unique ID for recipe    |



## unitConversions

<https://dba.stackexchange.com/a/62468/185504>

`toUnitValue = ((fromUnitOffset + fromOffset) * multiplicand / denominator) + toOffset`

| Column Name  | Datatype (mysql) | Datatype (sqlite) | Description                |
| ------------ | ---------------- | ----------------- | -------------------------- |
| ID           | int (pk)         | INTEGER (pk)      | unique id                  |
| fromUnit     | int (fk)         | INTEGER (fk)      | id of unit being converted |
| toUnit       | int (fk)         | INTEGER (fk)      | id of unit to convert to   |
| multiplicand | decimal(12,3)    | NUM               | multiplicand               |
| denominator  | decimal(12,3)    | NUM               | denominator                |
| fromOffset   | decimal(12,3)    | NUM               | from unit offset           |
| toOffset     | decimal(12,3)    | NUM               | to unit offset             |




