# Database Schema Documentation

want to have separate ingredient and inventory tables as they can have different
names or purchase quantities. IE, ingredient would be 1lb flour and inventory
would be 5lb bag flour, potentially with a partial remaining quantity.

TODO: Need to evaluate all VARCHARs for potential changes to text and any
integers for changes to floats.

## Recipe

store information about a specific recipe. Multiple versions of a recipe are allowed

| Column Name | Datatype | Description |
| ----------- | -------- | ----------- |
| RecipeID | integer (autoincrement) | unique ID for each recipe |
| RecipeName | VARCHAR | Name of recipe |
| RecipeDescription | VARCHAR | Short Description of recipe. |
| RecipeComments | VARCHAR | comments on recipe or history of recipe |
| RecipeSource | VARCHAR | source of recipe, include URL or other info |
| RecipeAuthor | VARCHAR | name of original creator of specific recipe (if known) |
| QuantityMade | integer | a specific quantity that this recipe makes. Allows for easy doubling or meal planning |
| QuantityMadeUnits | VARCHAR | a foreign key linked to the units table to select a unit of measure for QuantityMade |


## Ingredient

stores ingredients as used in recipes.

| Column Name | Datatype | Description |
| ----------- | -------- | ----------- |
| IngredientID | integer (autoincrement) | unique ID for each ingredient |
| IngredientName | VARCHAR | name of ingredient |
| InventoryID | integer | foreign key mapping ingredient to its precursor inventory item |

## Ingredient <-> Inventory

maps ingredients to inventory items

has composite primary key

| Column Name | Datatype | Description |
| ----------- | -------- | ----------- |
| IngredientID | integer | unique ID for each ingredient (pk,fk) |
| InventoryID | integer | unique ID for each inventory item (pk,fk) |


## Ingredient <-> Recipe

maps ingredients or sub recipes to recipes

| Column Name | Datatype | Description |
| ----------- | -------- | ----------- |
| IngRecMapID | integer (autoincrement) | unique ID for each ingredient <-> recipe mapping |
| RecipeID | integer | unique ID for each recipe (fk) |
| IngredientID | integer | unique ID for each ingredient (fk) |
| IngOrRec | integer or boolean | whether mapping is to another recipe or is to an ingredient |
| IngredientQuantity | integer | how much of ingredient is needed for recipe |
| IngredientQuantityUnits | VARCHAR | foreign key mapping to units table to select a unit for IngredientQuantity |



## Step

TODO: Somehow figure out how to reference specific ingredients in a step.
Probably done using encoding in the step instructions field of the IngredientID

| Column Name | Datatype | Description |
| ----------- | -------- | ----------- |
| StepID | integer (autoincrement) | unique ID for step |
| StepInstructions | VARCHAR | instructions for step |
| StepTime | integer | stores time of step in seconds. Will be displayed as other units |
| StepTypeID | integer | foreign key mapping to step type description table (fake enum) |
| temperature | integer | cooking temperature of step |
| temperatureUnits | VARCHAR | foreignKey mapping to units table |


## StepType

| Column Name | Datatype | Description |
| ----------- | -------- | ----------- |
| StepTypeID | integer (autoincrement) | unique ID for stepType |
| StepName | VARCHAR | name of step type (prep, cook, wait, other) |

## Step <-> Recipe

maps steps to recipes

has composite primary key

| Column Name | Datatype | Description |
| ----------- | -------- | ----------- |
| RecipeID | integer | unique ID for each ingredient (fk) |
| StepID | integer | unique ID for each step (fk) |

## Inventory

| Column Name     | Datatype (mysql) | Datatype (sqlite) | Description                     |
| --------------- | ---------------- | ----------------- | ------------------------------- |
| ID              | integer          | INTEGER           | unique ID for inventory item    |
| EAN             | CHAR(14)         | TEXT              | barcode data                    |
| Name            | VARCHAR          | TEXT              | short name of inventory item    |
| Description     | VARCHAR          | TEXT              | description of inventory item   |
| storedQty       | integer          | INTEGER           | quantity of item in inventory   |
| PackageQty      | integer          | INTEGER           | quantity of item in package     |
| PackageQtyUnits | integer          | INTEGER           | fk for referencing unit table   |

## Units

stores all units with a standardized PK and a human readable description

<unitsofmeasure.org/ucum.html>

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description                                     |
| ----------- | ---------------- | ----------------- | ----------------------------------------------- |
| ID          | int              | INTEGER           | unique ID for unit (follow ucum standard above) |
| Name        | VARCHAR          | TEXT              | print name of unit                              |
| Description | VARCHAR          | TEXT              | description of unit                             |

## Tags


| Column Name | Datatype (mysql) | Datatype (sqlite) | Description        |
| ----------- | ---------------- | ----------------- | ------------------ |
| ID          | int              | INTEGER           | unique ID for tag  |
| Name        | VARCHAR          | TEXT              | name of tag        |
| Description | VARCHAR          | TEXT              | description of tag |

## Tag <-> Recipe (tag\_recipe)

has composite primary key

maps tags to recipes

| Column Name | Datatype (mysql) | Datatype (sqlite) | Description               |
| ----------- | ---------------- | ----------------- | ------------------------- |
| tagID       | int              | INTEGER           | unique ID for tag (fk)    |
| recipeID    | int              | INTEGER           | unique ID for recipe (fk) |



