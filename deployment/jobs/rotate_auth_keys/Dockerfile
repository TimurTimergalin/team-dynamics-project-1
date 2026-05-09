FROM python:3.12-slim

WORKDIR /app

COPY jobs/rotateAuthKeys/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY jobs/rotateAuthKeys/rotate.py .

CMD ["python", "rotate.py"]
