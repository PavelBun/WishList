DROP INDEX IF EXISTS idx_unique_item_name_per_wishlist;
DROP INDEX IF EXISTS idx_items_is_booked;
DROP INDEX IF EXISTS idx_wishlists_access_token;
DROP INDEX IF EXISTS idx_items_wishlist_id;
DROP INDEX IF EXISTS idx_wishlists_user_id;

DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS pgcrypto;