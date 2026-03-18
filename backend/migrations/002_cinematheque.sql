CREATE TABLE IF NOT EXISTS movies (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	director TEXT NOT NULL,
	year INT,
	rating NUMERIC(3,1),
	cover TEXT
);

CREATE TABLE IF NOT EXISTS cinematheque_tags (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS movie_tags (
	movie_id INT REFERENCES movies(id) ON DELETE CASCADE,
	tag_id INT REFERENCES cinematheque_tags(id) ON DELETE CASCADE,
	PRIMARY KEY (movie_id, tag_id)
);