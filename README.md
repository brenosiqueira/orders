Orders
=====


How to test
-----------
Create order
```bash
# POST http://localhos/orders
curl -X POST -H "Content-Type: application/json" -d '{
  "id": "123",
  "number": "12312",
  "reference": "ABC1",
  "notes": "Nota",
  "price": 1000
}' "http://localhos/orders"

```

Create order item

```bash
# POST http://localhos/orders/:order_id/items
curl -X POST -H "Content-Type: application/json" -d '{
    "sku": "b79f2752-5bfb-11e6-8b77-86f30ca893d3",
    "unit_price": 1000,
    "quantity": 2
}' "http://localhost/orders/df144f65-34e8-4450-9324-a9e1f6d5178b/items"

```

```bash
# POST http://localhos/orders/:order_id/transactions
curl -X POST -H "Content-Type: application/json"  -d '{
    "external_id": "10",
    "amount": 1000,
    "type": "PAYMENT",
    "card_brand": "VISA",
    "card_bin": "1402",
    "card_last": "3211"
}' "http://localhost/orders/df144f65-34e8-4450-9324-a9e1f6d5178b/transactions"
```