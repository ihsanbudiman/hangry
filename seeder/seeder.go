package seeder

import (
	"fmt"
	"hangry/constants"
	"hangry/domain/models"
	"time"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	tx := db.Begin()

	defer func() {
		tx.Rollback()
	}()

	if err := SeedUsers(tx); err != nil {
		panic(err)
	}

	if err := SeedProducts(tx); err != nil {
		panic(err)
	}

	if err := SeedPromos(tx); err != nil {
		panic(err)
	}

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}

}

// SeedUsers inserts sample user data into the database
func SeedUsers(db *gorm.DB) error {
	users := []models.User{
		{Name: "Andi", Email: "andi@example.com", City: "Jakarta", IsLoyal: true},
		{Name: "Budi", Email: "budi@example.com", City: "Bandung", IsLoyal: false},
		{Name: "Citra", Email: "citra@example.com", City: "Surabaya", IsLoyal: false},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

// SeedProducts inserts sample product data into the database
func SeedProducts(db *gorm.DB) error {
	products := []models.Product{
		{Name: "Nasi Goreng", Price: 25000}, // Harga dalam Rupiah
		{Name: "Ayam Goreng", Price: 35000},
		{Name: "Es Teh Manis", Price: 5000},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			return err
		}
	}

	return nil
}

// seed promos
func SeedPromos(db *gorm.DB) error {

	var buyProductID uint = 1
	var freeProductID uint = 2
	limit := 100

	promos := []models.Promo{
		// Promo Buy X Get Y
		{
			Name:              "Test Promo Buy X Get Y",
			Description:       "Get one product free when you buy another",
			Segmentation:      constants.PROMOSEGMENTATIONALL,
			Type:              constants.PROMOTYPEBUYXGETY,
			MinOrderAmount:    0,
			DiscountValue:     0,
			MaxDiscountAmount: 0,
			BuyProductID:      &buyProductID,
			FreeProductID:     &freeProductID,
			BuyProductQty:     1,
			FreeProductQty:    1,
			StartDate:         time.Now(),
			EndDate:           time.Now().AddDate(0, 1, 0),
			MaxUsageLimit:     &limit,
			CurrentUsageCount: 0,
		},
		// Promo Percentage Discount
		{
			Name:              "Test Promo Percentage Discount",
			Description:       "Get a percentage off your order",
			Segmentation:      constants.PROMOSEGMENTATIONCITY,
			Type:              constants.PROMOTYPEPERCENTAGE,
			MinOrderAmount:    50000,
			DiscountValue:     15, // 15% discount
			MaxDiscountAmount: 10000,
			BuyProductID:      nil,
			FreeProductID:     nil,
			BuyProductQty:     0,
			FreeProductQty:    0,
			StartDate:         time.Now(),
			EndDate:           time.Now().AddDate(0, 1, 0),
			MaxUsageLimit:     &limit,
			CurrentUsageCount: 0,
			PromoCities: []models.PromoCity{
				{City: "Jakarta"},
				{City: "Bandung"},
			},
		},
		// Promo Percentage But for Loyalty User
		{
			Name:              "Test Promo Loyalty Discount",
			Description:       "Get a percentage off for loyal customers",
			Segmentation:      constants.PROMOSEGMENTATIONLOYALUSER,
			Type:              constants.PROMOTYPEPERCENTAGE,
			MinOrderAmount:    30000,
			DiscountValue:     10, // 10% discount
			MaxDiscountAmount: 5000,
			BuyProductID:      nil,
			FreeProductID:     nil,
			BuyProductQty:     0,
			FreeProductQty:    0,
			StartDate:         time.Now(),
			EndDate:           time.Now().AddDate(0, 1, 0),
			MaxUsageLimit:     &limit,
			CurrentUsageCount: 0,
			PromoCities: []models.PromoCity{
				{City: "Jakarta"},
				{City: "Bandung"},
			},
		},
		// Promo Percentage But for New User
		{
			Name:              "Test Promo New User Discount",
			Description:       "Get 20% discount for new users",
			Segmentation:      constants.PROMOSEGMENTATIONNEWUSER,
			Type:              constants.PROMOTYPEPERCENTAGE,
			MinOrderAmount:    0,
			DiscountValue:     20, // 20% discount
			MaxDiscountAmount: 20000,
			BuyProductID:      nil,
			FreeProductID:     nil,
			BuyProductQty:     0,
			FreeProductQty:    0,
			StartDate:         time.Now(),
			EndDate:           time.Now().AddDate(0, 1, 0),
			MaxUsageLimit:     &limit,
			CurrentUsageCount: 0,
		},
	}

	for _, promo := range promos {
		if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&promo).Error; err != nil {
			return err
		}
	}

	fmt.Println("Promos seeded successfully")

	return nil
}
