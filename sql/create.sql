CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(50) UNIQUE,
    email varchar(50) UNIQUE,
    full_name varchar(50),
    display_picture varchar(50)
);