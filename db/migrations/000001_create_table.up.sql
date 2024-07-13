CREATE TABLE IF NOT EXISTS auth_service_users(
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	username VARCHAR UNIQUE NOT NULL,
	email VARCHAR UNIQUE NOT NULL, 
	password_hash VARCHAR NOT NULL,
	full_name VARCHAR NOT NULL,
	bio TEXT,
	eco_points INT DEFAULT 0,
	created_at timestamp DEFAULT DEFACURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
	deleted_at timestamp
);


create table IF NOT EXISTS refresh_token(
    id uuid primary key DEFAULT gen_random_uuid(),
    user_id uuid references users(id),
    token text UNIQUE NOT NULL,
    expires_at bigint,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

