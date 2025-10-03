-- PostgreSQL doesn't support IF NOT EXISTS for CREATE DATABASE
-- Instead, we use a DO block to check if the database exists first

-- Create auth_db if it doesn't exist
SELECT 'CREATE DATABASE auth_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'auth_db')\gexec

-- Create product_db if it doesn't exist
SELECT 'CREATE DATABASE product_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'product_db')\gexec

-- Note: The default 'postgres' database already exists and is used by the gateway