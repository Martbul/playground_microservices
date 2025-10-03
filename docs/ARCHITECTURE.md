This is an API that uses microservices architecture
It has api-gateway and 3 services - client, auth and product service

api-gateway is the client for the APIs services(auth,client and products)
the APIs services(auth,client and products) act as a server for the api-gateway

                  ┌─────────────────┐
                  │  Client Service │ (port 8083)
                  └────────┬────────┘
                           │ HTTP
                           ▼
                  ┌─────────────────┐
                  │  API Gateway    │ (port 8080)
                  └───┬─────────┬───┘
                      │         │
              HTTP    │         │ HTTP
                      ▼         ▼
        ┌─────────────────┐   ┌─────────────────┐
        │ Auth Service    │   │ Product Service │
        │ (port 8081)     │   │ (port 8082)     │
        └──────┬──────────┘   └────────┬────────┘
               │                       │
          SQL  │                       │ SQL
               ▼                       ▼
        ┌────────────────────────────────┐
        │        PostgreSQL (5432)       │
        │   auth_db, product_db inside   │
        └────────────────────────────────┘
