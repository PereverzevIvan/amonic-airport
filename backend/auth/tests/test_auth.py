import pytest
import requests

import tests.utils.auth as auth


def test_login(api_url):
    user_auth = auth.Auth(api_url, auth.test_user)
    user_auth.login_user()


def test_login_wrong_password(api_url):
    user_auth = auth.Auth(api_url, auth.test_user_wrong_password)

    response = user_auth.login_user(do_assertion=False)

    # Проверка cookie токенов
    assert response.cookies.get("access-token") is None
    assert response.cookies.get("refresh-token") is None

    assert response.status_code == 401
    assert response.text == "Wrong email or password"


def test_login_wrong_email(api_url):
    user_auth = auth.Auth(api_url, auth.test_user_wrong_email)

    response = user_auth.login_user(do_assertion=False)

    # Проверка cookie токенов
    assert response.cookies.get("access-token") is None
    assert response.cookies.get("refresh-token") is None

    assert response.status_code == 401
    assert response.text == "Wrong email or password"


def test_login_no_email(api_url):
    user_auth = auth.Auth(api_url, user_to_login={"password": "123"})

    response = user_auth.login_user(do_assertion=False)

    # Проверка cookie токенов
    assert response.cookies.get("access-token") is None
    assert response.cookies.get("refresh-token") is None

    assert response.status_code == 400
    assert response.text == "email is required"


def test_login_no_password(api_url):
    user_auth = auth.Auth(api_url, user_to_login={"email": "j.doe@amonic.com"})

    response = user_auth.login_user(do_assertion=False)

    # Проверка cookie токенов
    assert response.cookies.get("access-token") is None
    assert response.cookies.get("refresh-token") is None

    assert response.status_code == 400
    assert response.text == "password is required"


def test_refresh_and_version_increase(api_url):
    user_auth = auth.Auth(api_url)

    user_auth.login_user()
    user_auth.refresh_user_tokens()


def test_logout_and_version_increase(api_url):
    user_auth = auth.Auth(api_url)

    user_auth.login_user()

    user_auth.logout_user()

    user_auth.login_user()


def test_logout_twice(api_url):
    user_auth = auth.Auth(api_url)

    user_auth.login_user()

    user_auth.logout_user()

    response = user_auth.logout_user(False)
    assert response.status_code == 401
    assert response.text == auth.response_text_invalid_token_version
