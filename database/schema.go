package database

import (
	"fmt"
)

func schema(password string) string {
	return fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS users (
	  id BIGSERIAL,
	  email varchar(255) NOT NULL UNIQUE,
	  name varchar(255) NOT NULL,
	  password varchar(255) NOT NULL,
	  created_at timestamp NOT NULL DEFAULT NOW(),
	  PRIMARY KEY (id)
	);

	CREATE TABLE IF NOT EXISTS files (
	  id BIGSERIAL,
	  path varchar(255) NOT NULL,
	  ext varchar(25) NOT NULL,
	  algorithms varchar(15) NOT NULL,
	  user_id BIGINT NOT NULL,
	  PRIMARY KEY (id)
	);

	DELETE FROM users WHERE email = 'test@user.com';
	INSERT INTO users (email, name, password) VALUES ('test@user.com', 'John', '%s');`, password)
}
