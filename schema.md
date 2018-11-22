# Database Schema Documentation

want to have separate ingredient and inventory tables as they can have different
names or purchase quantities. IE, ingredient would be 1lb flour and inventory
would be 5lb bag flour, potentially with a partial remaining quantity.

TODO: Need to evaluate all VARCHARs for potential changes to text

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
| QuantityMadeUnits | integer | a foreign key linked to the units table to select a unit of measure for QuantityMade |

## Ingredient

stores ingredients as used in recipes.

| Column Name | Datatype | Description |
| ----------- | -------- | ----------- |
| IngredientID | integer (autoincrement) | unique ID for each ingredient |
| IngredientName | VARCHAR | name of ingredient |
| InventoryID | integer | foreign key mapping ingredient to its precursor inventory item |


## Ingredient <-> Recipe

maps ingredients or sub recipes to recipes

| Column Name | Datatype | Description |
| ----------- | -------- | ----------- |
| IngRecMapID | integer (autoincrement) | unique ID for each ingredient <-> recipe mapping |
| RecipeID | integer | unique ID for each ingredient (fk) |
| IngredientID | integer | unique ID for each ingredient (fk) |
| IngOrRec | integer or boolean | whether mapping is to another recipe or is to an ingredient |
| IngredientQuantity | integer | how much of ingredient is needed for recipe |
| IngredientQuantityUnits | integer | foreign key mapping to units table to select a unit for IngredientQuantity |



## Step

## Inventory


## Units

stores all units with a standardized PK and a human readable description
