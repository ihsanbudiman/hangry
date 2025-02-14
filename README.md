
## Requirements

To run this project, you need to have the following installed:

1. [Go](https://golang.org/doc/install) version 1.23
2. [GNU Make](https://www.gnu.org/software/make/)
3. [oapi-codegen](https://github.com/deepmap/oapi-codegen)

  Install the latest version with:
  ```bash
  go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
  ```
4. [mock](https://github.com/uber-go/mock)

  Install the latest version with:
  ```bash
  go install go.uber.org/mock/mockgen@latest
  ```
5. [Docker](https://docs.docker.com/get-docker/) version 20

6. [Docker Compose](https://docs.docker.com/compose/install/) version 1.29

---

## Workflow

1. **Add to Cart**
   - Ensure a cart is created before using any promotional features.
   - Cart data is linked to user data, so a cart will be automatically created if one does not exist.
   - Refer to the API documentation for an example of how to delete an item from the cart.

2. **Create Promo**
   - You can create a promo. Refer to the API documentation or Postman file for examples.
   - You can also extend the promo date.

3. **Fetch Eligible Promo**
   - Once the cart is ready and a promo is created, you can fetch a list of eligible promos.

4. **Place Order**
   - You can place an order after adding items to the cart. Multiple promos can be applied.  
   - After four orders, a user is classified as a loyal user.  
   - There are two types of promos:  
     - **Buy X Get Y**: This promo can be redeemed if the user meets the required minimum product order. For example, buying at least three of item X grants an additional quantity of item Y.  
     - **Percentage Discount**: This promo can be redeemed if the user meets the minimum order amount specified in the promo.  
   - There are four promo categories:  
     - **All:** Available to all customers who meet the criteria.  
     - **City:** Applicable only to users in a specific city.  
     - **Loyal User:** Available exclusively to users classified as loyal.  
     - **New User:** Can be redeemed by newly registered users within their first month. After one month, this promo is no longer available.  

---

## Initiate The Project

To start working, execute:

```bash
make init
```

## Running

You can run the project using the following command:

```bash
docker compose up -d
```

The API will be accessible at [http://localhost:1323](http://localhost:1323).

## Testing

To run unit tests, execute the following command:

```bash
make test
```

## Seeding

To seed the database, use the following command:

```bash
make seed
```

*Note: Ensure the database is running first using `docker compose up -d`.*

The following data will be generated:

### Users

| Name  | Email             | City     | Is Loyal |
|-------|-------------------|----------|----------|
| Andi  | andi@example.com  | Jakarta  | true     |
| Budi  | budi@example.com  | Bandung  | false    |
| Citra | citra@example.com | Surabaya | false    |

### Products

| Name         | Price |
|--------------|-------|
| Nasi Goreng  | 25000 |
| Ayam Goreng  | 35000 |
| Es Teh Manis | 5000  |

### Promos

| Name                          | Description                                  | Segmentation | Type       | Min Order Amount | Discount Value | Max Discount Amount | Buy Product ID | Free Product ID | Buy Product Qty | Free Product Qty | Start Date | End Date                               | Max Usage Limit | Current Usage Count |
|-------------------------------|----------------------------------------------|--------------|------------|------------------|----------------|---------------------|----------------|-----------------|-----------------|------------------|------------|----------------------------------------|-----------------|---------------------|
| Test Promo Buy X Get Y        | Get one product free when you buy another    | ALL          | BUYXGETY   | 0                | 0              | 0                   | 1              | 2               | 1               | 1                | NOW()      | NOW() + INTERVAL '1 month'             | 100             | 0                   |
| Test Promo Percentage Discount| Get a percentage off your order              | CITY         | PERCENTAGE | 50000            | 15             | 10000               | NULL           | NULL            | 0               | 0                | NOW()      | NOW() + INTERVAL '1 month'             | 100             | 0                   |
| Test Promo Loyalty Discount   | Get a percentage off for loyal customers     | LOYALUSER    | PERCENTAGE | 30000            | 10             | 5000                | NULL           | NULL            | 0               | 0                | NOW()      | NOW() + INTERVAL '1 month'             | 100             | 0                   |
| Test Promo New User Discount  | Get 20% discount for new users               | NEWUSER      | PERCENTAGE | 0                | 20             | 20000               | NULL           | NULL            | 0               | 0                | NOW()      | NOW() + INTERVAL '1 month'             | 100             | 0                   |

### Promo Cities

| Promo ID | City    |
|----------|---------|
| 2        | Jakarta |
| 2        | Bandung |


