package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Recipe struct {
	Id           int      `json:"id"`
	Title        string   `json:"title"`
	Tags         []string `json:"tags"`
	Ingredients  []string `json:"ingredients"`
	Amount       []string `json:"amount"`
	Instructions string   `json:"instructions"`
	Likes        int      `json:"likes"`
	CreatorName  string   `json:"creatorName"`
}

type Ingredient struct {
	Id         int    `json:"id"`
	RecipeId   int    `json:"recipe_id"`
	Ingredient string `json:"ingredient"`
	Amount     string `json:"amount"`
}

type Tag struct {
	Id       int    `json:"id"`
	RecipeId int    `json:"recipe_id"`
	Tag      string `json:"tag"`
}

func initCookBookController() {
	fmt.Println("Initializing CookBookController...")

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/writeRecipe", writeRecipe)
	mux.HandleFunc("/getRecipe", getRecipe)
	mux.HandleFunc("/addLike", AddLike)
	mux.HandleFunc("/removeLike", RemoveLike)
	mux.HandleFunc("/getAllRecipes", getAllRecipes)
	mux.HandleFunc("/deleteRecipeById", deleteRecipeById)

	fmt.Println("CookBookController initialized. Listening on port 8085...")
	err := http.ListenAndServe(":"+BACKEND_PORT, corsMiddleware(mux))
	if err != nil {
		panic(err)
	}

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Check if the request method is OPTIONS
		if r.Method == http.MethodOptions {
			// If it is, respond with a 200 OK status and return
			w.WriteHeader(http.StatusOK)
			return
		}

		// Otherwise, call the original handler's ServeHTTP method
		next.ServeHTTP(w, r)
	})
}

func getAllRecipes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var recipeList []Recipe
	var err error

	var queryString = "SELECT * FROM recipes"
	var rows, _ = db.Query(queryString)

	for rows.Next() {
		var recipe Recipe
		err = rows.Scan(&recipe.Id, &recipe.Title, &recipe.Instructions, &recipe.Likes, &recipe.CreatorName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		recipeList = append(recipeList, recipe)
	}

	recipeList, err = collectTagsAndIngredients(recipeList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(recipeList)
	if err != nil {
		return
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pong")
	_, err := w.Write([]byte("pong"))
	if err != nil {
		panic(err)
	}
}

func writeRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Writing recipe...")

	var recipe Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var queryString = "INSERT INTO recipes (title, instructions, likes, creatorName) VALUES (?, ?, ?, ?)"

	res, err := db.Exec(queryString, recipe.Title, recipe.Instructions, recipe.Likes, recipe.CreatorName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recipeId, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range recipe.Ingredients {
		queryString = "INSERT INTO ingredients (recipe_id, ingredient, amount) VALUES (?, ?, ?)"
		_, err = db.Exec(queryString, recipeId, recipe.Ingredients[i], recipe.Amount[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for i := range recipe.Tags {
		queryString = "INSERT INTO tags (recipe_id, tag) VALUES (?, ?)"
		_, err = db.Exec(queryString, recipeId, recipe.Tags[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

/**
 * getRecipe is a function that retrieves a recipe from the database.
 * It can search for a recipe by ID, title, or tags.
 */
func getRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var recipe Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var recipeList []Recipe

	//search for id
	if recipe.Id != 0 { // search for id
		recipeList, err = getRecipesById(recipe.Id)
	} else if recipe.Title != "" { // search for title
		recipeList, err = getRecipesByTitle(recipe.Title)
	} else if recipe.Tags != nil { // search for tags
		recipeList, err = getRecipeByTags(recipe.Tags)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Getting recipe...")

	err = json.NewEncoder(w).Encode(recipeList)
	if err != nil {
		return
	}
}

func AddLike(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var recipe Recipe
	err = json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		panic(err)
	}

	var queryString = "UPDATE recipes SET likes = likes + 1 WHERE id = ?"
	_, err = db.Exec(queryString, recipe.Id)
	if err != nil {
		panic(err)
	}
}

func RemoveLike(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var recipe Recipe
	err = json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		panic(err)
	}

	var queryString = "UPDATE recipes SET likes = likes - 1 WHERE id = ?"
	_, err = db.Exec(queryString, recipe.Id)
	if err != nil {
		panic(err)
	}

}

func deleteRecipeById(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var recipe Recipe
	err = json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		panic(err)
	}

	queryStrings := []string{"DELETE FROM tags WHERE recipe_id = ?", "DELETE FROM ingredients WHERE recipe_id = ?", "DELETE FROM recipes WHERE id = ?"}

	for _, queryString := range queryStrings {
		_, err := db.Exec(queryString, recipe.Id)
		if err != nil {
			panic(err)
		}
	}
}

func getRecipesById(id int) ([]Recipe, error) {
	var recipeList []Recipe
	var err error

	var queryString = "SELECT * FROM recipes WHERE id = ?"
	var row = db.QueryRow(queryString, id)

	var recipe Recipe
	err = row.Scan(&recipe.Id, &recipe.Title, &recipe.Instructions, &recipe.Likes, &recipe.CreatorName)
	if err != nil {
		return nil, err
	}

	recipeList = append(recipeList, recipe)

	recipeList, err = collectTagsAndIngredients(recipeList)

	if err != nil {
		return nil, err
	}

	return recipeList, nil
}

func getRecipesByTitle(title string) ([]Recipe, error) {
	var recipeList []Recipe
	var err error

	var queryString = "SELECT * FROM recipes WHERE title LIKE ?"
	var rows, _ = db.Query(queryString, "%"+title+"%")

	for rows.Next() {
		var recipe Recipe
		err = rows.Scan(&recipe.Id, &recipe.Title, &recipe.Instructions, &recipe.Likes, &recipe.CreatorName)
		if err != nil {
			return nil, err
		}

		recipeList = append(recipeList, recipe)
	}

	recipeList, err = collectTagsAndIngredients(recipeList)

	if err != nil {
		return nil, err
	}

	return recipeList, nil
}

func collectTagsAndIngredients(recipeList []Recipe) ([]Recipe, error) {
	var err error

	for i := range recipeList {
		// Get the tags for each recipe
		var queryString = "SELECT * FROM tags WHERE recipe_id = ?"
		var rows, _ = db.Query(queryString, recipeList[i].Id)
		for rows.Next() {
			var tag Tag
			err = rows.Scan(&tag.Id, &tag.RecipeId, &tag.Tag)
			if err != nil {
				return nil, err
			}
			recipeList[i].Tags = append(recipeList[i].Tags, tag.Tag)
		}

		// Get the ingredients for each recipe
		queryString = "SELECT * FROM ingredients WHERE recipe_id = ?"
		rows, _ = db.Query(queryString, recipeList[i].Id)
		for rows.Next() {
			var ingredient Ingredient
			err = rows.Scan(&ingredient.Id, &ingredient.RecipeId, &ingredient.Ingredient, &ingredient.Amount)
			if err != nil {
				return nil, err
			}
			recipeList[i].Ingredients = append(recipeList[i].Ingredients, ingredient.Ingredient)
			recipeList[i].Amount = append(recipeList[i].Amount, ingredient.Amount)
		}
	}
	return recipeList, nil
}

func getTags(tags []string) ([]Tag, error) {
	// Convert the tags slice to a comma-separated string
	tagsString := strings.Join(tags, "','")

	// Prepare the SQL query
	queryString := fmt.Sprintf(`SELECT * FROM tags WHERE tag IN ('%s')`, tagsString)

	// Execute the query
	rows, err := db.Query(queryString)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	// Scan the results into a Tag slice
	var resultTags []Tag
	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.Id, &tag.RecipeId, &tag.Tag)
		if err != nil {
			return nil, err
		}
		resultTags = append(resultTags, tag)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resultTags, nil
}

func getRecipeByTags(tags []string) ([]Recipe, error) {
	var recipeList []Recipe
	var err error

	var foundTags, _ = getTags(tags)
	var sqlString = "SELECT * FROM tags WHERE tag = ?"

	var resultTags []Tag
	for i := range tags {
		var recipeId = foundTags[i].Tag
		var rows, _ = db.Query(sqlString, recipeId)

		for rows.Next() {
			var tag Tag
			err = rows.Scan(&tag.Id, &tag.RecipeId, &tag.Tag)
			if err != nil {
				return nil, err
			}
			resultTags = append(resultTags, tag)
		}
	}

	for i := range resultTags {
		var recipe Recipe
		var queryString = "SELECT * FROM recipes WHERE id = ?"
		var row = db.QueryRow(queryString, resultTags[i].RecipeId)
		err = row.Scan(&recipe.Id, &recipe.Title, &recipe.Instructions, &recipe.Likes, &recipe.CreatorName)
		if err != nil {
			return nil, err
		}

		var entryFound = false
		//Check if the recipe is already in the list
		for j := range recipeList {
			if recipe.Id == recipeList[j].Id {
				entryFound = true
				break
			}
		}

		if !entryFound {
			recipeList = append(recipeList, recipe)
		}
	}

	recipeList, err = collectTagsAndIngredients(recipeList)
	if err != nil {
		return nil, err
	}

	return recipeList, nil
}
