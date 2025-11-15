# Pet Clinic API (Go + PostgreSQL)

A minimal REST-style API for a Pet Clinic built in Go, backed by PostgreSQL.

---

## Overview

* **Language:** Go
* **Database:** PostgreSQL
* **Driver:** [github.com/lib/pq](https://github.com/lib/pq)
* **Module:** `petclinic`

---

## Quick Start

### 1. Install Prerequisites

* [Go](https://go.dev/dl/)
* [PostgreSQL](https://www.postgresql.org/download/)

---

### 2. Clone the Project

```bash
git clone <your-repo-url>
cd petclinic
```

---

### 3. Set Up the Database

Run the following SQL commands in `psql` as a superuser (e.g. `postgres`):

```sql
-- Create DB user and password
CREATE ROLE petuser WITH LOGIN PASSWORD 'petpass';

-- Create database
CREATE DATABASE petclinic OWNER petuser;

-- Connect to DB
\c petclinic

-- Create schema (tables and sample data)

-- OWNERS
CREATE TABLE owners (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    address TEXT
);

-- PETS
CREATE TABLE pets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    species VARCHAR(50) NOT NULL,
    breed VARCHAR(50),
    birth_date DATE,
    owner_id INT REFERENCES owners(id) ON DELETE CASCADE
);

-- VETS
CREATE TABLE vets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    specialization VARCHAR(100)
);

-- VISITS
CREATE TABLE visits (
    id SERIAL PRIMARY KEY,
    pet_id INT REFERENCES pets(id) ON DELETE CASCADE,
    vet_id INT REFERENCES vets(id) ON DELETE SET NULL,
    visit_date DATE NOT NULL DEFAULT CURRENT_DATE,
    description TEXT
);

-- Indexes
CREATE INDEX idx_pets_owner_id ON pets(owner_id);
CREATE INDEX idx_visits_pet_id ON visits(pet_id);
CREATE INDEX idx_visits_vet_id ON visits(vet_id);

-- Sample data
INSERT INTO owners (name, phone, address) VALUES
('John Doe', '9876543210', '123 Pet Street'),
('Jane Smith', '9988776655', '45 Bark Avenue');

INSERT INTO pets (name, species, breed, birth_date, owner_id) VALUES
('Buddy', 'Dog', 'Labrador', '2021-02-15', 1),
('Mittens', 'Cat', 'Siamese', '2020-08-12', 2);

INSERT INTO vets (name, specialization) VALUES
('Dr. Emily Clark', 'Surgery'),
('Dr. Ravi Patel', 'Dermatology');

INSERT INTO visits (pet_id, vet_id, visit_date, description) VALUES
(1, 1, '2025-01-12', 'Annual vaccination'),
(2, 2, '2025-03-20', 'Skin allergy checkup');

-- Privileges for the app user
GRANT USAGE ON SCHEMA public TO petuser;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO petuser;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO petuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO petuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO petuser;
```

---

### 4. Configure the App Connection

Update the database connection string in `db.go` if needed:

```
host=localhost port=5432 user=petuser password=petpass dbname=petclinic sslmode=disable
```

---

### 5. Run the Application

```bash
go get github.com/lib/pq
go mod tidy
go run .
```

**Expected output:**

```
Connected to PostgreSQL!
Server running at :8080
```

---

## API Routes

**Base URL:** `http://localhost:8080`

---

### Owners

* **GET** `/owners` — Returns list of owners

  ```bash
  curl http://localhost:8080/owners
  ```

* **GET** `/owners/id?id={id}` — Returns a single owner by ID

  ```bash
  curl "http://localhost:8080/owners/id?id=1"
  ```

* **POST** `/owners` — Creates a new owner

  ```bash
  curl -X POST http://localhost:8080/owners \
    -H "Content-Type: application/json" \
    -d '{"name":"Alice","phone":"1234567890","address":"1 Main St"}'
  ```

---

### Pets

* **GET** `/pets` — Returns list of pets

  ```bash
  curl http://localhost:8080/pets
  ```

* **GET** `/pets/id?id={id}` — Returns a single pet by ID

  ```bash
  curl "http://localhost:8080/pets/id?id=1"
  ```

* **POST** `/pets` — Creates a pet (dates use RFC3339 format `YYYY-MM-DDT00:00:00Z`)

  ```bash
  curl -X POST http://localhost:8080/pets \
    -H "Content-Type: application/json" \
    -d '{"name":"Rex","species":"Dog","breed":"Beagle","birth_date":"2023-05-01T00:00:00Z","owner_id":1}'
  ```

---

### Vets

* **GET** `/vets` — Returns list of all vets

  ```bash
  curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/vets
  ```

* **POST** `/vets` — Create a new vet
  
  ```bash
  curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_JWT_TOKEN" \
    -d '{"name":"Dr. Smith", "specialization":"Surgery"}' \
    http://localhost:8080/vets
  ```

* **GET** `/vets/id?id={id}` — Returns a single vet by ID

  ```bash
  curl -H "Authorization: Bearer YOUR_JWT_TOKEN" "http://localhost:8080/vets/id?id=1"
  ```

* **PUT** `/vets/id?id={id}` — Update a vet

  ```bash
  curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_JWT_TOKEN" \
    -d '{"name":"Dr. Smith", "specialization":"Advanced Surgery"}' \
    "http://localhost:8080/vets/id?id=1"
  ```

* **DELETE** `/vets/id?id={id}` — Delete a vet

  ```bash
  curl -X DELETE -H "Authorization: Bearer YOUR_JWT_TOKEN" "http://localhost:8080/vets/id?id=1"
  ```

### Visits

* **GET** `/visits` — Returns list of visits

  ```bash
  curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/visits
  ```

* **GET** `/visits/id?id={id}` — Returns a single visit by ID

  ```bash
  curl -H "Authorization: Bearer YOUR_JWT_TOKEN" "http://localhost:8080/visits/id?id=1"
  ```

* **POST** `/visits` — Creates a visit

  ```bash
  curl -X POST http://localhost:8080/visits \
    -H "Content-Type: application/json" \
    -d '{"pet_id":1,"vet_id":1,"visit_date":"2025-01-12T00:00:00Z","description":"Checkup"}'
  ```

---

## Notes

* Default credentials in `db.go`:

  ```
  user: petuser
  password: petpass
  database: petclinic
  ```
* For production, use environment variables for credentials.
* If you encounter permission errors, re-run the GRANT statements or transfer ownership:

  ```sql
  ALTER TABLE table_name OWNER TO petuser;
  ```

