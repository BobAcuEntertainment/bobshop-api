package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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
	if err := web.BindAndValidate(c, h.validate, &req); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return
	}

	product, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "Product created", dto.ToCreateProductResponse(product))
}

func (h *ProductHandler) Update(c *gin.Context) {
	productId, err := web.GetIDParam(c)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := web.BindAndValidate(c, h.validate, &req); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return
	}

	if err := h.service.Update(c.Request.Context(), productId, req); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, err)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.NoContent(c, "Product updated")
}

func (h *ProductHandler) Delete(c *gin.Context) {
	productId, err := web.GetIDParam(c)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), productId); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, err)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.NoContent(c, "Product deleted")
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	productId, err := web.GetIDParam(c)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	product, err := h.service.GetByID(c.Request.Context(), productId)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, err)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Product found", dto.ToProductResponse(product))
}

func (h *ProductHandler) List(c *gin.Context) {
	var filter dto.ListFilterRequest
	if err := web.BindAndValidate(c, h.validate, &filter); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return
	}

	var pagination dto.CursorPaginationRequest
	if err := web.BindAndValidate(c, h.validate, &pagination); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return
	}

	var sort dto.SortRequest
	if err := web.BindAndValidate(c, h.validate, &sort); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return
	}

	products, nextCursor, err := h.service.List(c.Request.Context(), filter, pagination, sort)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Products listed", dto.ToListProductsResponse(products, nextCursor))
}

func (h *ProductHandler) AddReview(c *gin.Context) {
	userID := web.GetUserID(c)

	productID, err := web.GetIDParam(c)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	var req dto.AddReviewRequest
	if err := web.BindAndValidate(c, h.validate, &req); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return
	}

	if err := h.service.AddReview(c.Request.Context(), req, userID, productID); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SimpleSuccess(c, "Review added")
}

func (h *ProductHandler) TrackRecentlyViewed(c *gin.Context) {
	userID := web.GetUserID(c)

	productID, err := web.GetIDParam(c)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	if err := h.service.TrackRecentlyViewedProduct(c.Request.Context(), userID, productID); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SimpleSuccess(c, "Viewed tracked")
}

func (h *ProductHandler) GetRecentlyViewed(c *gin.Context) {
	userID := web.GetUserID(c)

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
