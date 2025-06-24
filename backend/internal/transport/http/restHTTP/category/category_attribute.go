package category

// import (
// 	"context"
// 	"log/slog"
// 	"net/http"

// 	"github.com/Neimess/zorkin-store-project/internal/domain"
// 	"github.com/Neimess/zorkin-store-project/internal/transport/dto"
// 	"github.com/Neimess/zorkin-store-project/pkg/httputils"
// 	"github.com/go-chi/render"
// 	"github.com/mailru/easyjson"
// )

// type CategoryAttributeService interface {
// 	Create(ctx context.Context, ca *domain.CategoryAttribute) error
// 	GetByCategory(ctx context.Context, categoryID int64) ([]domain.CategoryAttribute, error)
// 	Update(ctx context.Context, ca *domain.CategoryAttribute) error
// 	Delete(ctx context.Context, categoryID, attributeID int64) error
// }

// type CategoryAttributeHandler struct {
// 	srv CategoryAttributeService
// 	log *slog.Logger
// }

// func NewCategoryAttributeHandler(srv CategoryAttributeService, log *slog.Logger) *CategoryAttributeHandler {
// 	return &CategoryAttributeHandler{
// 		srv: srv,
// 		log: log.With("component", "transport.http.restHTTP.category_attribute"),
// 	}
// }

// // CreateCategoryAttribute godoc
// // @Summary      Bind attribute to category
// // @Description  Создаёт связь «категория ↔ атрибут» с указанием обязательности и приоритета.
// // @Tags         category-attributes
// // @Accept       json
// // @Produce      json
// // @Security     BearerAuth
// // @Param        body  body  dto.CategoryAttributeRequest  true  "Mapping to create"
// // @Success      201   "Created"
// // @Failure      400   {object}  dto.ErrorResponse  "Bad request"
// // @Failure      401   {object}  dto.ErrorResponse  "Unauthorized"
// // @Failure      500   {object}  dto.ErrorResponse  "Internal server error"
// // @Router       /api/admin/category/{categoryID}/attribute [post]
// func (h *CategoryAttributeHandler) Create(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	log := h.log.With("op", "category_attribute.Create")
// 	defer r.Body.Close()
// 	// TODO Добавить обратки ошибок SQLSTATE 23505 (уникальное ограничение) и 23503 (нарушение внешнего ключа) и оных
// 	var req dto.CategoryAttributeRequest
// 	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
// 		log.Warn("invalid JSON", slog.Any("error", err))
// 		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
// 		return
// 	}

// 	d := &domain.CategoryAttribute{
// 		CategoryID:  req.CategoryID,
// 		AttributeID: req.AttributeID,
// 		IsRequired:  req.IsRequired,
// 		Priority:    req.Priority,
// 	}

// 	// TODO Добавить обратки ошибок SQLSTATE 23505 (уникальное ограничение) и 23503 (нарушение внешнего ключа) и оных
// 	if err := h.srv.Create(ctx, d); err != nil {
// 		log.Error("create failed", slog.Any("error", err))
// 		httputils.WriteError(w, http.StatusInternalServerError, "failed to create mapping")
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }

// // GetCategoryAttributes godoc
// // @Summary      List attributes of a category
// // @Description  Возвращает все атрибуты, привязанные к категории, в порядке приоритета.
// // @Tags         category-attributes
// // @Accept       json
// // @Produce      json
// // @Security     BearerAuth
// // @Param        categoryID  path  int  true  "Category ID"
// // @Success      200  {array}  dto.CategoryAttributeResponse
// // @Failure      400  {object}  dto.ErrorResponse  "Invalid ID"
// // @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// // @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// // @Router       /api/admin/category/{categoryID}/attribute [get]
// func (h *CategoryAttributeHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	log := h.log.With("op", "category_attribute.GetByCategory")

// 	categoryID, err := httputils.IDFromURL(r, "categoryID")
// 	if err != nil {
// 		log.Warn("invalid category_id", slog.Any("error", err))
// 		httputils.WriteError(w, http.StatusBadRequest, "invalid category_id")
// 		return
// 	}
// 	res, err := h.srv.GetByCategory(ctx, categoryID)
// 	// TODO Добавить обратки ошибок SQLSTATE 23505 (уникальное ограничение) и 23503 (нарушение внешнего ключа) и оных
// 	if err != nil {
// 		log.Error("get failed", slog.Any("error", err))
// 		httputils.WriteError(w, http.StatusInternalServerError, "failed to fetch mappings")
// 		return
// 	}
// 	render.Status(r, http.StatusOK)
// 	render.JSON(w, r, res)
// }

// // UpdateCategoryAttribute godoc
// // @Summary      Update mapping
// // @Description  Обновляет обязательность/приоритет связки «категория ↔ атрибут».
// // @Tags         category-attributes
// // @Accept       json
// // @Produce      json
// // @Security     BearerAuth
// // @Param        categoryID    path  int                           true  "Category ID"
// // @Param        attributeID   path  int                           true  "Attribute ID"
// // @Param        body          body  dto.CategoryAttributeRequest  true  "New values"
// // @Success      204  "No Content"
// // @Failure      400  {object}  dto.ErrorResponse  "Bad request"
// // @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// // @Failure      404  {object}  dto.ErrorResponse  "Mapping not found"
// // @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// // @Router       /api/admin/category/{categoryID}/attribute/{attributeID} [put]
// func (h *CategoryAttributeHandler) Update(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	log := h.log.With("op", "category_attribute.Update")
// 	defer r.Body.Close()

// 	var req dto.CategoryAttributeRequest
// 	// TODO Добавить обратки ошибок SQLSTATE 23505 (уникальное ограничение) и 23503 (нарушение внешнего ключа) и оных
// 	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
// 		log.Warn("invalid JSON", slog.Any("error", err))
// 		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
// 		return
// 	}

// 	d := &domain.CategoryAttribute{
// 		CategoryID:  req.CategoryID,
// 		AttributeID: req.AttributeID,
// 		IsRequired:  req.IsRequired,
// 		Priority:    req.Priority,
// 	}

// 	if err := h.srv.Update(ctx, d); err != nil {
// 		log.Error("update failed", slog.Any("error", err))
// 		httputils.WriteError(w, http.StatusInternalServerError, "update failed")
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// // DeleteCategoryAttribute godoc
// // @Summary      Delete mapping
// // @Description  Удаляет связь атрибута и категории.
// // @Tags         category-attributes
// // @Accept       json
// // @Produce      json
// // @Security     BearerAuth
// // @Param        categoryID   path  int  true  "Category ID"
// // @Param        attributeID  path  int  true  "Attribute ID"
// // @Success      204  "No Content"
// // @Failure      400  {object}  dto.ErrorResponse  "Bad request"
// // @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// // @Failure      404  {object}  dto.ErrorResponse  "Mapping not found"
// // @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// // @Router       /api/admin/category/{categoryID}/attribute/{attributeID} [delete]
// func (h *CategoryAttributeHandler) Delete(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	log := h.log.With("op", "category_attribute.Delete")

// 	categoryID, err := httputils.IDFromURL(r, "id")
// 	// TODO Добавить обратки ошибок SQLSTATE 23505 (уникальное ограничение) и 23503 (нарушение внешнего ключа) и оных
// 	if err != nil {
// 		httputils.WriteError(w, http.StatusBadRequest, "invalid category_id")
// 		return
// 	}
// 	attributeID, err := httputils.IDFromURL(r, "attribute_id")
// 	if err != nil {
// 		httputils.WriteError(w, http.StatusBadRequest, "invalid attribute_id")
// 		return
// 	}

// 	if err := h.srv.Delete(ctx, categoryID, attributeID); err != nil {
// 		log.Error("delete failed", slog.Any("error", err))
// 		httputils.WriteError(w, http.StatusInternalServerError, "delete failed")
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }
