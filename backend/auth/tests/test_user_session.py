import pytest
import requests
import jwt

# test_user = {
#     "email": "j.doe@amonic.com",
#     "password": "123"
# }

# test_user_wrong_password = {
#     "email": "j.doe@amonic.com",
#     "password": "1234"
# }

# test_user_wrong_email = {
#     "email": "wrong_email@amonic.com",
#     "password": "123"
# }

# def login_user(api_url, user_to_login):
#     response = requests.post(f'{api_url}/login', json=user_to_login)
    
#     assert response.status_code == 200
    
#     # Проверка cookie токенов
#     assert response.cookies.get('access-token') != None
#     assert response.cookies.get('refresh-token') != None

#     # Проверка параметров cookie токенов
#     access_token_cookie = next((c for c in response.cookies if c.name == "access-token"), None)
#     assert access_token_cookie is not None
#     assert not access_token_cookie.has_nonstandard_attr("HttpOnly")
    
#     refresh_token_cookie = next((c for c in response.cookies if c.name == "refresh-token"), None)
#     assert refresh_token_cookie is not None
#     assert refresh_token_cookie.has_nonstandard_attr("HttpOnly")
    
#     return response


# def test_login(api_url):
#     login_user(api_url, test_user)
    
    

# def test_login_wrong_password(api_url):
#     response = requests.post(f'{api_url}/login', json=test_user_wrong_password)
    
#     # Проверка cookie токенов
#     assert response.cookies.get('access-token') == None
#     assert response.cookies.get('refresh-token') == None

#     assert response.status_code == 401
#     assert response.text == "Wrong email or password"

# def test_login_wrong_email(api_url):
#     response = requests.post(f'{api_url}/login', json=test_user_wrong_email)
    
#     # Проверка cookie токенов
#     assert response.cookies.get('access-token') == None
#     assert response.cookies.get('refresh-token') == None

#     assert response.status_code == 401
#     assert response.text == "Wrong email or password"

# def test_login_no_email(api_url):
#     response = requests.post(f'{api_url}/login', json={
#         "password": "123"
#     })
    
#     # Проверка cookie токенов
#     assert response.cookies.get('access-token') == None
#     assert response.cookies.get('refresh-token') == None

#     assert response.status_code == 400
#     assert response.text == "email is required"

# def test_login_no_password(api_url):
#     response = requests.post(f'{api_url}/login', json={
#         "email": "j.doe@amonic.com"
#     })
    
#     # Проверка cookie токенов
#     assert response.cookies.get('access-token') == None
#     assert response.cookies.get('refresh-token') == None

#     assert response.status_code == 400
#     assert response.text == "password is required"

# def test_refresh_and_version_increase(api_url):
#     response = login_user(api_url, test_user)

#     access_token1 = response.cookies.get('access-token')
#     refresh_token1 = response.cookies.get('refresh-token')

#     refresh_response = requests.get(f'{api_url}/refresh', cookies=response.cookies)
#     assert refresh_response.status_code == 200

#     access_token2 = refresh_response.cookies.get('access-token')
#     refresh_token2 = refresh_response.cookies.get('refresh-token')
#     assert access_token1 != access_token2
#     assert refresh_token1 != refresh_token2

#     # Проверка версии
#     tokens_version1 = jwt.decode(access_token1, options={"verify_signature": False})["ver"]
#     tokens_version2 = jwt.decode(access_token2, options={"verify_signature": False})["ver"]
#     assert tokens_version1+1 == tokens_version2

# def test_logout_and_version_increase(api_url):
#     login_1_response = login_user(api_url, test_user)

#     access_token1 = login_1_response.cookies.get('access-token')
#     refresh_token1 = login_1_response.cookies.get('refresh-token')

#     # Выход
#     logout_response = requests.get(f'{api_url}/logout', cookies=login_1_response.cookies)
#     assert logout_response.status_code == 200
#     assert logout_response.cookies.get('access-token') == None
#     assert logout_response.cookies.get('refresh-token') == None

#     response_login_2 = login_user(api_url, test_user)
#     access_token2 = response_login_2.cookies.get('access-token')
#     refresh_token2 = response_login_2.cookies.get('refresh-token')
#     assert access_token1 != access_token2
#     assert refresh_token1 != refresh_token2

#     # Проверка версии
#     tokens_version1 = jwt.decode(access_token1, options={"verify_signature": False})["ver"]
#     tokens_version2 = jwt.decode(access_token2, options={"verify_signature": False})["ver"]
#     assert tokens_version1+2 == tokens_version2