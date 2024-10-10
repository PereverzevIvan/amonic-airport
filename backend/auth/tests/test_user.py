import pytest

import requests
import tests.utils.auth as auth
import tests.utils.user as user
import tests.utils.datetime_utils as datetime_utils

# 
def delete_user_by_email(mysql, user_email):
    mysql.cursor().execute(
        f"DELETE FROM `users` WHERE `email` = '{user_email}'"
    )
    mysql.commit()
     


def test_get_user(api_url):
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    user_client = user.UserClient(
        api_url, user_auth.tokens_cookies
    )

    _, user_1 = user_client.get_user(1)

    assert user_1.id == 1
    assert user_1.role_id == 1
    assert user_1.office_id == 1

    assert user_1.email == "j.doe@amonic.com"
    assert user_1.first_name == "John"
    assert user_1.last_name == "Doe"
    
    assert user_1.birthday != None
    assert user_1.active == True

# Создание
def test_create_user(api_url, mysql_conn):
    delete_user_by_email(mysql_conn, user.test_new_user.email)
    
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    user_client = user.UserClient(api_url, user_auth.tokens_cookies)

    _, created_user = user_client.create_user(user.test_new_user)

def test_try_non_admin_user_create_new_user(api_url, mysql_conn):
    delete_user_by_email(mysql_conn, user.test_new_user.email)
    
    user_auth = auth.Auth(api_url, auth.test_user)
    user_auth.login_user()

    user_client = user.UserClient(api_url, user_auth.tokens_cookies)

    response, created_user = user_client.create_user(user.test_new_user, False)
    assert response.status_code == 403
    assert created_user == None

def test_create_user_with_invalid_office(api_url, mysql_conn):
    delete_user_by_email(mysql_conn, user.test_new_user.email)
    
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    user_client = user.UserClient(api_url, user_auth.tokens_cookies)

    _, created_user = user_client.create_user(user.test_new_user)
    duplicate_respose, duplicated_user = user_client.create_user(user.test_new_user, False)

    assert duplicate_respose.status_code == 409
    assert duplicated_user == None

# Редактирование
def test_update_user(api_url, mysql_conn):
    delete_user_by_email(mysql_conn, user.test_new_user.email)
    delete_user_by_email(mysql_conn, "test-updated@mail.ru")
    
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    user_client = user.UserClient(api_url, user_auth.tokens_cookies)

    _, created_user = user_client.create_user(user.test_new_user)
    
    created_user.role_id = 1
    created_user.office_id = 4
    created_user.email = "test-updated@mail.ru"
    created_user.first_name = "test-updated-firstname"
    created_user.last_name = "test-updated-lastname"

    _, updated_user = user_client.update_user(created_user)

# Изменение IsActive
def test_deactivate_user(api_url, mysql_conn):
    delete_user_by_email(mysql_conn, user.test_new_user.email)
    
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    user_client = user.UserClient(api_url, user_auth.tokens_cookies)

    _, created_user = user_client.create_user(user.test_new_user)

    created_user.active = False
    _ = user_client.update_is_active(created_user)
    
    
    login_inactive_response = auth.Auth(
        api_url, auth.test_new_user).login_user(False)
    assert login_inactive_response.status_code == 403

    created_user.active = True
    _ = user_client.update_is_active(created_user)

    auth.Auth(api_url, auth.test_new_user).login_user()