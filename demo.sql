
-- ==========================
-- DATABASE: petclinic
-- ==========================

-- 1️⃣ OWNERS TABLE
CREATE TABLE owners (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    address TEXT
);

-- 2️⃣ PETS TABLE
CREATE TABLE pets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    species VARCHAR(50) NOT NULL,
    breed VARCHAR(50),
    birth_date DATE,
    owner_id INT REFERENCES owners(id) ON DELETE CASCADE
);

-- 3️⃣ VETS TABLE
CREATE TABLE vets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    specialization VARCHAR(100)
);

-- 4️⃣ VISITS TABLE
CREATE TABLE visits (
    id SERIAL PRIMARY KEY,
    pet_id INT REFERENCES pets(id) ON DELETE CASCADE,
    vet_id INT REFERENCES vets(id) ON DELETE SET NULL,
    visit_date DATE NOT NULL DEFAULT CURRENT_DATE,
    description TEXT
);

-- 5️⃣ OPTIONAL: INDEXES (For performance)
CREATE INDEX idx_pets_owner_id ON pets(owner_id);
CREATE INDEX idx_visits_pet_id ON visits(pet_id);
CREATE INDEX idx_visits_vet_id ON visits(vet_id);

-- 6️⃣ OPTIONAL: SAMPLE DATA (You can skip this if you want a clean DB)
INSERT INTO owners (name, phone, address)
VALUES
('John Doe', '9876543210', '123 Pet Street'),
('Jane Smith', '9988776655', '45 Bark Avenue');

INSERT INTO pets (name, species, breed, birth_date, owner_id)
VALUES
('Buddy', 'Dog', 'Labrador', '2021-02-15', 1),
('Mittens', 'Cat', 'Siamese', '2020-08-12', 2);

INSERT INTO vets (name, specialization)
VALUES
('Dr. Emily Clark', 'Surgery'),
('Dr. Ravi Patel', 'Dermatology');

INSERT INTO visits (pet_id, vet_id, visit_date, description)
VALUES
(1, 1, '2025-01-12', 'Annual vaccination'),
(2, 2, '2025-03-20', 'Skin allergy checkup');

-- USERS TABLE FOR AUTH
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
