package retriever

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func RetrieveDoc(link string, gid int) []string {
	composedLink := fmt.Sprintf("%s/export?exportFormat=csv&gid=%d", link, gid)
	response, err := http.Get(composedLink)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(b)

	return strings.Split(bodyString, "\r\n")
}

func ExtractDishFromPlan(dishes []string) []string {
	result := make([]string, 0)
	for _, row := range dishes[1:] {
		for _, dish := range strings.Split(row, ",")[1:] {
			if dish != "" {
				result = append(result, dish)
			}

		}

	}
	return result
}

func BuildRecipes(recipes []string) map[string][]string {
	result := make(map[string][]string)

	for _, recipe := range recipes[1:] {

		sanitizedRecipe := strings.ReplaceAll(strings.ReplaceAll(recipe, "\"", ""), "\\", "")
		recipeFields := strings.Split(sanitizedRecipe, ",")
		recipeName := recipeFields[0]
		result[recipeName] = mapTrim(recipeFields[2:], strings.Trim)
	}
	return result
}

func mapTrim(vs []string, f func(string, string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v, " ")
	}
	return vsm
}
