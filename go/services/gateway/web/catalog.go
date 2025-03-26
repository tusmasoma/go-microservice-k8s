package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tusmasoma/go-microservice-k8s/go/pkg/request"
	"github.com/tusmasoma/go-microservice-k8s/go/pkg/response"
	"github.com/tusmasoma/go-microservice-k8s/go/services/gateway/web/entity"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
	webpb "github.com/tusmasoma/go-microservice-k8s/proto/web"
)

func (w *web) routeCatalog(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		// 認証が必要な場合はこのコメントアウトを外す
		// r.Use(w.requireAuth)
		r.Get("/", w.listCatalogItems)
		r.Post("/", w.createCatalogItem)
		r.Get("/{catalogItemID}", w.getCatalogItem)
		r.Put("/{catalogItemID}", w.updateCatalogItem)
		r.Delete("/{catalogItemID}", w.deleteCatalogItem)
		r.Get("/search/by-name", w.listCatalogItemsByName) // TODO: refactoer route
		r.Post("/search/by-ids", w.listCatalogItemsByIDs)  // TODO: refactoer route
	})
}

func (w *web) listCatalogItems(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := w.catalog.ListCatalogItems(ctx, &catalog.ListCatalogItemsRequest{})
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	body := &webpb.ListCatalogItemsResponse{
		Items: entity.NewCatalogItems(out.GetItems()).Proto(),
	}
	response.OK(body, rw, r)
}

func (w *web) createCatalogItem(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &webpb.CreateCatalogItemRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	in := &catalog.CreateCatalogItemRequest{
		Name:  req.GetName(),
		Price: req.GetPrice(),
	}
	if _, err := w.catalog.CreateCatalogItem(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.Created(rw)
}

func (w *web) getCatalogItem(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	catalogItemID := request.Param(r, "catalogItemID")
	in := &catalog.GetCatalogItemRequest{
		Id: catalogItemID,
	}
	out, err := w.catalog.GetCatalogItem(ctx, in)
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	body := &webpb.GetCatalogItemResponse{
		Item: entity.NewCatalogItem(out.GetItem()).Proto(),
	}
	response.OK(body, rw, r)
}

func (w *web) updateCatalogItem(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	catalogItemID := request.Param(r, "catalogItemID")
	req := &webpb.UpdateCatalogItemRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	in := &catalog.UpdateCatalogItemRequest{
		Id:    catalogItemID,
		Name:  req.GetName(),
		Price: req.GetPrice(),
	}
	if _, err := w.catalog.UpdateCatalogItem(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.NoContent(rw)
}

func (w *web) deleteCatalogItem(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	catalogItemID := request.Param(r, "catalogItemID")
	in := &catalog.DeleteCatalogItemRequest{
		Id: catalogItemID,
	}
	if _, err := w.catalog.DeleteCatalogItem(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.NoContent(rw)
}

func (w *web) listCatalogItemsByName(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &webpb.ListCatalogItemsByNameRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	in := &catalog.ListCatalogItemsByNameRequest{
		Name: req.GetName(),
	}
	out, err := w.catalog.ListCatalogItemsByName(ctx, in)
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	body := &webpb.ListCatalogItemsByNameResponse{
		Items: entity.NewCatalogItems(out.GetItems()).Proto(),
	}
	response.OK(body, rw, r)
}

func (w *web) listCatalogItemsByIDs(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &webpb.ListCatalogItemsByIDsRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	in := &catalog.ListCatalogItemsByIDsRequest{
		Ids: req.GetIds(),
	}
	out, err := w.catalog.ListCatalogItemsByIDs(ctx, in)
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	body := &webpb.ListCatalogItemsByIDsResponse{
		Items: entity.NewCatalogItems(out.GetItems()).Proto(),
	}
	response.OK(body, rw, r)
}
