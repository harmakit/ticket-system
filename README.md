# Ticket System - Microservices Architecture

## Overview

This project consists of four microservices built using Go that together manage a ticket booking system. The microservices communicate with each other via gRPC, with Kafka handling event-driven interactions and PostgreSQL for storage. Each service is designed to handle a specific domain, and the system is set up for easy scaling and maintenance.

The services are:

1. **Event Service**: Manages events and their associated tickets.
2. **Checkout Service**: Handles user carts and the checkout process.
3. **Booking Service**: Manages seat reservations and availability.
4. **Notification Service**: Sends notifications when orders are updated.

Features
- Each microservice has a clean structure with separation of concerns, making the project easy to maintain and extend.
- Logical separation between layers like service, handler, repository, and API definitions.
- Each service uses Docker multi-stage builds to optimize the size of the production images, improving deployment efficiency.
- The project utilizes separate Docker Compose files for development with a debugger and production environments.
- Each microservice comes with a `Makefile` that includes a help section describing the main commands.
- All microservices communicate via gRPC, which provides fast and efficient inter-service communication with protocol buffers ensuring well-defined APIs.
- The `booking` service uses a cron-style scheduler to handle jobs like checking for expired bookings and restoring stock, ensuring up-to-date availability for users.
- Versioning is enforced in the `proto` files (such as `v1` for API and message definitions), allowing backward compatibility and clear evolution of the system.


### Microservices Summary

#### 1. Event Service
- Provides available events and tickets.
- Allows event creation and retrieval of specific event details. 
- gRPC methods:
  - CreateEvent: Creates a new event.
  - GetEvent: Fetches event details.
  - ListEvents: Lists available events.
  - GetTicket: Retrieves details of a ticket for a specific event.

#### 2. Checkout Service
- Manages user carts and the entire checkout process.
- Adds tickets to carts, modifies cart items, initiates booking requests, and handles order creation.
- Ensures that orders are placed and paid for successfully, with expired bookings being deleted.
- gRPC methods:
  - `GetOrder`: Retrieves details of a specific order.
  - `ListOrders`: Lists all user orders.
  - `GetUserCart`: Retrieves the current cart for a user.
  - `AddToCart`: Adds tickets to the user's cart.
  - `UpdateCart`: Updates items in the user's cart.
  - `PlaceOrder`: Places a new order and triggers the booking service.
  - `MarkOrderAsPaid`: Marks an order as paid.
  - `CancelOrder`: Cancels an order.

#### 3. Booking Service
- Manages the reservation of seats for events.
- Ensures availability of seats, monitors expired bookings, and restores stock when necessary.
- Supports the creation of stock for events.
- gRPC methods:
  - `CreateBooking`: Creates a seat reservation.
  - `GetBookings`: Retrieves bookings for a given order.
  - `ExpireBookings`: Expires old bookings and restores the stock.
  - `DeleteOrderBookings`: Deletes bookings associated with an order.
  - `CreateStock`: Creates a stock of available seats for an event.
  - `GetStocks`: Retrieves available stocks for events.
  - `GetStock`: Retrieves the stock for a specific event.
  - `DeleteStock`: Deletes a stock entry.

#### 4. Notification Service
- Sends notifications to users regarding the status of their orders (such as successful purchases).
- Uses Kafka as a message broker to handle notifications asynchronously.

## Project Structure

Each service is designed to be self-contained and is run as a separate container using Docker. The project uses `docker-compose` to bring up all services in a development environment. PostgreSQL is used as the database, with each service managing its own schema.

### Communication

- **gRPC** is used for synchronous communication between the services, providing type-safe and efficient APIs.
- **Kafka** handles asynchronous messaging, especially for notifications and other background tasks.

### Database

- **PostgreSQL** is the database for all services. Each service has its own schema, and migrations are applied via the `Makefile`.

### Local Development

Use `docker-compose` to run the services locally.

- Start the services:
  ```bash
  make start
  ```

- Stop the services:
  ```bash
  make stop
  ```

Each service has its own `Makefile` to simplify running tasks like database migrations or generating gRPC code.
