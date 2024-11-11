CREATE TABLE "users" (
    id bigserial PRIMARY KEY,
    name CHARACTER VARYING(255) NOT NULL,
    email CHARACTER VARYING(255) NOT NULL,
    username CHARACTER VARYING(255) NOT NULL,
    password TEXT NOT NULL,
    address TEXT,
    picture TEXT,
    phone character varying(255),
    two_factor_auth BOOLEAN,
    gender character varying(255),
    role character varying(255),
    is_active BOOLEAN,
    created_at DATE,
    updated_at DATE
)

CREATE TABLE otp (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT,
    "otp_code" VARCHAR(6),
    "is_verified" BIGINT,
    "is_expired" SMALLINT
);