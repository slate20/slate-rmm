#!/bin/bash

# This script is used to set up the environment for the application.

# Update the apt package list
sudo apt update && sudo apt upgrade -y

# Install Docker
sudo apt install -y docker.io

# Install Docker Compose
sudo apt install -y docker-compose

# Restart Docker service
sudo service docker restart

# Navigate to the server directory
cd server

# Set permissions for the server binary
chmod +x slatermm

# Create a systemd service file for the server
echo "[Unit]
Description=SlaterMM Server

[Service]
ExecStart=$(pwd)/slatermm
Restart=always
User=$(whoami)
Group=$(whoami)
Environment=PATH=$/usr/bin:/usr/local/bin
Environment=GO_ENV=production
WorkingDirectory=$(pwd)

[Install]
WantedBy=multi-user.target" | sudo tee /etc/systemd/system/slatermm.service

# Reload the systemd daemon
sudo systemctl daemon-reload

# Run Docker Services
docker-compose up -d

# Start the server
sudo systemctl start slatermm

# Enable the server to start on boot
sudo systemctl enable slatermm