
CREATE TABLE USER  (
	id INT NOT NULL UNIQUE AUTO_INCREMENT,
	username VARCHAR (127) NOT NULL UNIQUE,
	password VARCHAR (127) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);
CREATE TABLE FILE_META_INFO (
	id INT NOT NULL UNIQUE AUTO_INCREMENT,
	url VARCHAR (512) NOT NULL UNIQUE,
	bucket_name VARCHAR (256) NOT NULL,
	password VARCHAR (127) NOT NULL,
	file_type VARCHAR (64) NOT NULL,
	file_size INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);
CREATE TABLE ARTIST (
	id INT NOT NULL UNIQUE AUTO_INCREMENT,
	stage_name VARCHAR (127) NOT NULL,
	user_id INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	CONSTRAINT fk_artist_user FOREIGN KEY (user_id)  
	REFERENCES USER(id)  
	ON DELETE CASCADE  
	ON UPDATE CASCADE
);

CREATE TABLE GENRE_CD (
	genre_cd VARCHAR (4) NOT NULL,
	genre_desc VARCHAR (128) NOT NULL,
	PRIMARY KEY (genre_cd)
);
CREATE TABLE ALBUM (
	id INT NOT NULL UNIQUE AUTO_INCREMENT,
	title VARCHAR(25) NOT NULL,
	artist_id INT NOT NULL,
	genre_cd VARCHAR (4) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	CONSTRAINT fk_artist FOREIGN KEY (artist_id)  
	REFERENCES ARTIST(id)  
	ON DELETE CASCADE  
	ON UPDATE CASCADE,
	CONSTRAINT fk_album_genre FOREIGN KEY (genre_cd)  
	REFERENCES GENRE_CD(genre_cd)
	ON DELETE CASCADE  
	ON UPDATE CASCADE
);
CREATE TABLE TRACK (
	id INT NOT NULL UNIQUE AUTO_INCREMENT,
	track_name VARCHAR (127) NOT NULL,
	artist_id INT NOT NULL,
	album_id INT NOT NULL,
	stream_count INT,
	file_meta_id INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	CONSTRAINT fk_track_artist FOREIGN KEY (artist_id)  
	REFERENCES ARTIST(id)  
	ON DELETE CASCADE
	ON UPDATE CASCADE,
	CONSTRAINT fk_track_album FOREIGN KEY (album_id)  
	REFERENCES ALBUM(id)  
	ON DELETE CASCADE  
	ON UPDATE CASCADE,
	CONSTRAINT fk_track_file FOREIGN KEY (file_meta_id)  
	REFERENCES FILE_META_INFO(id)  
	ON DELETE CASCADE  
	ON UPDATE CASCADE
);

