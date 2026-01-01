CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "username" varchar(50) UNIQUE NOT NULL,
  "password_hash" varchar(255) NOT NULL,
  "is_admin" boolean DEFAULT true,
  "created_at" timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "books" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "cover_image_url" varchar(500),
  "title" varchar(200) NOT NULL,
  "author" varchar(100) NOT NULL,
  "isbn" varchar(20),
  "rating" integer NOT NULL CHECK (rating >= 1 AND rating <= 5),
  "created_at" timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "films" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "poster_image_url" varchar(500),
  "title" varchar(200) NOT NULL,
  "rating" integer NOT NULL CHECK (rating >= 1 AND rating <= 5),
  "created_at" timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "book_tags" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar(50) UNIQUE NOT NULL,
  "color" varchar(20)
);

CREATE TABLE IF NOT EXISTS "film_tags" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar(50) UNIQUE NOT NULL,
  "color" varchar(20)
);

CREATE TABLE IF NOT EXISTS "book_tag_links" (
  "book_id" uuid,
  "tag_id" uuid,
  PRIMARY KEY ("book_id", "tag_id"),
  FOREIGN KEY ("book_id") REFERENCES "books" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("tag_id") REFERENCES "book_tags" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "film_tag_links" (
  "film_id" uuid,
  "tag_id" uuid,
  PRIMARY KEY ("film_id", "tag_id"),
  FOREIGN KEY ("film_id") REFERENCES "films" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("tag_id") REFERENCES "film_tags" ("id") ON DELETE CASCADE
);

-- Create admin user (password: admin123)
INSERT INTO users (username, password_hash, is_admin) 
VALUES ('admin', '$2a$10$rKZLVqF5qH5qH5qH5qH5qOyJGKvZ8K.vZ8K.vZ8K.vZ8K.vZ8K.vZ8', true)
ON CONFLICT (username) DO NOTHING;

-- Insert sample books
INSERT INTO books (title, author, rating) VALUES
  ('The Problems Of Philosophy', 'Bertrand Russell', 5),
  ('Pokemon Adventures, Vol. 1', 'Hidenori Kusaka', 4),
  ('American Psycho', 'Bret Easton Ellis', 5)
ON CONFLICT DO NOTHING;

-- Insert sample films
INSERT INTO films (title, rating) VALUES
  ('The Matrix', 5),
  ('Blade Runner 2049', 5)
ON CONFLICT DO NOTHING;
