# Use a slim Python image for a smaller final image size
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Copy the dependencies file first to leverage Docker's layer caching
COPY requirements.txt .

# Install the Python dependencies
RUN pip install --no-cache-dir --trusted-host pypi.python.org -r requirements.txt

# Copy the application source code into the container
COPY . .

# Expose the port that the Uvicorn server will run on
EXPOSE 3004

# The command to run the application using Uvicorn.
# The default port for this service will be 3004.
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "3004"]
