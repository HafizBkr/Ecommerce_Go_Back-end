-- Table des utilisateurs/clients (avec adresse et numéro de téléphone)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_admin BOOLEAN DEFAULT false,
    points INTEGER DEFAULT 0,
    last_login TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Nouvelle colonne pour l'adresse
    address TEXT, -- Peut contenir l'adresse complète (numéro, rue, code postal, ville, pays)
    phone_number VARCHAR(20), -- Numéro de téléphone
    residence_city VARCHAR(100), -- Ville de résidence
    residence_country VARCHAR(100) -- Pays de résidence
);

-- Activer l'extension nécessaire pour générer des UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Création de la table "categories"
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Identifiant unique généré automatiquement
    nom VARCHAR(255) NOT NULL,                    -- Nom de la catégorie
    nombre_produits INT DEFAULT 0,                -- Nombre de produits dans la catégorie
    statut VARCHAR(50) DEFAULT 'actif',           -- Statut de la catégorie (par ex. actif, inactif)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Date de création
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Date de mise à jour
);

-- Table Produit
-- Modification de la table produits avec UUID comme identifiant
CREATE TABLE produits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- Utilisation d'UUID au lieu de SERIAL
    nom VARCHAR(255) NOT NULL,
    prix DECIMAL(10,2) NOT NULL,
    stock INTEGER DEFAULT 0,
    etat VARCHAR(50) NOT NULL CHECK (etat IN ('Très bon état', 'Reconditionné', 'Bon état', 'État correct')),
    photos TEXT[], -- Pour stocker les URLs des images
    categorie_id UUID REFERENCES categories(id) ON DELETE RESTRICT,
    localisation VARCHAR(255) NOT NULL, -- Pour "Paris 11ème"
    description TEXT,
    nombre_vues INTEGER DEFAULT 0,
    
    -- Ajouts utiles par rapport à votre version
    disponible BOOLEAN DEFAULT true,
    marque VARCHAR(100),
    modele VARCHAR(100),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index pour améliorer les performances
CREATE INDEX idx_produits_nom ON produits(nom);
CREATE INDEX idx_produits_categorie ON produits(categorie_id);
CREATE INDEX idx_produits_etat ON produits(etat);

-- Trigger pour mettre à jour updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_produits_updated_at
    BEFORE UPDATE ON produits
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE event_categories (
    id UUID PRIMARY KEY,
    label VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    start_time TIME,
    price DECIMAL(10, 2), -- Pour stocker le prix en CFA
    event_type_id UUID REFERENCES event_categories(id),
    available_seats INTEGER,
    image_url VARCHAR(255),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index pour améliorer les performances des recherches
CREATE INDEX idx_events_date ON events(start_date);
CREATE INDEX idx_events_type ON events(event_type_id);