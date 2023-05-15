import os
import requests

def get_base_url() -> str:
    return os.environ.get("MODRINTH_API_URL") or "https://api.modrinth.com/v2"

def get_resource(path: str) -> str:
    return requests.get(f"{get_base_url()}/{path}").text
