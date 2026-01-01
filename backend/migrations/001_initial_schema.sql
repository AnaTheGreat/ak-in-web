CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "username" varchar(50) UNIQUE NOT NULL,
  "password_hash" varchar(255) NOT NULL,
  "is_admin" boolean DEFAULT true,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "books" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "cover_image_url" varchar(500),
  "title" varchar(200) NOT NULL,
  "author" varchar(100) NOT NULL,
  "isbn" varchar(20),
  "rating" integer NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "films" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "poster_image_url" varchar(500),
  "title" varchar(200) NOT NULL,
  "rating" integer NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "book_tags" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar(50) UNIQUE NOT NULL,
  "color" varchar(20)
);

CREATE TABLE "film_tags" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar(50) UNIQUE NOT NULL,
  "color" varchar(20)
);

CREATE TABLE "book_tag_links" (
  "book_id" uuid,
  "tag_id" uuid,
  PRIMARY KEY ("book_id", "tag_id")
);

CREATE TABLE "film_tag_links" (
  "film_id" uuid,
  "tag_id" uuid,
  PRIMARY KEY ("film_id", "tag_id")
);

ALTER TABLE "book_tag_links" ADD FOREIGN KEY ("book_id") REFERENCES "books" ("id");

ALTER TABLE "book_tag_links" ADD FOREIGN KEY ("tag_id") REFERENCES "book_tags" ("id");

ALTER TABLE "film_tag_links" ADD FOREIGN KEY ("film_id") REFERENCES "films" ("id");

ALTER TABLE "film_tag_links" ADD FOREIGN KEY ("tag_id") REFERENCES "film_tags" ("id");
