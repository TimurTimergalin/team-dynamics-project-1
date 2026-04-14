FROM python:3.11-slim

# Build argument for the subdirectory containing the service
ARG PROJECT_DIR

# Set working directory
WORKDIR /app

# Copy shared libraries and generated protobuf code
COPY ./gen ./gen
COPY ./services ./services

# Set Python path to find modules in /app
ENV PYTHONPATH=/app:/app/gen/api/proto/python
WORKDIR /app/${PROJECT_DIR}

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Run the service
CMD ["python", "./main.py"]