package database

import (
	"fmt"
)

func schema() string {
	return fmt.Sprint(`
	CREATE TABLE IF NOT EXISTS clients (
	  id BIGSERIAL,
	  type varchar(25) NOT NULL,
	  redirect_uris varchar(255) NULL,
	  default_redirect_uri varchar(255) NULL,
	  allowed_grant_types varchar(255) NULL,
	  PRIMARY KEY (id)
	);

	CREATE TABLE IF NOT EXISTS access_tokens (
	  id BIGSERIAL,
	  client_id BIGINT NOT NULL,
	  expires timestamp NOT NULL DEFAULT (NOW() + interval '1 month'),
	  refresh_token varchar(55) NULL,
	  scopes varchar(255) NOT NULL,
	  secret varchar(55) NOT NULL,
	  user_id BIGINT NOT NULL,
	  redirect_uri varchar(255) NULL,
	  PRIMARY KEY (id)
	);

	CREATE TABLE IF NOT EXISTS users (
	  id BIGSERIAL,
	  name varchar(255) NOT NULL,
	  email varchar(255) NOT NULL UNIQUE,
	  password varchar(255) NOT NULL,
	  created_at timestamp NOT NULL DEFAULT NOW(),
	  PRIMARY KEY (id)
	);

	CREATE TABLE IF NOT EXISTS files (
	  id BIGSERIAL,
	  name varchar(255) NOT NULL,
	  path varchar(255) NOT NULL,
	  ext varchar(25) NOT NULL,
	  algorithms varchar(15) NOT NULL,
	  user_id BIGINT NOT NULL,
	  created_at timestamp NOT NULL DEFAULT NOW(),
	  PRIMARY KEY (id)
	);`)
}
