-- Create the snippets table
CREATE TABLE snippets (
                          id SERIAL PRIMARY KEY, -- SERIAL provides auto-incremented integer IDs in PostgreSQL
                          title VARCHAR(100) NOT NULL,
                          content TEXT NOT NULL,
                          created TIMESTAMP NOT NULL,
                          expires TIMESTAMP NOT NULL
);

-- Create an index on the created column in the snippets table
CREATE INDEX idx_snippets_created ON snippets(created);

-- Create the users table
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL,
                       hashed_password CHAR(60) NOT NULL,
                       created TIMESTAMP NOT NULL
);

-- Add a unique constraint on the email column in the users table
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

-- Insert a sample user into the users table
INSERT INTO users (name, email, hashed_password, created) VALUES (
                                                                     'Alice Jones',
                                                                     'alice@example.com',
                                                                     '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
                                                                     '2022-01-01 10:00:00'
                                                                 );