package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          *string            `json:"name" validate:"required,min=2,max=30"`
	ProfileUrl    *string            `json:"profile_url"`
	Password      *string            `json:"password"   validate:"required,min=6"`
	Email         *string            `json:"email"      validate:"email,required"`
	Phone         *string            `json:"phone"      validate:"required"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `josn:"refresh_token"`
	Created_At    time.Time          `json:"created_at"`
	Updated_At    time.Time          `json:"updtaed_at"`
	User_ID       string             `json:"user_id"`
	User_ROLE     *uint8             `json:"user_role"`
}

type Product struct {
	Product_ID     primitive.ObjectID `bson:"_id"`
	Product_Name   *string            `json:"product_name" binding:"required`
	Price          *uint64            `json:"price" binding:"required`
	Rating         *uint8             `json:"rating" binding:"required`
	Image          *string            `json:"image"`
	VednorID       *string            `json:"vendor_id"`
	Quantity       *uint64            `json:"quantity" binding:"required`
	IsActive       *bool              `json:"isactive" binding:"required`
	AvaliableStock *uint64            `json:"avaliable_stock" binding:"required`
}

type Address struct {
	Address_id primitive.ObjectID `bson:"_id"`
	UserID     *string            `json:"user_id"`
	House      *string            `json:"house_name" bson:"house_name"`
	Street     *string            `json:"street_name" bson:"street_name"`
	City       *string            `json:"city_name" bson:"city_name"`
	Pincode    *string            `json:"pin_code" bson:"pin_code"`
}

type ProductUser struct {
	ID         primitive.ObjectID `bson:"_id"`
	VednorID   *string            `json:"vendor_id"`
	UserID     *string            `json:"user_id"`
	Product_ID *string            `json:"product_id"`
}

type Order struct {
	OrderID   primitive.ObjectID `bson:"_id"`
	VednorID  string             `json:"vendor_id" bson:"vendor_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	ProductID string             `json:"product_id" bson:"product_id"`

	OrderStatus string    `json:"order_status" bson:"order_status"`
	OrdereredAt time.Time `json:"ordered_on" bson:"ordered_on"`
	Price       float64   `json:"total_price" bson:"total_price"`
	Discount    float64   `json:"discount" bson:"discount"`
	PaymentID   string    `json:"payment_id" bson:"payment_id"`
}

type Payment struct {
	PaymentID     primitive.ObjectID `bson:"_id"`
	PaymentMethod string             `json:"payment_method" bson:"payment_method"`
	TransactionID string             `json:"trans_id" bson:"trans_id"`
	PaymentAt     time.Time          `json:"payment_at" bson:"payment_at"`
	UserID        string             `json:"user_id" bson:"user_id"`
	PaymentStatus uint8              `json:"payment_status" bson:"payment_status"`
}

type OrderProduct struct {
}
