package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/pkg/hash"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedInitialData runs initial seeding if the database has no users
func SeedInitialData(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return // Data already seeded
	}
	SeedAllData(db, false)
}

// SeedAllData seeds comprehensive dummy data (Users, Products, Carts, Orders, Payments)
func SeedAllData(db *gorm.DB, force bool) {
	if force {
		log.Println("Force seeding requested: cleaning existing transactional data...")
		_ = db.Exec("TRUNCATE TABLE payments, order_items, orders, cart_items, carts, products, refresh_tokens, users CASCADE").Error
	}

	log.Println("Seeding comprehensive dummy data for Simple E-Commerce Platform...")

	defaultPassHash, err := hash.HashPassword("Password123")
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}

	now := time.Now()

	// ==========================================
	// 1. USERS: SELLERS
	// ==========================================
	seller1ID := uuid.New()
	seller1 := model.User{
		ID:           seller1ID,
		Name:         "TechGear Official Store",
		Email:        "seller@techgear.com",
		PasswordHash: defaultPassHash,
		Role:         model.RoleSeller,
		CreatedAt:    now.Add(-30 * 24 * time.Hour),
		UpdatedAt:    now.Add(-30 * 24 * time.Hour),
	}

	seller2ID := uuid.New()
	seller2 := model.User{
		ID:           seller2ID,
		Name:         "AudioHub Indonesia",
		Email:        "seller@audiohub.com",
		PasswordHash: defaultPassHash,
		Role:         model.RoleSeller,
		CreatedAt:    now.Add(-25 * 24 * time.Hour),
		UpdatedAt:    now.Add(-25 * 24 * time.Hour),
	}

	seller3ID := uuid.New()
	seller3 := model.User{
		ID:           seller3ID,
		Name:         "Minimalist Desk Setup",
		Email:        "seller@minimalist.com",
		PasswordHash: defaultPassHash,
		Role:         model.RoleSeller,
		CreatedAt:    now.Add(-20 * 24 * time.Hour),
		UpdatedAt:    now.Add(-20 * 24 * time.Hour),
	}

	seller4ID := uuid.New()
	seller4 := model.User{
		ID:           seller4ID,
		Name:         "Demo Seller Store",
		Email:        "seller@example.com",
		PasswordHash: defaultPassHash,
		Role:         model.RoleSeller,
		CreatedAt:    now.Add(-15 * 24 * time.Hour),
		UpdatedAt:    now.Add(-15 * 24 * time.Hour),
	}

	sellers := []model.User{seller1, seller2, seller3, seller4}
	for _, s := range sellers {
		_ = db.Create(&s).Error
	}

	// ==========================================
	// 2. USERS: BUYERS
	// ==========================================
	buyer1ID := uuid.New()
	buyer1 := model.User{
		ID:           buyer1ID,
		Name:         "Naabil Fikry",
		Email:        "buyer@example.com",
		PasswordHash: defaultPassHash,
		Role:         model.RoleBuyer,
		CreatedAt:    now.Add(-10 * 24 * time.Hour),
		UpdatedAt:    now.Add(-10 * 24 * time.Hour),
	}

	buyer2ID := uuid.New()
	buyer2 := model.User{
		ID:           buyer2ID,
		Name:         "Budi Santoso",
		Email:        "budi@example.com",
		PasswordHash: defaultPassHash,
		Role:         model.RoleBuyer,
		CreatedAt:    now.Add(-9 * 24 * time.Hour),
		UpdatedAt:    now.Add(-9 * 24 * time.Hour),
	}

	buyer3ID := uuid.New()
	buyer3 := model.User{
		ID:           buyer3ID,
		Name:         "Siti Rahmah",
		Email:        "siti@example.com",
		PasswordHash: defaultPassHash,
		Role:         model.RoleBuyer,
		CreatedAt:    now.Add(-8 * 24 * time.Hour),
		UpdatedAt:    now.Add(-8 * 24 * time.Hour),
	}

	buyer4ID := uuid.New()
	buyer4 := model.User{
		ID:           buyer4ID,
		Name:         "Alex Pratama",
		Email:        "alex@example.com",
		PasswordHash: defaultPassHash,
		Role:         model.RoleBuyer,
		CreatedAt:    now.Add(-7 * 24 * time.Hour),
		UpdatedAt:    now.Add(-7 * 24 * time.Hour),
	}

	buyers := []model.User{buyer1, buyer2, buyer3, buyer4}
	for _, b := range buyers {
		_ = db.Create(&b).Error
	}

	// ==========================================
	// 3. PRODUCTS (12 Items)
	// ==========================================
	imgP1 := "https://images.unsplash.com/photo-1587829741301-dc798b83add3?auto=format&fit=crop&w=800&q=80"
	imgP2 := "https://images.unsplash.com/photo-1615663245857-ac93bb7c39e7?auto=format&fit=crop&w=800&q=80"
	imgP3 := "https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?auto=format&fit=crop&w=800&q=80"
	imgP4 := "https://images.unsplash.com/photo-1583863788434-e58a36330cf0?auto=format&fit=crop&w=800&q=80"
	imgP5 := "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?auto=format&fit=crop&w=800&q=80"
	imgP6 := "https://images.unsplash.com/photo-1590658268037-6bf12165a8df?auto=format&fit=crop&w=800&q=80"
	imgP7 := "https://images.unsplash.com/photo-1545454675-3531b543be5d?auto=format&fit=crop&w=800&q=80"
	imgP8 := "https://images.unsplash.com/photo-1590658006821-04f4008d5717?auto=format&fit=crop&w=800&q=80"
	imgP9 := "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?auto=format&fit=crop&w=800&q=80"
	imgP10 := "https://images.unsplash.com/photo-1544816155-12df9643f363?auto=format&fit=crop&w=800&q=80"
	imgP11 := "https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=800&q=80"
	imgP12 := "https://images.unsplash.com/photo-1546868871-7041f2a55e12?auto=format&fit=crop&w=800&q=80"

	prod1ID := uuid.New()
	prod2ID := uuid.New()
	prod3ID := uuid.New()
	prod4ID := uuid.New()
	prod5ID := uuid.New()
	prod6ID := uuid.New()
	prod7ID := uuid.New()
	prod8ID := uuid.New()
	prod9ID := uuid.New()
	prod10ID := uuid.New()
	prod11ID := uuid.New()
	prod12ID := uuid.New()

	products := []model.Product{
		// TechGear Store Products
		{
			ID:          prod1ID,
			SellerID:    seller1ID,
			Name:        "Mechanical Keyboard TKL RGB",
			Description: "Keyboard mechanical tenkeyless dengan Gateron Yellow Pro switches, double-shot PBT keycaps, dan hotswappable PCB 5-pin. Nyaman untuk mengetik dan gaming.",
			Price:       650000,
			Stock:       14,
			ImageURL:    &imgP1,
			IsActive:    true,
			CreatedAt:   now.Add(-6 * 24 * time.Hour),
			UpdatedAt:   now.Add(-6 * 24 * time.Hour),
		},
		{
			ID:          prod2ID,
			SellerID:    seller1ID,
			Name:        "Wireless Ergonomic Mouse 26K DPI",
			Description: "Mouse nirkabel ultra-lightweight 58 gram dengan sensor optik 26.000 DPI, koneksi 2.4GHz tri-mode, dan baterai tahan hingga 80 jam pemakaian.",
			Price:       420000,
			Stock:       25,
			ImageURL:    &imgP2,
			IsActive:    true,
			CreatedAt:   now.Add(-6 * 24 * time.Hour),
			UpdatedAt:   now.Add(-6 * 24 * time.Hour),
		},
		{
			ID:          prod3ID,
			SellerID:    seller1ID,
			Name:        "Monitor ScreenBar Eye-Care Light",
			Description: "Lampu monitor asimetris anti-silau dengan kontrol kecerahan sentuh nirkabel, suhu warna 2700K-6500K, dan daya USB Type-C.",
			Price:       350000,
			Stock:       18,
			ImageURL:    &imgP3,
			IsActive:    true,
			CreatedAt:   now.Add(-5 * 24 * time.Hour),
			UpdatedAt:   now.Add(-5 * 24 * time.Hour),
		},
		{
			ID:          prod4ID,
			SellerID:    seller1ID,
			Name:        "USB-C GaN Fast Charger 100W 4-Port",
			Description: "Kepala charger GaN 4-port (3x USB-C, 1x USB-A) dengan dukungan Power Delivery 3.0 dan PPS. Mampu mengisi laptop dan smartphone sekaligus.",
			Price:       345000,
			Stock:       30,
			ImageURL:    &imgP4,
			IsActive:    true,
			CreatedAt:   now.Add(-5 * 24 * time.Hour),
			UpdatedAt:   now.Add(-5 * 24 * time.Hour),
		},

		// AudioHub Indonesia Products
		{
			ID:          prod5ID,
			SellerID:    seller2ID,
			Name:        "Studio Monitor Wireless Headphones ANC",
			Description: "Headphone over-ear nirkabel dengan Active Noise Cancellation (ANC), driver titanium 40mm, audio resolusi tinggi, dan busa memori berpori.",
			Price:       1250000,
			Stock:       8,
			ImageURL:    &imgP5,
			IsActive:    true,
			CreatedAt:   now.Add(-4 * 24 * time.Hour),
			UpdatedAt:   now.Add(-4 * 24 * time.Hour),
		},
		{
			ID:          prod6ID,
			SellerID:    seller2ID,
			Name:        "Hi-Res In-Ear Monitor (IEM) Dual Driver",
			Description: "Earphone IEM audiophile dengan konfigurasi hybrid 1 Balanced Armature + 1 Dynamic Driver, kabel braided berlapis perak 2-pin lepas pasang.",
			Price:       480000,
			Stock:       20,
			ImageURL:    &imgP6,
			IsActive:    true,
			CreatedAt:   now.Add(-4 * 24 * time.Hour),
			UpdatedAt:   now.Add(-4 * 24 * time.Hour),
		},
		{
			ID:          prod7ID,
			SellerID:    seller2ID,
			Name:        "Active Desktop Studio Speakers 50W",
			Description: "Sepasang speaker monitor aktif 50W RMS dengan input optik, Bluetooth 5.0 aptX, dan output subwoofer terpisah. Bass bulat dan vokal jernih.",
			Price:       1850000,
			Stock:       5,
			ImageURL:    &imgP7,
			IsActive:    true,
			CreatedAt:   now.Add(-3 * 24 * time.Hour),
			UpdatedAt:   now.Add(-3 * 24 * time.Hour),
		},
		{
			ID:          prod8ID,
			SellerID:    seller2ID,
			Name:        "USB Condenser Cardioid Podcast Microphone",
			Description: "Mikrofon rekaman USB plug-and-play pola polar cardioid dengan shock mount, pop filter logam, dan kenop gain terpadu. Cocok untuk meeting dan podcast.",
			Price:       590000,
			Stock:       12,
			ImageURL:    &imgP8,
			IsActive:    true,
			CreatedAt:   now.Add(-3 * 24 * time.Hour),
			UpdatedAt:   now.Add(-3 * 24 * time.Hour),
		},

		// Minimalist Desk Setup Products
		{
			ID:          prod9ID,
			SellerID:    seller3ID,
			Name:        "Adjustable Aluminum Laptop Stand",
			Description: "Stand laptop bahan paduan aluminium kokoh dapat dilipat, mendukung sirkulasi udara maksimal, kompatibel laptop ukuran 11 hingga 17 inci.",
			Price:       175000,
			Stock:       30,
			ImageURL:    &imgP9,
			IsActive:    true,
			CreatedAt:   now.Add(-2 * 24 * time.Hour),
			UpdatedAt:   now.Add(-2 * 24 * time.Hour),
		},
		{
			ID:          prod10ID,
			SellerID:    seller3ID,
			Name:        "Premium Felt & Vegan Leather Desk Mat 90x40",
			Description: "Alas meja kerja premium ukuran 90x40cm dua sisi: kulit sintetis tahan air dan felt wol lembut. Memberikan sentuhan elegan dan rapi pada meja kerja.",
			Price:       185000,
			Stock:       22,
			ImageURL:    &imgP10,
			IsActive:    true,
			CreatedAt:   now.Add(-2 * 24 * time.Hour),
			UpdatedAt:   now.Add(-2 * 24 * time.Hour),
		},
		{
			ID:          prod11ID,
			SellerID:    seller3ID,
			Name:        "Under-Desk Steel Cable Management Tray",
			Description: "Baki manajemen kabel bawah meja berbahan baja galvanis tanpa perlu bor (clamp-on). Menampung stopkontak dan adaptor dengan rapi.",
			Price:       120000,
			Stock:       40,
			ImageURL:    &imgP11,
			IsActive:    true,
			CreatedAt:   now.Add(-1 * 24 * time.Hour),
			UpdatedAt:   now.Add(-1 * 24 * time.Hour),
		},
		{
			ID:          prod12ID,
			SellerID:    seller3ID,
			Name:        "Solid Walnut Wood Wrist Rest for TKL",
			Description: "Sandaran pergelangan tangan ergonomis berbahan kayu kenari solid (natural walnut) dengan kaki karet anti-selip. Menghilangkan pegal saat mengetik.",
			Price:       145000,
			Stock:       15,
			ImageURL:    &imgP12,
			IsActive:    true,
			CreatedAt:   now.Add(-1 * 24 * time.Hour),
			UpdatedAt:   now.Add(-1 * 24 * time.Hour),
		},
	}

	for _, p := range products {
		_ = db.Create(&p).Error
	}

	// ==========================================
	// 4. CARTS & CART ITEMS
	// ==========================================
	// Buyer 1 Cart (Naabil Fikry) has 2 items from TechGear Store
	cart1ID := uuid.New()
	cart1 := model.Cart{
		ID:        cart1ID,
		BuyerID:   buyer1ID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_ = db.Create(&cart1).Error

	_ = db.Create(&model.CartItem{
		ID:        uuid.New(),
		CartID:    cart1ID,
		ProductID: prod2ID, // Mouse
		Quantity:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}).Error

	_ = db.Create(&model.CartItem{
		ID:        uuid.New(),
		CartID:    cart1ID,
		ProductID: prod3ID, // ScreenBar
		Quantity:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}).Error

	// Buyer 2, 3, 4 empty or active carts
	cart2 := model.Cart{ID: uuid.New(), BuyerID: buyer2ID, CreatedAt: now, UpdatedAt: now}
	_ = db.Create(&cart2).Error

	cart3ID := uuid.New()
	cart3 := model.Cart{ID: cart3ID, BuyerID: buyer3ID, CreatedAt: now, UpdatedAt: now}
	_ = db.Create(&cart3).Error
	_ = db.Create(&model.CartItem{
		ID:        uuid.New(),
		CartID:    cart3ID,
		ProductID: prod6ID, // IEM from AudioHub
		Quantity:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}).Error

	cart4 := model.Cart{ID: uuid.New(), BuyerID: buyer4ID, CreatedAt: now, UpdatedAt: now}
	_ = db.Create(&cart4).Error

	// ==========================================
	// 5. SAMPLE ORDERS IN VARIOUS STATUSES & PAYMENTS
	// ==========================================
	addrBandung := model.ShippingAddress{
		RecipientName: "Naabil Fikry",
		Phone:         "081234567890",
		Address:       "Jl. Dago Asri No. 12, Coblong",
		City:          "Bandung",
		Province:      "Jawa Barat",
		PostalCode:    "40135",
	}

	addrJakarta := model.ShippingAddress{
		RecipientName: "Budi Santoso",
		Phone:         "081398765432",
		Address:       "Apartemen Senopati Suites Lt. 15, Kebayoran Baru",
		City:          "Jakarta Selatan",
		Province:      "DKI Jakarta",
		PostalCode:    "12190",
	}

	addrSurabaya := model.ShippingAddress{
		RecipientName: "Siti Rahmah",
		Phone:         "085712348765",
		Address:       "Jl. Kertajaya Indah Timur No. 45, Sukolilo",
		City:          "Surabaya",
		Province:      "Jawa Timur",
		PostalCode:    "60117",
	}

	addrYogyakarta := model.ShippingAddress{
		RecipientName: "Alex Pratama",
		Phone:         "082134567891",
		Address:       "Jl. Kaliurang KM 5.5, Depok, Sleman",
		City:          "Yogyakarta",
		Province:      "DI Yogyakarta",
		PostalCode:    "55281",
	}

	inv1 := "INV-20261107-DEMO001"
	inv2 := "INV-20261107-DEMO002"
	inv3 := "INV-20261107-DEMO003"
	inv4 := "INV-20261107-DEMO004"
	inv5 := "INV-20261107-DEMO005"
	inv6 := "INV-20261107-DEMO006"

	payUrl1 := fmt.Sprintf("http://localhost:5173/orders/%s?doku_checkout=1&invoice=%s&amount=650000", buyer1ID.String(), inv1)
	payUrl2 := fmt.Sprintf("https://jokul.doku.com/checkout/link/%s", inv2)
	payUrl3 := fmt.Sprintf("https://jokul.doku.com/checkout/link/%s", inv3)
	payUrl4 := fmt.Sprintf("https://jokul.doku.com/checkout/link/%s", inv4)
	payUrl5 := fmt.Sprintf("https://jokul.doku.com/checkout/link/%s", inv5)
	payUrl6 := fmt.Sprintf("http://localhost:5173/orders/%s?doku_checkout=1&invoice=%s&amount=175000", buyer1ID.String(), inv6)

	note1 := "Tolong packing kayu dan bubble wrap tebal ya min."
	note3 := "Kirim hari kerja ya, terima kasih."

	// Order 1: Pending & Unpaid (Ready for Buyer to Pay with DOKU)
	order1ID := uuid.New()
	ord1 := model.Order{
		ID:              order1ID,
		BuyerID:         buyer1ID,
		SellerID:        seller1ID,
		TotalAmount:     650000,
		Status:          model.StatusPending,
		PaymentStatus:   model.PaymentStatusUnpaid,
		InvoiceNumber:   &inv1,
		PaymentURL:      &payUrl1,
		ShippingAddress: addrBandung,
		Notes:           &note1,
		CreatedAt:       now.Add(-2 * time.Hour),
		UpdatedAt:       now.Add(-2 * time.Hour),
	}
	_ = db.Create(&ord1).Error
	_ = db.Create(&model.OrderItem{
		ID:          uuid.New(),
		OrderID:     order1ID,
		ProductID:   prod1ID,
		ProductName: "Mechanical Keyboard TKL RGB",
		UnitPrice:   650000,
		Quantity:    1,
		Subtotal:    650000,
		CreatedAt:   now.Add(-2 * time.Hour),
	}).Error
	_ = db.Create(&model.Payment{
		ID:            uuid.New(),
		OrderID:       order1ID,
		InvoiceNumber: inv1,
		Amount:        650000,
		Status:        "PENDING",
		CreatedAt:     now.Add(-2 * time.Hour),
		UpdatedAt:     now.Add(-2 * time.Hour),
	}).Error

	// Order 2: Confirmed & Paid via DOKU BCA VA (Ready for Seller to Process)
	order2ID := uuid.New()
	ord2 := model.Order{
		ID:              order2ID,
		BuyerID:         buyer2ID,
		SellerID:        seller1ID,
		TotalAmount:     345000,
		Status:          model.StatusConfirmed,
		PaymentStatus:   model.PaymentStatusPaid,
		InvoiceNumber:   &inv2,
		PaymentURL:      &payUrl2,
		ShippingAddress: addrJakarta,
		CreatedAt:       now.Add(-5 * time.Hour),
		UpdatedAt:       now.Add(-4 * time.Hour),
	}
	_ = db.Create(&ord2).Error
	_ = db.Create(&model.OrderItem{
		ID:          uuid.New(),
		OrderID:     order2ID,
		ProductID:   prod4ID,
		ProductName: "USB-C GaN Fast Charger 100W 4-Port",
		UnitPrice:   345000,
		Quantity:    1,
		Subtotal:    345000,
		CreatedAt:   now.Add(-5 * time.Hour),
	}).Error
	rawPayload2, _ := json.Marshal(map[string]any{"channel": "VIRTUAL_ACCOUNT_BCA", "status": "SUCCESS"})
	_ = db.Create(&model.Payment{
		ID:              uuid.New(),
		OrderID:         order2ID,
		InvoiceNumber:   inv2,
		Amount:          345000,
		PaymentMethod:   "VIRTUAL_ACCOUNT_BCA",
		Channel:         "BCA",
		Status:          "SUCCESS",
		DokuReferenceID: "TRX-DOKU-0002",
		RawPayload:      rawPayload2,
		CreatedAt:       now.Add(-5 * time.Hour),
		UpdatedAt:       now.Add(-4 * time.Hour),
	}).Error

	// Order 3: Processing & Paid via DOKU Mandiri VA (Seller preparing package)
	order3ID := uuid.New()
	ord3 := model.Order{
		ID:              order3ID,
		BuyerID:         buyer3ID,
		SellerID:        seller2ID,
		TotalAmount:     1250000,
		Status:          model.StatusProcessing,
		PaymentStatus:   model.PaymentStatusPaid,
		InvoiceNumber:   &inv3,
		PaymentURL:      &payUrl3,
		ShippingAddress: addrSurabaya,
		Notes:           &note3,
		CreatedAt:       now.Add(-24 * time.Hour),
		UpdatedAt:       now.Add(-20 * time.Hour),
	}
	_ = db.Create(&ord3).Error
	_ = db.Create(&model.OrderItem{
		ID:          uuid.New(),
		OrderID:     order3ID,
		ProductID:   prod5ID,
		ProductName: "Studio Monitor Wireless Headphones ANC",
		UnitPrice:   1250000,
		Quantity:    1,
		Subtotal:    1250000,
		CreatedAt:   now.Add(-24 * time.Hour),
	}).Error
	rawPayload3, _ := json.Marshal(map[string]any{"channel": "VIRTUAL_ACCOUNT_MANDIRI", "status": "SUCCESS"})
	_ = db.Create(&model.Payment{
		ID:              uuid.New(),
		OrderID:         order3ID,
		InvoiceNumber:   inv3,
		Amount:          1250000,
		PaymentMethod:   "VIRTUAL_ACCOUNT_MANDIRI",
		Channel:         "MANDIRI",
		Status:          "SUCCESS",
		DokuReferenceID: "TRX-DOKU-0003",
		RawPayload:      rawPayload3,
		CreatedAt:       now.Add(-24 * time.Hour),
		UpdatedAt:       now.Add(-20 * time.Hour),
	}).Error

	// Order 4: Shipped & Paid via DOKU QRIS (Waiting for Buyer to Confirm Delivery)
	order4ID := uuid.New()
	ord4 := model.Order{
		ID:              order4ID,
		BuyerID:         buyer4ID,
		SellerID:        seller1ID,
		TotalAmount:     770000,
		Status:          model.StatusShipped,
		PaymentStatus:   model.PaymentStatusPaid,
		InvoiceNumber:   &inv4,
		PaymentURL:      &payUrl4,
		ShippingAddress: addrYogyakarta,
		CreatedAt:       now.Add(-48 * time.Hour),
		UpdatedAt:       now.Add(-24 * time.Hour),
	}
	_ = db.Create(&ord4).Error
	_ = db.Create(&model.OrderItem{
		ID:          uuid.New(),
		OrderID:     order4ID,
		ProductID:   prod2ID,
		ProductName: "Wireless Ergonomic Mouse 26K DPI",
		UnitPrice:   420000,
		Quantity:    1,
		Subtotal:    420000,
		CreatedAt:   now.Add(-48 * time.Hour),
	}).Error
	_ = db.Create(&model.OrderItem{
		ID:          uuid.New(),
		OrderID:     order4ID,
		ProductID:   prod3ID,
		ProductName: "Monitor ScreenBar Eye-Care Light",
		UnitPrice:   350000,
		Quantity:    1,
		Subtotal:    350000,
		CreatedAt:   now.Add(-48 * time.Hour),
	}).Error
	rawPayload4, _ := json.Marshal(map[string]any{"channel": "QRIS", "status": "SUCCESS"})
	_ = db.Create(&model.Payment{
		ID:              uuid.New(),
		OrderID:         order4ID,
		InvoiceNumber:   inv4,
		Amount:          770000,
		PaymentMethod:   "QRIS",
		Channel:         "QRIS",
		Status:          "SUCCESS",
		DokuReferenceID: "TRX-DOKU-0004",
		RawPayload:      rawPayload4,
		CreatedAt:       now.Add(-48 * time.Hour),
		UpdatedAt:       now.Add(-47 * time.Hour),
	}).Error

	// Order 5: Completed & Paid via DOKU ShopeePay (Completed Order History)
	order5ID := uuid.New()
	ord5 := model.Order{
		ID:              order5ID,
		BuyerID:         buyer2ID,
		SellerID:        seller3ID,
		TotalAmount:     330000,
		Status:          model.StatusCompleted,
		PaymentStatus:   model.PaymentStatusPaid,
		InvoiceNumber:   &inv5,
		PaymentURL:      &payUrl5,
		ShippingAddress: addrJakarta,
		CreatedAt:       now.Add(-7 * 24 * time.Hour),
		UpdatedAt:       now.Add(-5 * 24 * time.Hour),
	}
	_ = db.Create(&ord5).Error
	_ = db.Create(&model.OrderItem{
		ID:          uuid.New(),
		OrderID:     order5ID,
		ProductID:   prod10ID,
		ProductName: "Premium Felt & Vegan Leather Desk Mat 90x40",
		UnitPrice:   185000,
		Quantity:    1,
		Subtotal:    185000,
		CreatedAt:   now.Add(-7 * 24 * time.Hour),
	}).Error
	_ = db.Create(&model.OrderItem{
		ID:          uuid.New(),
		OrderID:     order5ID,
		ProductID:   prod12ID,
		ProductName: "Solid Walnut Wood Wrist Rest for TKL",
		UnitPrice:   145000,
		Quantity:    1,
		Subtotal:    145000,
		CreatedAt:   now.Add(-7 * 24 * time.Hour),
	}).Error
	rawPayload5, _ := json.Marshal(map[string]any{"channel": "SHOPEEPAY", "status": "SUCCESS"})
	_ = db.Create(&model.Payment{
		ID:              uuid.New(),
		OrderID:         order5ID,
		InvoiceNumber:   inv5,
		Amount:          330000,
		PaymentMethod:   "SHOPEEPAY",
		Channel:         "SHOPEEPAY",
		Status:          "SUCCESS",
		DokuReferenceID: "TRX-DOKU-0005",
		RawPayload:      rawPayload5,
		CreatedAt:       now.Add(-7 * 24 * time.Hour),
		UpdatedAt:       now.Add(-7 * 24 * time.Hour),
	}).Error

	// Order 6: Cancelled & Expired Payment (Cancelled History & Stock Restored)
	order6ID := uuid.New()
	ord6 := model.Order{
		ID:              order6ID,
		BuyerID:         buyer1ID,
		SellerID:        seller3ID,
		TotalAmount:     175000,
		Status:          model.StatusCancelled,
		PaymentStatus:   model.PaymentStatusExpired,
		InvoiceNumber:   &inv6,
		PaymentURL:      &payUrl6,
		ShippingAddress: addrBandung,
		CreatedAt:       now.Add(-3 * 24 * time.Hour),
		UpdatedAt:       now.Add(-3 * 24 * time.Hour),
	}
	_ = db.Create(&ord6).Error
	_ = db.Create(&model.OrderItem{
		ID:          uuid.New(),
		OrderID:     order6ID,
		ProductID:   prod9ID,
		ProductName: "Adjustable Aluminum Laptop Stand",
		UnitPrice:   175000,
		Quantity:    1,
		Subtotal:    175000,
		CreatedAt:   now.Add(-3 * 24 * time.Hour),
	}).Error
	rawPayload6, _ := json.Marshal(map[string]any{"reason": "Payment expired after 60 minutes", "status": "EXPIRED"})
	_ = db.Create(&model.Payment{
		ID:            uuid.New(),
		OrderID:       order6ID,
		InvoiceNumber: inv6,
		Amount:        175000,
		Status:        "EXPIRED",
		RawPayload:    rawPayload6,
		CreatedAt:     now.Add(-3 * 24 * time.Hour),
		UpdatedAt:     now.Add(-3 * 24 * time.Hour),
	}).Error

	log.Printf("Successfully seeded: %d Sellers, %d Buyers, %d Products, %d Carts, 6 Sample Orders (all statuses)",
		len(sellers), len(buyers), len(products), len(buyers))
}
