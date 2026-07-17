package cart

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/justKody/ecom-go/service/auth"
	"github.com/justKody/ecom-go/types"
	"github.com/justKody/ecom-go/utils"
)

type Handler struct {
	orderStore   types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(orderStore types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		orderStore:   orderStore,
		productStore: productStore,
		userStore:    userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods("POST")
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	// userId from jwt token
	userID, err := auth.GetUserIdFromContext(r.Context())

	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(cart); err != nil {
		err := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	productIds, err := getCardItemsIds(cart.Items)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// get products
	ps, err := h.productStore.GetProductByIDs(productIds)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(ps, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"order_id":    orderID,
		"total_price": totalPrice,
	})
}
