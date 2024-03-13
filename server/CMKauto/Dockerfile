# Use official CheckMK image as base
FROM checkmk/check-mk-raw:latest

# Install expect
RUN apt update && apt install -y expect

# Copy the automation script and entrypoint script
COPY automation_pw.exp /automation_pw.exp
COPY entrypoint.sh /entrypoint.sh

# Make the scripts executable
RUN chmod +x /automation_pw.exp /entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/bin/bash", "/entrypoint.sh"]