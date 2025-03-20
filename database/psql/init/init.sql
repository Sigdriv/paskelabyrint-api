-- -- Drop the existing "Paskelabyrint" database if it exists
-- DROP DATABASE IF EXISTS Paskelabyrint;

-- -- Create the "Paskelabyrint" database
-- CREATE DATABASE Paskelabyrint;

-- -- Ensure the database exists before switching to it
-- \connect Paskelabyrint;

-- Create the "users" table
CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Sessions (
    id VARCHAR(200) PRIMARY KEY UNIQUE NOT NULL,
    email VARCHAR(100) NOT NULL REFERENCES Users(email),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
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
INSERT INTO Users (name, email, password) 
VALUES ('Sigurd Drivstuen', 'sigdriv06@gmail.com', '$2a$12$fQ0mzw.K0vSyZLKIR7n3Ou2XpN/CH6kTARyPa/MmvhC0KaE3hlXGy');
