# protos-todo

This project implements a microservices architecture consisting of three main components: **API Gateway**, **Auth Service**, and **Todo Service**.
![Project Screenshot](diagram.jpg)

## Architecture

### 1. API Gateway
- Serves as a single entry point for clients, providing a convenient interface for interacting with microservices.
- Communicates with **Auth Service** and **Todo Service** via gRPC, ensuring high performance and low latency.

### 2. Auth Service
- Responsible for user authentication and authorization.

### 3. Todo Service
- Manages todo items and performs CRUD operations.
- Has its own database for storing task data.

## Technologies

- **HTTP & Gin**:
    - Used to create a RESTful API, allowing easy handling of HTTP requests.

- **gRPC**:
    - Facilitates efficient communication between services (API Gateway, Auth Service, and Todo Service).

- **GORM**:
    - An Object-Relational Mapping (ORM) library that simplifies database interactions.

- **PostgreSQL**:
    - Utilized as a relational database for both Todo Service and Auth Service.

- **Docker & Docker Compose**:
    - Docker is used for containerizing services, while Docker Compose manages multiple services together, enabling easy deployment of PostgreSQL as an image.

## Database Structure

- **Todo Service Database**:
    - Stores data about tasks, including titles, descriptions, and completion status.

- **Auth Service Database**:
    - Maintains user information, including credentials for authentication.

## Getting Started

To get started with the project, clone the repository and run the following commands to set up the services using Docker:

```bash
docker-compose up