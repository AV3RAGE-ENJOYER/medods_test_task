-- +goose Up

-- USERS

CREATE TABLE public.users (
    "id" BIGINT PRIMARY KEY,
    "email" TEXT NOT NULL
);

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

ALTER TABLE ONLY public.users ALTER COLUMN "id" SET DEFAULT nextval('public.users_id_seq'::regclass);

-- AUTH

CREATE TABLE public.auth (
    "id" BIGINT PRIMARY KEY,
    "user_id" BIGINT,
    "password_hash" TEXT NOT NULL,

    FOREIGN KEY("user_id") REFERENCES public.users("id")
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE SEQUENCE public.auth_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.auth_id_seq OWNED BY public.auth.id;

ALTER TABLE ONLY public.auth ALTER COLUMN "id" SET DEFAULT nextval('public.auth_id_seq'::regclass);


-- REFRESH TOKENS

CREATE TABLE public.refresh_tokens (
    "id" BIGINT PRIMARY KEY,
    "refresh_token_hash" TEXT NOT NULL,
    "expires_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "last_ip" TEXT NOT NULL
);

CREATE SEQUENCE public.refresh_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN "id" SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);

-- INSERT TEST DATA

WITH add_user AS (
    INSERT INTO public.users("email")
    VALUES ('admin@gmail.com')
    RETURNING id)

INSERT INTO public.auth("user_id", "password_hash")
SELECT "id", '$2a$12$7TAiwtMzHZg49781OZwni.CTTeBIWKYhjkNIh/1uL8MdRK9RFMwmK' FROM add_user;

-- +goose Down

DROP TABLE public.auth;
DROP TABLE public.refresh_tokens;
DROP TABLE public.users;