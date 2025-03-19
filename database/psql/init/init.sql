-- -- Drop the existing "Paskelabyrint" database if it exists
-- DROP DATABASE IF EXISTS Paskelabyrint;

-- -- Create the "Paskelabyrint" database
-- CREATE DATABASE Paskelabyrint;

-- -- Ensure the database exists before switching to it
-- \connect Paskelabyrint;

-- Create the "users" table
CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    email VARCHAR(120) NOT NULL,
    count_participants INT NOT NULL,
    youngest_participant_age INT,
    oldest_participant_age INT,
    team_name VARCHAR(120) NOT NULL UNIQUE,
    created_by_id INT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert a sample user into the users table
INSERT INTO Users (username, email, password_hash) 
VALUES ('testuser', 'testuser@example.com', 'hashed_password_here');
