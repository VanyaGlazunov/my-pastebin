import random
import string
from locust import HttpUser, task, between

def random_string(length=100):
    """Generates string for paste"""
    letters = string.ascii_letters + string.digits
    return ''.join(random.choice(letters) for i in range(length))

class PastebinUser(HttpUser):
    wait_time = between(1, 2)
    
    def on_start(self):
        """Init state for new user"""
        self.created_paste_ids = []

    @task(1)
    def create_paste(self):
        """
        Simulates new paste creation
        """
        expires_in = random.choice(["10m", "1h", "1d"])
        syntax = random.choice(["text", "python", "go", "json"])
        payload = {
            "content": f"Random content from Locust: {random_string(256)}",
            "expires_in": expires_in,
            "syntax": syntax
        }
        headers = {"Content-Type": "application/json"}

        with self.client.post("/api/v1/paste", json=payload, headers=headers, name="/api/v1/paste", catch_response=True) as response:
            if response.status_code == 201:
                try:
                    paste_id = response.json().get("id")
                    if paste_id:
                        self.created_paste_ids.append(paste_id)
                except response.JSONDecodeError:
                    response.failure("Response is not valid JSON")

    @task(5)
    def read_paste(self):
        """
        Simulates reading random created paste
        """
        if not self.created_paste_ids:
            return

        random_id = random.choice(self.created_paste_ids)

        self.client.get(f"/api/v1/paste/{random_id}", name="/api/v1/paste/[id]")