import pytest
import requests
import tests.utils.datetime_utils as datetime_utils

class User:
    id = None
    role_id = None
    office_id = None

    email = None
    password = None

    first_name = None
    last_name = None

    birthday = None
    active = None

    def __init__(self, json: dict):
        self.id = json.get("id")

        self.role_id = json.get("role_id")
        self.office_id = json.get("office_id")

        self.email = json.get("email")
        self.password = json.get("password")

        self.first_name = json["first_name"]
        self.last_name = json["last_name"]
        
        
        if json.get("birthday") is not None:
            self.birthday = datetime_utils.parse(json["birthday"])
        
        self.active = json.get("active")

    def to_user_params(self):
        return {
            "office_id": self.office_id,
            "email": self.email,
            "password": self.password,
            "first_name": self.first_name,
            "last_name": self.last_name,
            "birthday": datetime_utils.to_str_no_tz(self.birthday),
        }

    def to_update_user_params(self):
        return {
            "id": self.id,
            "role_id": self.role_id,
            "office_id": self.office_id,
            "email": self.email,
            "first_name": self.first_name,
            "last_name": self.last_name,
        }

    @staticmethod
    def parse_users(json: list):
        return [User(session) for session in json]

test_new_user_params = {
    "office_id": 1,
    "email": "test@mail.ru",
    "password": "example",
    "first_name": "test firstname",
    "last_name": "test lastname",
    "birthday": "2000-12-01T00:00:00Z"
}

test_new_user = User(test_new_user_params)


class UserClient:
    cookies = None

    def __init__(self, api_url, jwt_cookies):
        self.api_url = api_url
        self.cookies = jwt_cookies

    def get_user(self, user_id, do_assertion=True):
        response = requests.get(f"{self.api_url}/user/{user_id}", cookies=self.cookies)

        if do_assertion:
            assert response.status_code == 200

        return response, User(response.json())

    def create_user(self, user: User, do_assertion=True):
        response = requests.post(
            f"{self.api_url}/user",
            cookies=self.cookies,
            json=user.to_user_params()
        )

        if do_assertion:
            assert response.status_code == 201
        elif response.status_code != 201:
            return response, None
        
        created_user = User(response.json())
        
        if do_assertion:
            assert created_user.id != None and created_user.id > 0
            assert created_user.role_id == 2 # default user role
            assert created_user.office_id == user.office_id

            assert created_user.email == user.email
            assert created_user.first_name == user.first_name
            assert created_user.last_name == user.last_name
            
            assert created_user.birthday.date() == user.birthday.date()
            assert created_user.active == True

        return response, created_user

    def update_user(self, user: User, do_assertion=True):
        response = requests.patch(
            f"{self.api_url}/user",
            cookies=self.cookies,
            json=user.to_update_user_params()
        )

        if do_assertion:
            # assert response.text == "asdf"
            assert response.status_code == 200
        elif response.status_code != 200:
            return response, None
        
        updated_user = User(response.json())
        
        if do_assertion:
            assert updated_user.id == user.id
            assert updated_user.role_id == user.role_id
            assert updated_user.office_id == user.office_id

            assert updated_user.email == user.email
            assert updated_user.first_name == user.first_name
            assert updated_user.last_name == user.last_name
            
            assert updated_user.birthday.date() == user.birthday.date()

        return response, updated_user

    def update_is_active(self, user: User, do_assertion=True):
        response = requests.put(
            f"{self.api_url}/user/active",
            cookies=self.cookies,
            json={
                "id": user.id,
                "is_active": user.active
            }
        )

        if do_assertion:
            assert response.status_code == 200
        elif response.status_code != 200:
            return response
        
        if do_assertion:
            _, updated_user = self.get_user(user.id)
            assert updated_user.active == user.active

        return response