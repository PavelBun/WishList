CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email TEXT NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE wishlists (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           title TEXT NOT NULL,
                           description TEXT NOT NULL DEFAULT '',
                           event_date DATE NOT NULL,
                           access_token UUID NOT NULL DEFAULT gen_random_uuid(),
                           created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                           updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                           UNIQUE(access_token)
);

CREATE TABLE items (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       wishlist_id UUID NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
                       title TEXT NOT NULL,
                       description TEXT NOT NULL DEFAULT '',
                       product_link TEXT NOT NULL DEFAULT '',
                       priority INTEGER NOT NULL DEFAULT 3 CHECK (priority BETWEEN 1 AND 5),
                       is_booked BOOLEAN NOT NULL DEFAULT FALSE,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);
CREATE INDEX idx_items_wishlist_id ON items(wishlist_id);
CREATE INDEX idx_wishlists_access_token ON wishlists(access_token);
CREATE INDEX idx_items_is_booked ON items(is_booked);