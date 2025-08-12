# Use a Python Alpine image as it's lightweight.
FROM python:3.9-alpine

# Install FFmpeg, which is a system-level dependency for video transcoding.
# We also install build-base for any C extensions that might be needed by pip packages.
RUN apk update && apk add --no-cache ffmpeg build-base

# Set the working directory in the container
WORKDIR /app

# Copy the dependencies file first to leverage Docker's layer caching.
COPY requirements.txt .

# Install the Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application source code
COPY . .

# This command will be executed when the container starts.
# It runs the main consumer script.
CMD ["python", "main.py"]
