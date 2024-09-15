import pytest
import requests
import tests.utils.datetime_utils as datetime_utils

class UserSession:
    id = None
    user_id = None
    login_at = None
    logout_at = None
    invalid_logout_reason = None
    crash_reason_type = None

    def __init__(self, json: dict):
        self.id = json['id']
        self.user_id = json['user_id']
        if json['login_at'] != None:
            self.login_at = datetime_utils.parse(json['login_at'])
        if json['logout_at'] != None:
            self.logout_at = datetime_utils.parse(json['logout_at'])

        self.invalid_logout_reason = json['invalid_logout_reason']
        self.crash_reason_type = json['crash_reason_type']
    
    @staticmethod
    def parse_user_sessions(json: list):
        return [UserSession(session) for session in json]

class UserSessionClient:
    cookies = None

    def __init__(self, api_url, jwt_cookies):
        self.api_url = api_url
        self.cookies = jwt_cookies
    
    # def get

    def get_current_user_sessions(self, do_assertion = True):
        response = requests.get(f'{self.api_url}/user-sessions', 
                                cookies=self.cookies)
        
        if do_assertion:
            assert response.status_code == 200
        
        user_sessions = UserSession.parse_user_sessions(response.json())
        assert len(user_sessions) > 0


        return response, user_sessions

    