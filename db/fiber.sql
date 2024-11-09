CREATE TABLE "users" (
    id bigserial PRIMARY KEY,
    name CHARACTER VARYING(255) NOT NULL,
    username CHARACTER VARYING(255) NOT NULL,
    password TEXT NOT NULL,
    picture TEXT,
    phone character varying(255),
    two_factor_auth BOOLEAN,
    gender character varying(255),
    role character varying(255),
    is_active BOOLEAN,
    created_at DATETIME,
    updated_at DATETIME
)