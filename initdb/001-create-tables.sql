-- Creat the Host Information table
CREATE TABLE IF NOT EXISTS host_information (
    host_id SERIAL PRIMARY KEY,
    hostname VARCHAR(255) NOT NULL,
    ip_address VARCHAR(15),
    os VARCHAR(50),
    os_version VARCHAR(50),
    hardware_specs JSONB,
    agent_version VARCHAR(50),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create the User Table
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL
);

-- Create the Command Logs table
CREATE TABLE IF NOT EXISTS command_logs (
    command_id SERIAL PRIMARY KEY,
        host_id INT NOT NULL,
        command TEXT NOT NULL,
        timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        response_status VARCHAR(50),
        response TEXT,
        FOREIGN KEY (host_id) REFERENCES host_information (host_id)
);