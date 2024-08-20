package handler

import (
	"html/template"
	"net/http"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/entity"
	"github.com/tusmasoma/microservice-k8s-demo/catalog/usecase"
)

type CatalogItemHandler interface {
	ListCatalogItems(w http.ResponseWriter, r *http.Request)
}

type catalogItemHandler struct {
	cuc usecase.CatalogItemUseCase
}

func NewCatalogItemHandler(cuc usecase.CatalogItemUseCase) CatalogItemHandler {
	return &catalogItemHandler{
		cuc: cuc,
	}
}

func (ch *catalogItemHandler) ListCatalogItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	items, err := ch.cuc.ListCatalogItems(ctx)
	if err != nil {
		log.Error("Failed to list catalog items", log.Ferror(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		"gateway/web/templates/layout.html",
		"gateway/web/templates/list.html",
	)
	if err != nil {
		log.Error("Failed to parse template", log.Ferror(err))
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Items []entity.CatalogItem
	}{
		Items: items,
	}

	if err = tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Error("Failed to execute template", log.Ferror(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
