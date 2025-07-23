package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	validator "github.com/go-playground/validator/v10"

	"bobshop/internal/modules/product/application"
	"bobshop/internal/modules/product/delivery/http/dto"
	"bobshop/internal/modules/product/domain"
	"bobshop/internal/platform/response"
	"bobshop/internal/platform/web"
)

type ProductHandler struct {
	service  *application.ProductService
	validate *validator.Validate
}

func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	web.BindAndValidate(c, h.validate, &req)

	product := dto.ToDomain(&req)
	if err := h.service.Create(c.Request.Context(), product); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "Product created", dto.ToResponse(product))
}

func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	var req dto.UpdateProductRequest
	web.BindAndValidate(c, h.validate, &req)

	updateFields := bson.M{
		"updated_at": time.Now(),
	}
	if req.Name != nil {
		updateFields["name"] = *req.Name
	}
	if req.Price != nil {
		updateFields["price"] = *req.Price
	}

	if err := h.service.UpdatePartial(c.Request.Context(), id, updateFields); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, err)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SimpleSuccess(c, "Product updated")
}

func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, err)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Deleted(c, "Product deleted")
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	product, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Product found", gin.H{"product": dto.FromDomain(product)})
}

func (h *ProductHandler) List(c *gin.Context) {
	var filter domain.ListFilter
	// Parse query params for filter, pagination, sort
	name := c.Query("name")
	if name != "" {
		filter.Name = &name
	}
	// Similarly for others: categories = c.QueryArray("categories") etc.

	var pagination domain.CursorPagination
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10 // default limit
	}

	// using new service of pagiantion
	pagination.Cursor = cursor
	pagination.Limit = limit

	sortStr := c.Query("sort")
	sort := domain.SortBy(sortStr)

	products, nextCursor, err := h.service.List(c.Request.Context(), &filter, &pagination, sort)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Products listed", gin.H{"products": dto.FromDomainList(products), "next_cursor": nextCursor})
}

func (h *ProductHandler) AddReview(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	var req dto.AddReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "mismatched fields", err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return
	}

	userID := uuid.MustParse(c.GetString("userId"))

	review := req.ToDomain(productID, userID)

	if err := h.service.AddReview(c.Request.Context(), review); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SimpleSuccess(c, "Review added")
}

func (h *ProductHandler) TrackRecentlyViewed(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	userID := uuid.MustParse(c.GetString("user_id")) // dry

	if err := h.service.TrackRecentlyViewedProduct(c.Request.Context(), userID, productID); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SimpleSuccess(c, "Viewed tracked")
}

func (h *ProductHandler) GetRecentlyViewed(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id")) // dry

	// dry
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}

	ids, err := h.service.GetRecentlyViewedProducts(c.Request.Context(), userID, limit)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Recently viewed", ids)
}
