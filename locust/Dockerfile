FROM python:3.11-slim

WORKDIR /mnt/locust

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY locustfile.py .

EXPOSE 8089

ENTRYPOINT [ "locust" ]
CMD [ "-f", "locustfile.py", "--host", "http://app:8080" ]