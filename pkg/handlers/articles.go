package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/namanag0502/go-blog/pkg/models"
	"github.com/namanag0502/go-blog/pkg/utils"
)

func GetArticles() *[]models.Article {
	f, err := os.Open("pkg/data/articles.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer f.Close()
	var articles []models.Article
	json.NewDecoder(f).Decode(&articles)
	return &articles
}

func Home(w http.ResponseWriter, r *http.Request) {
	articles := GetArticles()
	if articles == nil {
		utils.WriteJsonError(w, "No articles found", http.StatusNotFound, nil)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/home.tmpl",
	}

	ts := template.Must(template.New("home").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate,
	}).ParseFiles(files...))

	data := &models.ArticlesResponse{
		Articles: *articles,
	}

	if err := ts.ExecuteTemplate(w, "base", data); err != nil {
		utils.WriteJsonError(w, "Template execution error", http.StatusInternalServerError, err)
		return
	}
}

func ArticleView(w http.ResponseWriter, r *http.Request) {
	id := utils.GetID(w, r)
	if id == 0 {
		return
	}

	articles := GetArticles()
	if articles == nil {
		utils.WriteJsonError(w, "No articles found", http.StatusNotFound, nil)
		return
	}

	article := &models.Article{}
	for _, a := range *articles {
		if a.ID == id {
			article = &a
			break
		}
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/view.tmpl",
	}

	ts := template.Must(template.New("view").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate,
	}).ParseFiles(files...))

	data := &models.ArticleResponse{
		Article: article,
	}

	if err := ts.ExecuteTemplate(w, "base", data); err != nil {
		utils.WriteJsonError(w, "Template execution error", http.StatusInternalServerError, err)
		return
	}

}

func ArticleCreateForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/create.tmpl",
	}

	ts := template.Must(template.New("form").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate,
	}).ParseFiles(files...))

	if err := ts.ExecuteTemplate(w, "base", nil); err != nil {
		utils.WriteJsonError(w, "Template execution error", http.StatusInternalServerError, err)
		return
	}
}

func ArticleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJsonError(w, "Invalid request method", http.StatusMethodNotAllowed, nil)
		http.Redirect(w, r, "/new", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")
	publishedDate := r.FormValue("publishedDate")

	date, err := utils.ParseDate(publishedDate)
	if err != nil {
		utils.WriteJsonError(w, "Invalid date format", http.StatusBadRequest, err)
		return
	}

	newArticle := &models.Article{
		ID:            utils.GenerateNewID(),
		Title:         title,
		Content:       content,
		PublishedDate: date,
	}

	// Save new article to JSON file
	articles := GetArticles()
	if articles == nil {
		articles = &[]models.Article{}
	}

	*articles = append(*articles, *newArticle)

	f, err := os.Create("pkg/data/articles.json")
	if err != nil {
		utils.WriteJsonError(w, "Error creating file", http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(&articles); err != nil {
		utils.WriteJsonError(w, "Error encoding JSON", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func ArticleUpdateForm(w http.ResponseWriter, r *http.Request) {
	id := utils.GetID(w, r)
	if id == 0 {
		return
	}

	articles := GetArticles()
	if articles == nil {
		utils.WriteJsonError(w, "No articles found", http.StatusNotFound, nil)
		return
	}

	article := &models.Article{}
	for _, a := range *articles {
		if a.ID == id {
			article = &a
			break
		}
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/edit.tmpl",
	}

	ts := template.Must(template.New("edit").Funcs(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
	}).ParseFiles(files...))

	data := &models.ArticleResponse{
		Article: article,
	}

	if err := ts.ExecuteTemplate(w, "base", data); err != nil {
		utils.WriteJsonError(w, "Template execution error", http.StatusInternalServerError, err)
		return
	}
}

func ArticleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJsonError(w, "Invalid request method", http.StatusMethodNotAllowed, nil)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	id := utils.GetID(w, r)
	if id == 0 {
		return
	}

	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")
	publishedDate := r.FormValue("publishedDate")

	date, err := utils.ParseDate(publishedDate)
	if err != nil {
		utils.WriteJsonError(w, "Invalid date format", http.StatusBadRequest, err)
		return
	}

	articles := GetArticles()
	if articles == nil {
		utils.WriteJsonError(
			w, "No articles found", http.StatusNotFound, nil)
		return
	}
	var updatedArticles []models.Article
	for _, a := range *articles {
		if a.ID == id {
			updatedArticle := models.Article{
				ID:            id,
				Title:         title,
				Content:       content,
				PublishedDate: date,
			}
			updatedArticles = append(updatedArticles, updatedArticle)
		} else {
			updatedArticles = append(updatedArticles, a)
		}
	}

	f, err := os.Create("pkg/data/articles.json")
	if err != nil {
		utils.WriteJsonError(w, "Error creating file", http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(&updatedArticles); err != nil {
		utils.WriteJsonError(w, "Error encoding JSON", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func ArticleDelete(w http.ResponseWriter, r *http.Request) {
	id := utils.GetID(w, r)
	if id == 0 {
		return
	}

	articles := GetArticles()
	if articles == nil {
		utils.WriteJsonError(w, "No articles found", http.StatusNotFound, nil)
		return
	}

	var updatedArticles []models.Article
	for _, a := range *articles {
		if a.ID != id {
			updatedArticles = append(updatedArticles, a)
		}
	}

	if len(updatedArticles) == len(*articles) {
		utils.WriteJsonError(w, "Article not found", http.StatusNotFound, nil)
		return
	}

	f, err := os.Create("pkg/data/articles.json")
	if err != nil {
		utils.WriteJsonError(w, "Error creating file", http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(&updatedArticles); err != nil {
		utils.WriteJsonError(w, "Error encoding JSON", http.StatusInternalServerError, err)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	articles := GetArticles()
	if articles == nil {
		utils.WriteJsonError(w, "No articles found", http.StatusNotFound, nil)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/dashboard.tmpl",
	}

	ts := template.Must(template.New("home").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate,
	}).ParseFiles(files...))

	data := &models.ArticlesResponse{
		Articles: *articles,
	}

	if err := ts.ExecuteTemplate(w, "base", data); err != nil {
		utils.WriteJsonError(w, "Template execution error", http.StatusInternalServerError, err)
		return
	}
}
