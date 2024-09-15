import pytest
import requests

def test_hello(api_url):
    response = requests.get(f'{api_url}/hello')

    assert response.status_code == 200
    # assert False, f"Response: {response.text} important"