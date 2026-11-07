package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleBuyer  UserRole = "buyer"
	RoleSeller UserRole = "seller"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"`
	Role         UserRole  `gorm:"type:varchar(20);not null;check:role IN ('buyer', 'seller')" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SellerID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"seller_id"`
	Seller      *User          `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:numeric(14,2);not null;check:price >= 0" json:"price"`
	Stock       int            `gorm:"type:integer;not null;default:0;check:stock >= 0" json:"stock"`
	ImageURL    *string        `gorm:"type:text" json:"image_url"`
	IsActive    bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Cart struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BuyerID   uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"buyer_id"`
	Items     []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CartID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_cart_product" json:"cart_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_cart_product" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int       `gorm:"type:integer;not null;check:quantity > 0" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusConfirmed  OrderStatus = "confirmed"
	StatusProcessing OrderStatus = "processing"
	StatusShipped    OrderStatus = "shipped"
	StatusCompleted  OrderStatus = "completed"
	StatusCancelled  OrderStatus = "cancelled"
	StatusRejected   OrderStatus = "rejected"
)

type ShippingAddress struct {
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	City          string `json:"city"`
	Province      string `json:"province,omitempty"`
	PostalCode    string `json:"postal_code"`
}

func (a ShippingAddress) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ShippingAddress) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed for ShippingAddress")
	}
	return json.Unmarshal(bytes, a)
}

type Order struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BuyerID         uuid.UUID       `gorm:"type:uuid;not null;index" json:"buyer_id"`
	Buyer           *User           `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	SellerID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"seller_id"`
	Seller          *User           `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	TotalAmount     float64         `gorm:"type:numeric(14,2);not null;check:total_amount >= 0" json:"total_amount"`
	Status          OrderStatus     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ShippingAddress ShippingAddress `gorm:"type:jsonb;not null" json:"shipping_address"`
	Notes           *string         `gorm:"type:text" json:"notes,omitempty"`
	Items           []OrderItem     `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type OrderItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID   uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductName string    `gorm:"type:varchar(150);not null" json:"product_name"`
	UnitPrice   float64   `gorm:"type:numeric(14,2);not null;check:unit_price >= 0" json:"unit_price"`
	Quantity    int       `gorm:"type:integer;not null;check:quantity > 0" json:"quantity"`
	Subtotal    float64   `gorm:"type:numeric(14,2);not null;check:subtotal >= 0" json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
}

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"type:text;not null;uniqueIndex" json:"token_hash"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

