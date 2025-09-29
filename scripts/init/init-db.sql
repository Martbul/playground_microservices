CREATE DATABASE IF NOT EXISTS auth_db;
CREATE DATABASE IF NOT EXISTS product_db;

\c auth_db
-- optional: tables or demo data

\c product_db
-- optional: tables or demo data

GRANT ALL PRIVILEGES ON DATABASE auth_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE product_db TO postgres;
