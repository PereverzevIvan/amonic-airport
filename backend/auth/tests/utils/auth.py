import requests
import jwt


test_user = {
    "email": "j.doe@amonic.com",
    "password": "123"
}

test_user_wrong_password = {
    "email": "j.doe@amonic.com",
    "password": "1234"
}

test_user_wrong_email = {
    "email": "wrong_email@amonic.com",
    "password": "123"
}

class Auth:
    access_token = ""
    refresh_token = ""
    cookies = None
    # Логаут увеличивает версию токена еще на 1,
    # что должно учитываться проверяться в login запросе
    was_last_request_logout = False 

    def __init__(self, api_url):
        self.api_url = api_url

    def login_user(self, user_to_login=test_user, do_assertion = True):
        response = requests.post(f'{self.api_url}/login', 
                                 json=user_to_login)
        
        if do_assertion:
            self.__assert_tokens_response(response)
            if self.was_last_request_logout:
                self.__assert_refresh_response(response)
        
        self.access_token = response.cookies.get('access-token')
        self.refresh_token = response.cookies.get('refresh-token')
        self.tokens_cookies = response.cookies
        
        self.was_last_request_logout = False

        return response

    def refresh_user_tokens(self, do_assertion = True):
        response = requests.get(f'{self.api_url}/refresh', 
                                cookies=self.tokens_cookies)
        
        if do_assertion:
            self.__assert_refresh_response(response)
        
        self.access_token = response.cookies.get('access-token')
        self.refresh_token = response.cookies.get('refresh-token')
        self.tokens_cookies = response.cookies

        return response

    def logout_user(self, do_assertion = True):
        response = requests.get(f'{self.api_url}/logout', 
                                cookies=self.tokens_cookies)
        
        if do_assertion:
            self.__assert_logout_response(response)
        
        self.was_last_request_logout = True

        return response
    
    def __assert_tokens_response(self, response: requests.Response):
        assert response.status_code == 200

        new_access_token = response.cookies.get('access-token')
        new_refresh_token = response.cookies.get('refresh-token')
        # Проверка cookie токенов
        assert new_access_token != None
        assert new_refresh_token != None

        # Проверка параметров cookie токенов
        access_token_cookie = next((c for c in response.cookies if c.name == "access-token"), None)
        assert access_token_cookie is not None
        assert not access_token_cookie.has_nonstandard_attr("HttpOnly")
        
        refresh_token_cookie = next((c for c in response.cookies if c.name == "refresh-token"), None)
        assert refresh_token_cookie is not None
        assert refresh_token_cookie.has_nonstandard_attr("HttpOnly")

    def __assert_refresh_response(self, response: requests.Response):
        # Проверка ответа и новых токенов
        self.__assert_tokens_response(response)

        # Сравнение новых и старых токенов
        new_access_token = response.cookies.get('access-token')
        new_refresh_token = response.cookies.get('refresh-token')
        
        assert new_access_token != self.access_token
        assert new_refresh_token != self.refresh_token

        # Проверка версии токенов
        old_tokens_version = jwt.decode(self.access_token, options={"verify_signature": False})["ver"]
        new_tokens_version = jwt.decode(new_access_token, options={"verify_signature": False})["ver"]
        assert old_tokens_version + 1 + self.was_last_request_logout == new_tokens_version
    
    def __assert_logout_response(self, response: requests.Response):
        assert response.status_code == 200

        assert response.cookies.get('access-token') == None
        assert response.cookies.get('refresh-token') == None

        