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

def test_logout_user_session(api_url):
    user_auth = auth.Auth(api_url, auth.test_user)
    user_auth.login_user()
    
    user_session_client = user_session.UserSessionClient(api_url, user_auth.tokens_cookies)
    
    _, user_sessions = user_session_client.get_current_user_sessions()
    login_session = user_sessions[0]

    assert login_session.logout_at == None
    
    # Выход из аккаунта -> в последнюю сессию записывается время выхода 
    user_auth.logout_user()

    # Входим в аккаунт чтобы получить доступ к списку сессий
    user_auth.login_user()
    # Обновляем куки для авторизации
    user_session_client.cookies = user_auth.tokens_cookies

    _, user_sessions = user_session_client.get_current_user_sessions()
    assert len(user_sessions) >= 1
    logout_session = user_sessions[1]


    assert login_session.id == logout_session.id

    assert logout_session.logout_at != None
    assert login_session.login_at == logout_session.login_at

def test_login_2_times_to_check_invalid_logout_marked(api_url):
    user_auth = auth.Auth(api_url, auth.test_user)
    # Авторизовываемся 2 раза, 
    # чтобы предпоследняя сессия пометилась как невалидная
    user_auth.login_user()
    user_auth.login_user()

    user_session_client = user_session.UserSessionClient(api_url, user_auth.tokens_cookies)
    
    _, user_sessions = user_session_client.get_current_user_sessions()
    assert len(user_sessions) >= 1
    
    login_session = user_sessions[0]
    assert login_session.invalid_logout_reason == None or login_session.invalid_logout_reason == ""
    
    invalid_logout_session = user_sessions[1]

    # Проверка, что предпоследняя сессия помечена как невалидная
    assert invalid_logout_session.invalid_logout_reason == "undefined"
    assert invalid_logout_session.logout_at == None
    assert invalid_logout_session.crash_reason_type == 0
    
    
def test_update_invalid_logout_session(api_url):
    user_auth = auth.Auth(api_url, auth.test_user)

    user_auth = auth.Auth(api_url, auth.test_user)
    # Авторизовываемся 2 раза, 
    # чтобы предпоследняя сессия пометилась как невалидная
    user_auth.login_user()
    user_auth.login_user()

    user_session_client = user_session.UserSessionClient(api_url, user_auth.tokens_cookies)
    
    _, user_sessions = user_session_client.get_current_user_sessions()
    assert len(user_sessions) >= 1
    invalid_logout_session = user_sessions[1]

    # Проверка, что предпоследняя сессия помечена как невалидная
    assert invalid_logout_session.invalid_logout_reason == "undefined"
    assert invalid_logout_session.logout_at == None
    assert invalid_logout_session.crash_reason_type == 0

    # Устанавливаем невалидную причину выхода
    invalid_logout_session.invalid_logout_reason = "test"
    invalid_logout_session.crash_reason_type = 1

    _ = user_session_client.update_user_session(invalid_logout_session)
    _, new_user_sessions = user_session_client.get_current_user_sessions()
    new_invalid_logout_session = new_user_sessions[1]

    assert new_invalid_logout_session.invalid_logout_reason == "test"
    assert new_invalid_logout_session.logout_at == None
    assert new_invalid_logout_session.crash_reason_type == 1 