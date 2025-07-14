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

// listCatalogItems godoc
// @Summary      List all catalog items
// @Description  Get a list of all catalog items
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Success      200  {object}  webpb.ListCatalogItemsResponse
// @Failure      500  {object}  web.Error
// @Router       /api/v1/catalog [get]
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

// createCatalogItem godoc
// @Summary      Create a new catalog item
// @Description  Create a new catalog item with name and price
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Param        request body webpb.CreateCatalogItemRequest true "Catalog item data"
// @Success      201
// @Failure      400  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/catalog [post]
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

// getCatalogItem godoc
// @Summary      Get a catalog item
// @Description  Get a catalog item by ID
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Param        catalogItemID path string true "Catalog Item ID"
// @Success      200  {object}  webpb.GetCatalogItemResponse
// @Failure      404  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/catalog/{catalogItemID} [get]
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

// updateCatalogItem godoc
// @Summary      Update a catalog item
// @Description  Update a catalog item by ID
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Param        catalogItemID path string true "Catalog Item ID"
// @Param        request body webpb.UpdateCatalogItemRequest true "Updated catalog item data"
// @Success      204
// @Failure      400  {object}  web.Error
// @Failure      404  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/catalog/{catalogItemID} [put]
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

// deleteCatalogItem godoc
// @Summary      Delete a catalog item
// @Description  Delete a catalog item by ID
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Param        catalogItemID path string true "Catalog Item ID"
// @Success      204
// @Failure      404  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/catalog/{catalogItemID} [delete]
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

// listCatalogItemsByName godoc
// @Summary      Search catalog items by name
// @Description  Search catalog items by name
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Param        name query string true "Item name to search"
// @Success      200  {object}  webpb.ListCatalogItemsByNameResponse
// @Failure      400  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/catalog/search/by-name [get]
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

// listCatalogItemsByIDs godoc
// @Summary      Get catalog items by IDs
// @Description  Get multiple catalog items by their IDs
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Param        request body webpb.ListCatalogItemsByIDsRequest true "List of catalog item IDs"
// @Success      200  {object}  webpb.ListCatalogItemsByIDsResponse
// @Failure      400  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/catalog/search/by-ids [post]
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
