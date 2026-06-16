-- Ekstensi untuk UUID jika belum aktif
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ==========================================
-- Tabel Users
-- ==========================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    phone_number VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_users_email ON users(email);

-- ==========================================
-- Tabel Categories
-- ==========================================
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE UNIQUE INDEX idx_categories_slug ON categories(slug);

-- ==========================================
-- Tabel Charities (Campaigns)
-- ==========================================
CREATE TABLE charities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE RESTRICT,
    category_id UUID REFERENCES categories(id) ON DELETE RESTRICT,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    target_amount DECIMAL(15,2) NOT NULL,
    current_amount DECIMAL(15,2) DEFAULT 0 NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    cover_image_url VARCHAR(255),
    contact_person VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    organization_name VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending' NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_charities_slug ON charities(slug);
CREATE INDEX idx_charities_status_date ON charities(status, created_at DESC);
CREATE INDEX idx_charities_category_id ON charities(category_id);

-- ==========================================
-- Tabel Charity Images
-- ==========================================
CREATE TABLE charity_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    charity_id UUID REFERENCES charities(id) ON DELETE CASCADE,
    image_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_charity_images_charity_id ON charity_images(charity_id);

-- ==========================================
-- Tabel Donations
-- ==========================================
CREATE TABLE donations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    charity_id UUID REFERENCES charities(id) ON DELETE RESTRICT,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Null jika donasi sebagai guest
    amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' NOT NULL,
    payment_method VARCHAR(100),
    payment_reference VARCHAR(255),
    anonymous BOOLEAN DEFAULT FALSE NOT NULL,
    message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_donations_charity_id_status ON donations(charity_id, status);
CREATE INDEX idx_donations_user_id ON donations(user_id);

-- ==========================================
-- Tabel Articles
-- ==========================================
CREATE TABLE articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID REFERENCES users(id) ON DELETE RESTRICT,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    cover_image_url VARCHAR(255),
    is_published BOOLEAN DEFAULT FALSE NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_articles_slug ON articles(slug);
CREATE INDEX idx_articles_published ON articles(is_published, published_at DESC);
