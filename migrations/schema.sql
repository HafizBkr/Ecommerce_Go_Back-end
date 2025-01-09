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


CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    nom VARCHAR(100) NOT NULL,
    sous_categories TEXT[], -- Comme montré dans l'image 2
    nombre_produits INTEGER DEFAULT 0,
    statut VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Produit
CREATE TABLE produits (
    id SERIAL PRIMARY KEY,
    nom VARCHAR(255) NOT NULL,
    prix DECIMAL(10,2) NOT NULL,
    stock INTEGER DEFAULT 0,
    etat VARCHAR(50) NOT NULL, -- Pour "Très bon état", "Reconditionné", etc.
    photos TEXT[], -- Pour stocker les URLs des images comme montré dans le formulaire
    categorie_id INTEGER REFERENCES categories(id),
    localisation VARCHAR(255), -- Pour "Paris 11ème" comme dans le formulaire
    description TEXT,
    nombre_vues INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);