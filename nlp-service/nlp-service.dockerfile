FROM python:3.11-slim-bullseye

WORKDIR /app

COPY requirements.txt .

RUN pip3 install -r requirements.txt

COPY . .

EXPOSE 80

CMD ["python3", "main.py"]
