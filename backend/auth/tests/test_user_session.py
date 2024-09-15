import pytest

import requests
import tests.utils.auth as auth
import tests.utils.user_session as user_session
import tests.utils.datetime_utils as datetime_utils

def test_login_new_user_session(api_url):
    time_before_login = datetime_utils.now() 

    user_auth = auth.Auth(api_url, auth.test_user)
    user_auth.login_user()

    time_after_login = datetime_utils.now(True)

    user_session_client = user_session.UserSessionClient(api_url, user_auth.tokens_cookies)
    
    _, user_sessions = user_session_client.get_current_user_sessions()
    current_session = user_sessions[0]

    assert len(user_sessions) >= 1
    assert current_session.logout_at == None

    assert time_before_login <= current_session.login_at 
    assert current_session.login_at <= time_after_login 

    
    
    

