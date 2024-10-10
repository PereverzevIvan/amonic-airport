import datetime
import pytest
import requests
import tests.utils.datetime_utils as datetime_utils

class Airport:
    id = None
    country_id = None
    iata_code = None
    name = None

    def __init__(self, json: dict):
        self.id = json.get("contry_id")

        self.country_id = json.get("contry_id")
        self.iata_code = json.get("iata_code")
        self.name = json.get("name")

    @staticmethod
    def parse_airports(json: list):
        return [Airport(session) for session in json]


class AirportClient:
    cookies = None

    def __init__(self, api_url, jwt_cookies):
        self.api_url = api_url
        self.cookies = jwt_cookies

    def get_airports(self, do_assertion=True) -> tuple[requests.Response, list[Airport]]:
        
        response = requests.get(
            f"{self.api_url}/airports", 
            cookies=self.cookies)

        if do_assertion:
            assert response.status_code == 200
        elif response.status_code != 200:
            return response

        airports = Airport.parse_airports(response.json())
        if do_assertion:
            assert len(airports) > 0
            for airport in airports:
                assert len(airport.iata_code) == 3 

            
                
        return response, airports
    