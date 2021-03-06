package gql

import (
	"fmt"
	"github.com/aPruner/my-fridge-server/app/db"
	"github.com/graphql-go/graphql"
	"log"
)

type Resolver struct {
	database *db.Db
}

func (r *Resolver) UserQueryResolver(p graphql.ResolveParams) (interface{}, error) {
	// Type-check the name
	username, ok := p.Args["username"].(string)
	if ok {
		users, err := r.database.GetUsersByUsername(username)
		if err != nil {
			return nil, err
		}
		return users, nil
	}
	err := fmt.Errorf("type-checking error: username was not a string")
	log.Print(err)
	return nil, err
}

func (r *Resolver) FoodItemsQueryResolver(p graphql.ResolveParams) (interface{}, error) {
	// Type-check the potential arguments
	foodItemsInput, ok := p.Args["foodItemsQueryInput"].(map[string]interface{})
	var err error
	if ok {
		householdId, householdIdOk := foodItemsInput["householdId"].(int)
		shoppingListId, shoppingListIdOk := foodItemsInput["shoppingListId"].(int)
		if householdIdOk && shoppingListIdOk {
			err = fmt.Errorf("type-checking error: please provide only one of householdId or shoppingListId")
			log.Print(err)
			return nil, err
		} else if householdIdOk {
			foodItems, err := r.database.GetFoodItemsByHouseholdId(householdId)
			if err != nil {
				return nil, err
			}
			return foodItems, nil
		} else if shoppingListIdOk {
			foodItems, err := r.database.GetFoodItemsByShoppingListId(shoppingListId)
			if err != nil {
				return nil, err
			}
			return foodItems, nil
		}
	}
	err = fmt.Errorf("type-checking error: please provide a valid foodItemsInput object (one of householdId and shoppingListId non-null)")
	log.Print(err)
	return nil, err
}

func (r *Resolver) HouseholdIdQueryResolver(p graphql.ResolveParams) (interface{}, error) {
	// Type-check the userId
	userId, ok := p.Args["userId"].(int)
	if ok {
		householdId, err := r.database.GetHouseholdIdByUserId(userId)
		if err != nil {
			return nil, err
		}
		return householdId, nil
	}
	err := fmt.Errorf("type-checking error: userId was not an int")
	log.Print(err)
	return nil, err
}

// TODO: Implement this
func (r *Resolver) HouseholdQueryResolver(p graphql.ResolveParams) (interface{}, error) {
	err := fmt.Errorf("type-checking error: params were invalid")
	return nil, err
}

func (r *Resolver) ShoppingListQueryResolver(p graphql.ResolveParams) (interface{}, error) {
	householdId, ok := p.Args["householdId"].(int)
	if ok {
		shoppingLists, err := r.database.GetShoppingListsByHouseholdId(householdId)
		if err != nil {
			return nil, err
		}
		return shoppingLists, nil
	}

	err := fmt.Errorf("type-checking error: householdId was not an int")
	log.Print(err)
	return nil, err
}

func (r *Resolver) CreateFoodItemMutationResolver(p graphql.ResolveParams) (interface{}, error) {
	// TODO: Figure out optional params here
	name, nameOk := p.Args["name"].(string)
	category, categoryOk := p.Args["category"].(string)
	amount, amountOk := p.Args["amount"].(int)
	unit, unitOk := p.Args["unit"].(string)
	householdId, householdIdOk := p.Args["householdId"].(int)
	shoppingListId, shoppingListIdOk := p.Args["shoppingListId"].(int)

	if nameOk && categoryOk && amountOk && unitOk && householdIdOk && shoppingListIdOk {
		newFoodItemId, err := r.database.CreateFoodItem(name, category, amount, unit, householdId, shoppingListId)
		if err != nil {
			return nil, err
		}
		return newFoodItemId, nil
	}
	err := fmt.Errorf("type-checking error: a combination of name, category, amount, householdId, and shoppingListId was misformed")
	log.Print(err)
	// TODO: In these cases, the server should probably throw a 400 bad request
	return nil, err
}

func (r *Resolver) UpdateFoodItemMutationResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	var err error
	if ok {
		err = r.database.UpdateFoodItem(id, p)
		if err != nil {
			return nil, err
		}
		// TODO: Maybe should return the updated FoodItem instead?
		return nil, nil
	}
	err = fmt.Errorf("type-checking error: id was not an int")
	log.Print(err)
	return nil, err
}

func (r *Resolver) DeleteFoodItemMutationResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	var err error
	if ok {
		err = r.database.DeleteFoodItem(id)
		if err != nil {
			return nil, err
		}
		// TODO: Proper graphql convention is to return the deleted row (FoodItem)
		return nil, nil
	}
	err = fmt.Errorf("type-checking error: id was not an int")
	log.Print(err)
	return nil, err
}

func (r *Resolver) CreateShoppingListResolver(p graphql.ResolveParams) (interface{}, error) {
	// TODO: Figure out optional params here
	name, nameOk := p.Args["name"].(string)
	description, descriptionOk := p.Args["description"].(string)
	userId, userIdOk := p.Args["userId"].(int)
	householdId, householdIdOk := p.Args["householdId"].(int)

	if nameOk && descriptionOk && householdIdOk && userIdOk {
		newFoodItemId, err := r.database.CreateShoppingList(name, description, userId, householdId)
		if err != nil {
			return nil, err
		}
		return newFoodItemId, nil
	}
	// TODO: In these cases, the server should probably throw a 400 bad request
	err := fmt.Errorf("type-checking error: a combination of name, description, category, amount, and householdId was misformed")
	log.Print(err)
	return nil, err
}

func (r *Resolver) UpdateShoppingListResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	var err error
	if ok {
		err = r.database.UpdateShoppingList(id, p)
		if err != nil {
			return nil, err
		}
		// TODO: Maybe should return the updated ShoppingList instead?
		return nil, nil
	}
	err = fmt.Errorf("type-checking error: id was not an int")
	log.Print(err)
	return nil, err
}

func (r *Resolver) DeleteShoppingListResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	var err error
	if ok {
		err = r.database.DeleteShoppingList(id)
		if err != nil {
			return nil, err
		}
		// TODO: Proper graphql convention is to return the deleted row (ShoppingList)
		return nil, nil
	}
	err = fmt.Errorf("type-checking error: id was not an int")
	log.Print(err)
	return nil, err
}
