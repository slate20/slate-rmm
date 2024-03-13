#!/bin/bash

# Call the original entrypoint script in the background
echo "Running original entrypoint script..."
/docker-entrypoint.sh &

# Wait for the CheckMK service to start
echo "Waiting for CheckMK service to start..."
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' -L http://localhost:5000/${CMK_SITE_ID}/check_mk/)" != "200" ]]; do
    sleep 1
done
echo "CheckMK service started"

# Run the expect script
echo "Running expect script..."
expect automation_pw.exp $AUTOMATION_SECRET
echo "Expect script finished"

# Keep the container running
tail -f /dev/null