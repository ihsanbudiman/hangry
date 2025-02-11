package db

import (
	"context"
	"database/sql"
	"fmt"
	"hangry/domain/models"
	"hangry/repository"

	"gorm.io/gorm"
)

type promoRepostory struct {
	db *gorm.DB
}

// GetPromoByUserCart implements repository.PromoRepository.
func (r *promoRepostory) GetPromoByUserCart(ctx context.Context, tx *gorm.DB, input repository.GetPromoByUserCartInput) ([]models.Promo, int64, error) {

	baseQuery := `
		with 
		user_data as (
				select 
						*
				from users u 
				where u.id = @userId
		),
		items as (
				select
						p.price,
						ci.*
				from cart_items ci 
				join carts c on c.id = ci.cart_id
				join products p on p.id = ci.product_id 
				where c.user_id  = @userId
		),
		summary as (
				select
						sum(i.price * i.quantity) total
				from items i
		)
		select
				%[1]v
		from promos p
		where (
						(
								p."type" = 'PERCENTAGE_DISCOUNT' 
								and (select total from summary) >= p.min_order_amount 
						)
						or (
								p."type" = 'BUY_X_GET_Y_FREE' 
								and (select 
										coalesce(i.quantity,0)
										from items i
										where i.product_id = p.buy_product_id 
								) >= p.buy_product_qty 
						)
				)
				and (
						p.segmentation = 'ALL'
						or (
								p.segmentation = 'CITY' 
								and lower((select city from user_data)) in (select lower(city) from promo_cities pc where pc.promo_id = p.id)
						)
						or (
								p.segmentation = 'LOYAL_USER'
								and (select is_loyal from user_data) = true
						)
						or (
								p.segmentation = 'NEW_USER'
								and (select created_at from user_data) > (current_timestamp - interval '1 month')
						)
				)
		%[2]v
		%[3]v
		%[4]v;
	`

	selectFields := `
		p.*
	`

	additionalCondition := ""
	if len(input.PromoIds) > 0 {
		additionalCondition += " and p.id in @promoIds"
	}

	if input.IsAvailable != nil && *input.IsAvailable {
		additionalCondition += " and p.start_date <= current_timestamp and p.end_date >= current_timestamp and p.current_usage_count < p.max_usage_limit"
	}

	limit := ""
	if input.Page != nil && input.PerPage != nil {
		offset := (*input.Page - 1) * *input.PerPage
		limit = fmt.Sprintf(" limit %d offset %d", *input.PerPage, offset)
	}

	groupBy := " group by p.id"

	query := fmt.Sprintf(baseQuery, selectFields, additionalCondition, groupBy, limit)
	db := tx
	if db == nil {
		db = r.db.WithContext(ctx)
	}

	var promos []models.Promo

	if err := db.Raw(query, sql.Named("userId", input.Cart.UserID), sql.Named("promoIds", input.PromoIds)).Scan(&promos).Error; err != nil {
		return nil, 0, err
	}

	countQuery := fmt.Sprintf(baseQuery, "count(*)", additionalCondition, "", "")
	var count int64
	if err := db.Raw(countQuery, sql.Named("userId", input.Cart.UserID), sql.Named("promoIds", input.PromoIds)).Scan(&count).Error; err != nil {
		return nil, 0, err
	}

	return promos, count, nil
}

// GetPromoByPromoID implements repository.PromoRepository.
func (r *promoRepostory) GetPromoByPromoID(ctx context.Context, tx *gorm.DB, promoID uint) (models.Promo, error) {
	db := tx
	if db == nil {
		db = r.db.WithContext(ctx)
	}

	var promo models.Promo
	if err := db.Where("id = ?", promoID).First(&promo).Error; err != nil {
		return promo, err
	}

	return promo, nil
}

func (r *promoRepostory) SaveCities(ctx context.Context, tx *gorm.DB, promoID uint, cities []string) error {
	db := tx
	if db == nil {
		db = r.db.WithContext(ctx)
	}

	pc := make([]models.PromoCity, len(cities))
	for i, city := range cities {
		pc[i] = models.PromoCity{
			PromoID: promoID,
			City:    city,
		}
	}

	return db.Create(&pc).Error
}

// CreatePromo implements PromoRepository.
func (p *promoRepostory) Save(ctx context.Context, tx *gorm.DB, promo *models.Promo) error {
	db := p.db
	if tx != nil {
		db = tx
	}

	return db.Save(promo).Error
}

func NewPromoRepository(db *gorm.DB) repository.PromoRepository {
	return &promoRepostory{
		db: db,
	}
}
