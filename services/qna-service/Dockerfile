# Use a slim Python image for a smaller final image size
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Copy the dependencies file first to leverage Docker's layer caching
COPY requirements.txt .

# Install the Python dependencies
# --no-cache-dir reduces image size
# --trusted-host is sometimes needed in certain network environments
RUN pip install --no-cache-dir --trusted-host pypi.python.org -r requirements.txt

# Copy the application source code into the container
COPY . .

# Expose the port that the Uvicorn server will run on
EXPOSE 3003

# The command to run the application using Uvicorn.
# --host 0.0.0.0 makes the server accessible from outside the container.
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "3003"]
