# How to start
`docker compose up -d`

# Create topic
```
docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic inventory-request
```

```
docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic inventory-result
```


# Test message
`kafka-console-producer --broker-list localhost:9092 --topic inventory-request --property "parse.key=true" --property "key.separator=:"
`

```
echo '"payload": { "products": [ { "id": "1", "name": "Product 1", "quantity": 10 }, { "id": "2", "name": "Product 2", "quantity": 20 } ], "action": "reserve", "request_id": "123" }' | kafka-console-producer --broker-list localhost:9092 --topic <topic-name> --property "parse.key=true" --property "key.separator=:" --property "key=payload"
```


echo '{"payload": {"products":[{"product_code":"1","product_name":"product1","quantity":10},{"product_code":"2","product_name":"product2","quantity":5}],"action":"HoldInventory","request_id":"12345"}' | sed 's/,\s*}/}/' | kafka-console-producer --broker-list localhost:9092 --topic inventory-request --property "parse.key=true" --property "key.separator=:"

```
$ echo '{"products":[{"id":"1","name":"product1","quantity":10},{"id":"2","name":"product2","quantity":5}],"action":"HoldInventory","request_id":"12345"}' | kafka-console-producer --broker-list localhost:9092 --topic inventory-request --property "parse.key=true" --property "key.separator=:"
```

```
kafka-console-producer --broker-list localhost:9092 --topic inventory-request --property "parse.key=true" --property "key.separator=:" --property "key=payload" --property "value={\"products\":[{\"product_id\":\"product-123\",\"quantity\":5}],\"action\":\"update\",\"request_id\":\"12345\"}"

```