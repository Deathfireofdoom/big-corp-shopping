# BIG-CORP-SHOPPING
## Overview
This project serves as a demonstration of the power of Go's concurrency and the capabilities of building high-performance distributed systems by combining Go's native concurrency support with technologies like Kafka and Redis.

The repository encompasses the entire microservice architecture, with infrastructure residing in the `/infrastructure` directory and microservices residing in `/micro-services`.

## Services
Below is a description of each service that is currently avaiable.

## Rest API

The Rest API serves as the bridge between the outside world and the micro-service architecture. While it could be classified as one of the micro-services itself, its primary purpose is to provide easy endpoints for frontend developers to integrate a shopping-cart feature into their online shop.

The Rest API supports basic CRUD operations, including adding and removing products, getting the current cart, and placing an order. While the real magic happens in the cart-service and the inventory-service, there are some cool features here too.

When adding a product to the cart, the Rest API addresses the potential problems that arise in online shopping carts, such as ensuring that the product is in stock and avoiding race conditions when two customers attempt to add the same product simultaneously.

To address these issues, the normal HTTP call is upgraded to a websocket to allow for asynchronous communication while still providing responses. Here's how it works:

1. The client makes a GET-request to `/cart/update-product`, sending information such as `userID`, `productCode`, and `quantity`.
2. The Rest API converts this information into a `cart-request` and sends it via a `channel` to the `kafkaCommunicator`, which publishes requests and listens on requests-responses.
3. The connection is upgraded to a websocket, which starts listening on the `response-channel` for the response from the request sent in step 2.
4. When the response arrives in the `response-channel`, it is published to the websocket, allowing the client to know whether the product was added successfully or not.

In addition to the `/cart` endpoint, the Rest API also has two other endpoints:
- `/test` - functions to test the health of the Rest API, such as availability.
- `/migration` - functions to initialize the database by loading the `inventory` and creating necessary views and tables.


## Cart service
The cart-service is responsible for managing the shopping carts of users, both keeping the state of the cart and altering it, such as adding and removing products.

To avoid false promises, such as adding a product that was just snatched by another user, the cart-service uses the concept of `pending-requests`. A `pending-request` is a request for either adding or removing a product from the cart that has not yet been processed by the `inventory-service`.

Adding a product to the cart follows this flow:

1. The cart-service receives a `cart-request` from Kafka, published by the `rest-api`. The Kafka message is published in a channel that the `cart-service` listens on.
2. The cart-service handles this request by doing the following two things: adding a `pending-request` to the cart and publishing an `inventory-request` to the Kafka-bridge, which is later handled by the `inventory-service`.
3. When a response with the same `requestID` as the one published arrives from the `inventory-service`, via Kafka, if successful, the cart-service will update the cart and publish a `cart-request-response` for the `rest-api`.

_When removing a product, the request will be applied to the cart directly and return a response to the rest-api. When the inventory-request is then handled, no additional response will be published. This is because removing a product should not be dependent on the inventory-service._

Since the cart-service works asynchronously, each `cart` is stored in a `cartHandle` that has a `mutex` lock. This way, we can avoid any races. Every time a change is being made to the cart-object, the mutex-lock is checked out.



## Inventory service
The `inventory-service` plays a crucial role in ensuring that a user cannot add an out-of-stock product to their cart. When the `cart-service` receives a request to add a product to a cart, it publishes an `inventory-request` to Kafka. The `inventory-service` listens for these requests and checks the product's availability before placing a hold on it, if possible. Once the inventory status has been determined, the `inventory-service` sends the result back to the `cart-service` via Kafka.

The `inventory-service` relies on a database that contains two primary tables and one primary view:

`inventory - table`: This table stores the current inventory and has columns such as `product_code` and `quantity`.

`hold - table`: This table stores all current holds and has columns such as `product_code`, `user_id`, `created_at`, and `hold_quantity`. When a new hold is created, the `inventory-service` adds a row to this table.

`available_products - view`: This view shows the current available products, which is calculated as the inventory quantity minus the hold quantity.

### Avoiding Negative Inventory Counts due to Concurrent Holds

To prevent negative inventory counts caused by concurrent writes to the `hold` table, the `inventory-service` employs distributed locks with Redis. If the requested hold quantity is greater than five or the available inventory count after the hold is less than five, the `inventory-service` requires the workers to check out the Redis lock for that product code. This approach enables concurrent writes for high-stock products while avoiding negative inventory counts for low-stock products.

## Infrastructure
Below is a description of the infrastructure used for this microservice architecture.

### Kafka
Kafka serves as the backbone and heart of this microservice architecture, which is logical given that the essence of a microservice architecture is to decouple all components with the aid of an event bus.

Although there are numerous other options, I have always been in contact with Kafka during my tenure as a Data Engineer and Data Ops due to its proficiency in handling big data. Hence, I opted to use Kafka. In the future, I intend to integrate ksqlDB to enable real-time analytics, but that is a task for another day.

### Postgres
I am not sure if I need to justify this selection, but I used Postgres to store information about the inventory because it is naturally tabular in format, necessitating a SQL database. The decision to use Postgres over other SQL databases was simply because I am biased and prefer Postgres.

### Redis
Redis was added to the mix primarily due to the distributed mutex locks it provides out of the box, rather than as a cache solution.

At the project's outset, I relied solely on local mutexes from the sync package in Go to prevent scenarios where different workers attempted to reserve the same product at the same time, resulting in a negative inventory count. This solution worked as long as only one instance of the inventory service was running. However, it did not function if I wanted to create a scalable and distributed system, which I did.

Therefore, by integrating Redis into the infrastructure, the system can now avoid race conditions on a distributed level.

