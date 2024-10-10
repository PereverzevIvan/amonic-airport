import datetime
import pytest
import requests
import tests.utils.datetime_utils as datetime_utils

class Schedule:
    id = None
    aircraft_id = None
    aircraft = None
    route_id = None
    route = None

    economy_price = None
    confirmed = None
    flight_number = None
    outbound = None

    def __init__(self, json: dict):
        self.id = json.get("id")

        self.aircraft_id = json.get("aircraft_id")
        self.aircraft = json.get("aircraft")
        self.route_id = json.get("route_id")
        self.route = json.get("route")
        
        self.economy_price = json.get("economy_price")
        self.confirmed = bool(json.get("confirmed"))
        self.flight_number = json.get("flight_number")
        
        if json.get("outbound") is not None:
            self.outbound = datetime_utils.parse(json["outbound"])

    @staticmethod
    def parse_schedules(json: list):
        return [Schedule(session) for session in json]

# test_new_schedule_params = {
#     "office_id": 1,
#     "email": "test@mail.ru",
#     "password": "example",
#     "first_name": "test firstname",
#     "last_name": "test lastname",
#     "birthday": "2000-12-01T00:00:00Z"
# }

# test_new_schedule = Schedule(test_new_schedule_params)

class SchedulesParams:
    def __init__(self, 
                 outbound: str = None, 
                 flight_number: str = None, 
                 departure_airport_id: int = None, 
                 arrival_airport_id: int = None, 
                 sort_by: str = None):
        self.outbound = outbound
        self.flight_number = flight_number
        self.departure_airport_id = departure_airport_id
        self.arrival_airport_id = arrival_airport_id
        self.sort_by = sort_by
    
    def to_dict(self):
        params = {}

        if self.outbound is not None:
            params["outbound"] = self.outbound
        if self.flight_number is not None:
            params["flight_number"] = self.flight_number
        if self.departure_airport_id is not None:
            params["from"] = self.departure_airport_id
        if self.arrival_airport_id is not None:
            params["to"] = self.arrival_airport_id
        if self.sort_by is not None:
            params["sort_by"] = self.sort_by
        
        return params

class ScheduleUpdateParams:
    def __init__(self, 
                 date: datetime.date = None, 
                 time: datetime.time = None, 
                 economy_price: int = None
                 ):
        self.date = date
        self.time = time
        self.economy_price = economy_price
    
    def to_dict(self):
        params = {}

        if self.date is not None:
            params["date"] = self.date.strftime("%Y-%m-%d")
        if self.time is not None:
            params["time"] = self.time.strftime("%H:%M")
        if self.economy_price is not None:
            params["economy_price"] = self.economy_price
        
        return params



class ScheduleClient:
    cookies = None

    def __init__(self, api_url, jwt_cookies):
        self.api_url = api_url
        self.cookies = jwt_cookies

    def get_schedules(self, 
            params: SchedulesParams = SchedulesParams(), 
            do_assertion=True) -> tuple[requests.Response, list[Schedule]]:
        
        response = requests.get(
            f"{self.api_url}/schedules", 
            params=params.to_dict(), 
            cookies=self.cookies)

        if do_assertion:
            assert response.status_code == 200
        elif response.status_code != 200:
            return response

        schedules = Schedule.parse_schedules(response.json())
        if do_assertion:
            assert len(schedules) > 0

            if params.outbound is not None:
                assert schedules[0].outbound.date().strftime("%Y-%m-%d") == params.outbound

            if params.flight_number is not None:
                assert schedules[0].flight_number == params.flight_number

            if params.departure_airport_id is not None:
                assert schedules[0].route["departure_airport_id"] == params.departure_airport_id

            if params.arrival_airport_id is not None:
                assert schedules[0].route["arrival_airport_id"] == params.arrival_airport_id

            if params.sort_by is not None:
                if params.sort_by == "date_time":
                    sorted_schedules = sorted(schedules, key=lambda x: x.outbound, reverse=True)
                    assert schedules == sorted_schedules

                elif params.sort_by == "confirmed":
                    sorted_schedules = sorted(schedules, key=lambda x: x.confirmed, reverse=True)
                    assert schedules == sorted_schedules

                elif params.sort_by == "ticket_price":
                    sorted_schedules = sorted(schedules, key=lambda x: x.economy_price)
                    assert schedules == sorted_schedules
                
        return response, schedules

    def get_schedule(self, schedule_id: int, do_assertion=True):
        response = requests.get(f"{self.api_url}/schedule/{schedule_id}", cookies=self.cookies)

        if do_assertion:
            assert response.status_code == 200

        return response, Schedule(response.json())
    
    def update_schedule(
            self, 
            schedule_id: int, 
            params: ScheduleUpdateParams = ScheduleUpdateParams(), 
            do_assertion=True):
        
        response = requests.patch(
            f"{self.api_url}/schedule/{schedule_id}", 
            json=params.to_dict(), 
            cookies=self.cookies)

        if do_assertion:
            assert response.status_code == 200

            _, updated_schedule = self.get_schedule(schedule_id)

            if params.economy_price is not None:
                assert updated_schedule.economy_price == params.economy_price
            if params.date is not None:
                assert updated_schedule.outbound.strftime("%Y-%m-%d") == params.date.strftime("%Y-%m-%d")
            if params.time is not None:
                assert updated_schedule.outbound.strftime("%H:%M") == params.time.strftime("%H:%M")


        return response
    
    def update_schedule_confirmed(self, schedule_id: int, set_confirmed: bool, do_assertion=True):
        response = requests.put(
            f"{self.api_url}/schedule/{schedule_id}", 
            json={
                "confirmed": set_confirmed
            }, 
            cookies=self.cookies)
        
        if do_assertion:
            assert response.status_code == 200

        return response

    def upload_csv(self, file_path: str, do_assertion=True):
        response = requests.post(
            f"{self.api_url}/schedules/upload",
            files={"file": open(file_path, "rb")},
            cookies=self.cookies
        )

        if do_assertion:
            assert response.status_code == 200

        return response, response.json()