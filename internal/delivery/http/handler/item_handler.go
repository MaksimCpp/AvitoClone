package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MaksimCpp/AvitoClone/internal/config"
	errorresponse "github.com/MaksimCpp/AvitoClone/internal/delivery/http/error"
	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
	itemusecase "github.com/MaksimCpp/AvitoClone/internal/usecase/item"
	itemimageusecase "github.com/MaksimCpp/AvitoClone/internal/usecase/item_image"
	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	createUseCase itemusecase.CreateItemUseCase
	deleteUseCase itemusecase.DeleteItemUseCase
	getByIDUseCase itemusecase.GetItemByIDUseCase
	listByUserIDUseCase itemusecase.ListItemsByUserIDUseCase
	listUseCase itemusecase.ListItemsUseCase

	uploadImageUseCase itemimageusecase.UploadImageUseCase
	deleteImageUseCase itemimageusecase.DeleteImageUseCase
	listImagesByItemIDUseCase itemimageusecase.ListImagesByItemIDUseCase

	cfg *config.Config
}

func NewItemHandler(
	createUseCase itemusecase.CreateItemUseCase,
	deleteUseCase itemusecase.DeleteItemUseCase,
	getByIDUseCase itemusecase.GetItemByIDUseCase,
	listByUserIDUseCase itemusecase.ListItemsByUserIDUseCase,
	listUseCase itemusecase.ListItemsUseCase,

	uploadImageUseCase itemimageusecase.UploadImageUseCase,
	deleteImageUseCase itemimageusecase.DeleteImageUseCase,
	listImagesByItemIDUseCase itemimageusecase.ListImagesByItemIDUseCase,

	cfg *config.Config,
) *ItemHandler {
	return &ItemHandler{
		createUseCase: createUseCase,
		deleteUseCase: deleteUseCase,
		getByIDUseCase: getByIDUseCase,
		listByUserIDUseCase: listByUserIDUseCase,
		listUseCase: listUseCase,

		uploadImageUseCase: uploadImageUseCase,
		deleteImageUseCase: deleteImageUseCase,
		listImagesByItemIDUseCase: listImagesByItemIDUseCase,

		cfg: cfg,
	}
}

type CreateItemRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Price float64 `json:"price"`
}

type ItemResponse struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Price float64 `json:"price"`
}

type ItemDetailResponse struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price float64 `json:"price"`
}

type UploadImageResponse struct {
	ImageURL string `json:"image_url"`
}

type ImagesResponse struct {
	ID int `json:"id"`
	ImageURL string `json:"image_url"`
}

// @Summary Create item
// @Description Create item
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateItemRequest true "Create item request"
// @Success 201 {object} ItemResponse
// @Failure 400 {object} errorresponse.ErrorResponse
// @Failure 401 {object} errorresponse.ErrorResponse
// @Router /items [post]
func (h *ItemHandler) Create(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req CreateItemRequest
	c.ShouldBindJSON(&req)

	if len(req.Title) < 2 {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Title length < 2."})
		return
	}

	if len(req.Description) < 2 {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Description length < 2."})
		return
	}

	itemEntity := item.Item{
		UserID: userID,
		Title: req.Title,
		Description: req.Description,
		Price: req.Price,
	}
	result, err := h.createUseCase.Execute(c.Request.Context(), &itemEntity)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}
	response := ItemResponse{
		ID: result.ID,
		Title: result.Title,
		Price: result.Price,
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Delete item
// @Description Delete item
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} errorresponse.ErrorResponse
// @Failure 401 {object} errorresponse.ErrorResponse
// @Router /items/{id} [delete]
func (h *ItemHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	base := 10
	bitSize := 64
	id, err := strconv.ParseInt(
		idParam, base, bitSize,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid item id."})
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errorresponse.ErrorResponse{Detail: "Unauthorized."})
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid user id."})
		return
	}

	itemEntity, err := h.getByIDUseCase.Execute(c.Request.Context(), int(id))
	if err != nil {
		c.JSON(http.StatusNotFound, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}
	if userID != itemEntity.UserID {
		c.JSON(http.StatusForbidden, errorresponse.ErrorResponse{Detail: "Invalid user id."})
		return
	}

	err = h.deleteUseCase.Execute(c.Request.Context(), int(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Get item by id
// @Description Get item by id
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Security BearerAuth
// @Success 200 {object} ItemDetailResponse
// @Failure 404 {object} errorresponse.ErrorResponse
// @Failure 401 {object} errorresponse.ErrorResponse
// @Router /items/{id} [get]
func (h *ItemHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	base := 10
	bitSize := 64
	id, err := strconv.ParseInt(
		idParam, base, bitSize,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid item id."})
	}

	itemEntity, err := h.getByIDUseCase.Execute(c.Request.Context(), int(id))
	if err != nil {
		if errors.Is(err, item.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, errorresponse.ErrorResponse{Detail: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	response := ItemDetailResponse{
		ID: itemEntity.ID,
		Title: itemEntity.Title,
		Description: itemEntity.Description,
		Price: itemEntity.Price,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary List items by user id
// @Description List items by user id
// @Tags items
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} []ItemResponse
// @Failure 400 {object} errorresponse.ErrorResponse
// @Router /users/{user_id}/items [get]
func (h *ItemHandler) ListByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	base := 10
	bitSize := 64
	userID, err := strconv.ParseInt(
		userIDParam, base, bitSize,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid user id."})
	}

	items, err := h.listByUserIDUseCase.Execute(c.Request.Context(), int(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	var response []ItemResponse

	for _, itemEntity := range items {
		response = append(response, ItemResponse{
			ID: itemEntity.ID,
			Title: itemEntity.Title,
			Price: itemEntity.Price,
		})
	}

	c.JSON(http.StatusOK, response)
}

// @Summary List my items
// @Description List my items
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []ItemResponse
// @Failure 400 {object} errorresponse.ErrorResponse
// @Failure 401 {object} errorresponse.ErrorResponse
// @Router /items/me [get]
func (h *ItemHandler) ListMyItems(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errorresponse.ErrorResponse{Detail: "Unauthorized."})
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid user id."})
		return
	}

	items, err := h.listByUserIDUseCase.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	var response []ItemResponse

	for _, itemEntity := range items {
		response = append(response, ItemResponse{
			ID: itemEntity.ID,
			Title: itemEntity.Title,
			Price: itemEntity.Price,
		})
	}


	c.JSON(http.StatusOK, response)
}

// @Summary List items
// @Description List items
// @Tags items
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} []ItemResponse
// @Failure 400 {object} errorresponse.ErrorResponse
// @Failure 401 {object} errorresponse.ErrorResponse
// @Router /items [get]
func (h *ItemHandler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	if limit < 1 || offset < 0 {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Limit or offset invalid."})
	}

	items, err := h.listUseCase.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	var response []ItemResponse

	for _, itemEntity := range items {
		response = append(response, ItemResponse{
			ID: itemEntity.ID,
			Title: itemEntity.Title,
			Price: itemEntity.Price,
		})
	}

	c.JSON(http.StatusOK, response)
}

// UploadImage godoc
//
//	@Summary Upload image
//	@Description Upload image for item
//	@Tags images
//	@Accept	 mpfd
//	@Produce json
//	@Security BearerAuth
//	@Param item_id	path int true "Item ID"
//	@Param image formData file true "Image file"
//	@Success 201 {object} UploadImageResponse
//	@Failure 400 {object} errorresponse.ErrorResponse
//	@Failure 401 {object} errorresponse.ErrorResponse
//	@Router	 /images/{item_id} [post]
func (h *ItemHandler) UploadImage(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errorresponse.ErrorResponse{Detail: "Unauthorized."})
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid user id."})
		return
	}

	itemID, err := strconv.Atoi(
		c.Param("item_id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid item id."})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Image is required."})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Cannot opened image."})
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	input := itemimageusecase.UploadInput{
		ItemID: itemID,
		UserID: userID,
		File: file,
		Filename: fileHeader.Filename,
		ContentType: contentType,
		Size: fileHeader.Size,
	}

	imageURL, err := h.uploadImageUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UploadImageResponse{ImageURL: imageURL})
}

// ListImages godoc
//
//	@Summary		Get item images
//	@Description	Get images for item
//	@Tags			images
//	@Produce		json
//	@Param			id	path	int	true	"Item ID"
//	@Success		200	{array}		ImagesResponse
//	@Failure		400	{object}	errorresponse.ErrorResponse
//	@Router			/items/{id}/images [get]
func (h *ItemHandler) ListImagesByItemID(c *gin.Context) {
	itemID, err := strconv.Atoi(
		c.Param("id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid item id."})
		return
	}

	images, err := h.listImagesByItemIDUseCase.Execute(c.Request.Context(), itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	var imagesResponse []ImagesResponse

	for _, image := range images {
		imagesResponse = append(imagesResponse, ImagesResponse{
			ID: image.ID,
			ImageURL: fmt.Sprintf(
				"http://%s/%s/%s",
				h.cfg.MinioEndpointImage,
				h.cfg.MinioBucket,
				image.ObjectName,
			),
		})
	}

	c.JSON(http.StatusOK, imagesResponse)
}

// DeleteImage godoc
//
//	@Summary		Delete image
//	@Description	Delete image by id
//	@Tags			images
//	@Security		BearerAuth
//	@Param			image_id	path	int	true	"Image ID"
//	@Success		204
//	@Failure		400	{object}	errorresponse.ErrorResponse
//	@Failure		401	{object}	errorresponse.ErrorResponse
//	@Failure		403	{object}	errorresponse.ErrorResponse
//	@Router			/images/{image_id} [delete]
func (h *ItemHandler) DeleteImage(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid user id."})
		return
	}

	userID := userIDValue.(int)

	imageID, err := strconv.Atoi(
		c.Param("image_id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid image id."})
		return
	}

	err = h.deleteImageUseCase.Execute(c.Request.Context(), imageID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
