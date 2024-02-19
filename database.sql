CREATE TYPE public.enum_user_registration_type AS ENUM (
	'phone_number'
);

CREATE TABLE public."user" (
	id varchar(512) NOT NULL,
	phone_number varchar(20) NULL,
	registration_type public.enum_user_registration_type NOT NULL,
	full_name varchar(255) NOT NULL,
	CONSTRAINT user_pkey PRIMARY KEY (id),
  	CONSTRAINT phone_number UNIQUE (phone_number)
);

CREATE TABLE public.user_credential (
	user_id varchar(512) NOT NULL,
	"method" varchar(32) DEFAULT 'bcrypt'::character varying NOT NULL,
	hash text NOT NULL,
	salt varchar(1024) NULL,
	CONSTRAINT user_credential_pkey PRIMARY KEY (user_id, method),
	CONSTRAINT user_credential_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);