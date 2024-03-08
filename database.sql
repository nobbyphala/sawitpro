/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.user_profile (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	full_name varchar NOT NULL,
	phone_number varchar NOT NULL,
	"password" varchar(60) NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	success_count int8 NOT NULL DEFAULT 0,
	CONSTRAINT user_profile_un UNIQUE (phone_number),
	CONSTRAINT user_table_pk PRIMARY KEY (id)
);
