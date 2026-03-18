CREATE TABLE IF NOT EXISTS books (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	author TEXT NOT NULL,
	isbn TEXT,
	year INT,
	rating NUMERIC(3,1),
	cover TEXT
);

CREATE TABLE IF NOT EXISTS library_tags (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS book_tags (
	book_id INT REFERENCES books(id) ON DELETE CASCADE,
	tag_id INT REFERENCES library_tags(id) ON DELETE CASCADE,
	PRIMARY KEY (book_id, tag_id)
);