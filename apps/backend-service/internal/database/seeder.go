package database

import (
	"log"

	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/pkg/hash"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedInitialData(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return // Data already seeded
	}

	log.Println("Seeding initial demo data for learning e-commerce platform...")

	passHash, _ := hash.HashPassword("Password123")

	// 1. Demo Seller
	sellerID := uuid.New()
	seller := model.User{
		ID:           sellerID,
		Name:         "TechGear Official Store",
		Email:        "seller@techgear.com",
		PasswordHash: passHash,
		Role:         model.RoleSeller,
	}
	_ = db.Create(&seller).Error

	// 2. Demo Buyer
	buyerID := uuid.New()
	buyer := model.User{
		ID:           buyerID,
		Name:         "Naabil Fikry",
		Email:        "buyer@example.com",
		PasswordHash: passHash,
		Role:         model.RoleBuyer,
	}
	_ = db.Create(&buyer).Error

	// Create buyer cart
	cart := model.Cart{
		ID:      uuid.New(),
		BuyerID: buyerID,
	}
	_ = db.Create(&cart).Error

	// 3. Demo Products
	img1 := "https://images.unsplash.com/photo-1587829741301-dc798b83add3?auto=format&fit=crop&w=800&q=80"
	img2 := "https://images.unsplash.com/photo-1615663245857-ac93bb7c39e7?auto=format&fit=crop&w=800&q=80"
	img3 := "https://images.unsplash.com/photo-1546868871-7041f2a55e12?auto=format&fit=crop&w=800&q=80"
	img4 := "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?auto=format&fit=crop&w=800&q=80"
	img5 := "https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?auto=format&fit=crop&w=800&q=80"
	img6 := "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?auto=format&fit=crop&w=800&q=80"

	products := []model.Product{
		{
			ID:          uuid.New(),
			SellerID:    sellerID,
			Name:        "Mechanical Keyboard TKL RGB",
			Description: "Keyboard mechanical tenkeyless dengan Gateron Yellow Pro switches, double-shot PBT keycaps, dan hotswappable PCB 5-pin. Sangat nyaman untuk mengetik dan gaming.",
			Price:       650000,
			Stock:       14,
			ImageURL:    &img1,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			SellerID:    sellerID,
			Name:        "Wireless Ergonomic Mouse 26K DPI",
			Description: "Mouse nirkabel ultra-lightweight 58 gram dengan sensor optik 26.000 DPI, koneksi 2.4GHz tri-mode, dan baterai tahan hingga 80 jam pemakaian.",
			Price:       420000,
			Stock:       25,
			ImageURL:    &img2,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			SellerID:    sellerID,
			Name:        "Minimalist Smart Watch V2",
			Description: "Jam tangan pintar dengan display AMOLED 1.43 inci, pelacak detak jantung, monitor SpO2, tahan air 5ATM, dan baterai hingga 14 hari.",
			Price:       890000,
			Stock:       9,
			ImageURL:    &img3,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			SellerID:    sellerID,
			Name:        "Studio Monitor Wireless Headphones",
			Description: "Headphone nirkabel dengan Active Noise Cancellation (ANC), driver titanium 40mm, audio resolusi tinggi, dan busa memori kulit berpori.",
			Price:       1250000,
			Stock:       6,
			ImageURL:    &img4,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			SellerID:    sellerID,
			Name:        "Monitor ScreenBar Eye-Care Light",
			Description: "Lampu monitor asimetris anti-silau dengan kontrol kecerahan sentuh nirkabel, suhu warna 2700K-6500K, dan daya USB Type-C.",
			Price:       350000,
			Stock:       18,
			ImageURL:    &img5,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			SellerID:    sellerID,
			Name:        "Adjustable Aluminum Laptop Stand",
			Description: "Stand laptop bahan paduan aluminium kokoh dapat dilipat, mendukung sirkulasi udara maksimal, kompatibel laptop ukuran 11 hingga 17 inci.",
			Price:       175000,
			Stock:       30,
			ImageURL:    &img6,
			IsActive:    true,
		},
	}

	for _, p := range products {
		_ = db.Create(&p).Error
	}

	log.Println("Initial demo data seeded successfully (1 Seller, 1 Buyer, 6 Products)")
}

