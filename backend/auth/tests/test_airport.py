import datetime
import pytest

import requests
import tests.utils.auth as auth
import tests.utils.user as user
import tests.utils.airport as airport
import tests.utils.datetime_utils as datetime_utils
     

def test_get_scedules(api_url):
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    airport_client = airport.AirportClient(api_url, user_auth.tokens_cookies)

    _, _ = airport_client.get_airports()