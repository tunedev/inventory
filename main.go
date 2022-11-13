package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type Car struct {
	Wheels  int
	Model   string
	Make    string
	MgtDate time.Time
}

type Product struct {
	id          string
	Name        string
	Price       float32
	CarInfo     Car
	AmntInStock int
}

func (p *Product) getProductInfo() {
	fmt.Println("Name:", p.Name)
	fmt.Println("Price:", p.Price)
	fmt.Println("Available Quantity:", p.AmntInStock)
	fmt.Println("--------------------")
	fmt.Println("Car Info:")
	fmt.Println("Model:", p.CarInfo.Model)
	fmt.Println("Make:", p.CarInfo.Make)
	fmt.Println("Manufactured At:", p.CarInfo.MgtDate)
	fmt.Println(p.CarInfo.Wheels, "Wheel Drive")
}

func (p *Product) sellFromProduct(qty int) error {
	if qty > p.AmntInStock {
		return errors.New("insufficient Cars available")
	}
	p.AmntInStock -= qty
	return nil
}

type OrderItem struct {
	orderId    string
	carInfo    Car
	quantity   int
	totalPrice float32
	unitPrice  float32
}

type Order struct {
	id         string
	totalPrice float32
	totalQty   int
	orderItems []OrderItem
}

type Store struct {
	products []Product
	Orders   []Order
}

type CreateNewProductParam struct {
	car   Car
	price float32
	qty   int
}

func (s *Store) addNewProduct(param CreateNewProductParam) {
	car := param.car
	newProduct := Product{
		id:          generateUniqueId(),
		Name:        car.Make + car.Model,
		CarInfo:     car,
		Price:       param.price,
		AmntInStock: param.qty,
	}
	s.products = append(s.products, newProduct)
}

func (s *Store) getTotalCarsLeft() int {
	result := 0
	for _, product := range s.products {
		result += product.AmntInStock
	}
	return result
}

func (s *Store) getTotalPriceOfCarsLeft() float32 {
	var result float32
	for _, product := range s.products {
		result += product.Price * float32(product.AmntInStock)
	}
	return result
}

func (s *Store) getProductById(productId string) (Product, error) {
	for _, p := range s.products {
		if p.id == productId {
			return p, nil
		}
	}
	return Product{}, errors.New(fmt.Sprintf("product with id: %v Not Found", productId))
}

type NewOrderParam struct {
	productId string
	qty       int
}

func (s *Store) createOrder(cart []NewOrderParam) {
	orderId := generateUniqueId()
	var orderItems []OrderItem
	var totalPrice float32
	totalQty := 0
	for _, o := range cart {
		product, errGettingProduct := s.getProductById(o.productId)
		if errGettingProduct != nil {
			log.Fatal(errGettingProduct)
		}
		errSellingProduct := product.sellFromProduct(o.qty)
		if errSellingProduct != nil {
			log.Fatal(errSellingProduct)
		}
		newOrderItem := OrderItem{
			orderId:    orderId,
			carInfo:    product.CarInfo,
			quantity:   o.qty,
			totalPrice: product.Price * float32(o.qty),
			unitPrice:  product.Price,
		}
		totalPrice += newOrderItem.totalPrice
		totalQty += newOrderItem.quantity
		orderItems = append(orderItems, newOrderItem)
	}

	s.Orders = append(s.Orders, Order{
		id:         orderId,
		totalPrice: totalPrice,
		totalQty:   totalQty,
		orderItems: orderItems,
	})
}

func (s *Store) getTotalCarsSold() int {
	result := 0
	for _, o := range s.Orders {
		result += o.totalQty
	}
	return result
}

func (s *Store) getTotalPriceOfCarsSold() float32 {
	var result float32
	for _, o := range s.Orders {
		result += o.totalPrice
	}
	return result
}

func main() {

}

func generateUniqueId() string {
	newUniqueId, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(newUniqueId)
}
