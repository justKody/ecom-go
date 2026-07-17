// business logic
package cart

import (
	"fmt"

	"github.com/justKody/ecom-go/types"
)

func getCardItemsIds(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0")
		}
		productIds[i] = item.ProductID
	}
	return productIds, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {

	// ! map optimization to avoid loops everytime
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}
	// check if all the product in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	// calcualte the total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce the quantity in our db
	for _, item := range items {
		product := productMap[item.ProductID]
		if err := h.productStore.ReduceProductQuantity(product.ID, item.Quantity); err != nil {
			return 0, 0, err
		}
	}
	// create the order
	order := types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "123 Main St, Anytown, USA",
	}
	if err := h.orderStore.CreateOrder(order); err != nil {
		return 0, 0, err
	}
	// create order items
	for _, item := range items {
		orderItem := types.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		}
		if err := h.orderStore.CreateOrderItem(orderItem); err != nil {
			return 0, 0, err
		}
	}
	return order.ID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product not found")
		}
		if item.Quantity > product.Quantity {
			return fmt.Errorf("product %s is out of stock", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	totalPrice := 0.0
	for _, item := range cartItems {
		product := products[item.ProductID] // ! this is safe because we checked if the product is in the map
		totalPrice += product.Price * float64(item.Quantity)
	}
	return totalPrice
}
