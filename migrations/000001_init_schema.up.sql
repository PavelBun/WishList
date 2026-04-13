CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE wishlists (
                           id SERIAL PRIMARY KEY,
                           user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           title VARCHAR(255) NOT NULL,
                           description TEXT,
                           event_date DATE NOT NULL,
                           access_token UUID NOT NULL DEFAULT gen_random_uuid(),
                           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                           UNIQUE(access_token)
);

CREATE TABLE items (
                       id SERIAL PRIMARY KEY,
                       wishlist_id INTEGER NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       product_link VARCHAR(512),
                       priority INTEGER NOT NULL DEFAULT 0,
                       is_booked BOOLEAN NOT NULL DEFAULT FALSE,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);
CREATE INDEX idx_items_wishlist_id ON items(wishlist_id);
CREATE INDEX idx_wishlists_access_token ON wishlists(access_token);
CREATE INDEX idx_items_is_booked ON items(is_booked);