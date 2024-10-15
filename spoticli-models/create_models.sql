CREATE TABLE USER  (
	id INT NOT NULL UNIQUE AUTO_INCREMENT,
	username VARCHAR (127) NOT NULL UNIQUE,
	password VARCHAR (127) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);
CREATE TABLE ARTIST (
	id INT NOT NULL UNIQUE AUTO_INCREMENT,
	stage_name VARCHAR (127) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);
CREATE TABLE ALBUBM (
	id INT NOT NULL UNIQUE AUTO_INCREMENT,
	title VARCHAR(25) NOT NULL,
	artist_id INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);
