import datetime
import pytest

import requests
import tests.utils.auth as auth
import tests.utils.user as user
import tests.utils.schedule as schedule
import tests.utils.datetime_utils as datetime_utils


def test_get_scedule_by_id(api_url):
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    schedule_client = schedule.ScheduleClient(api_url, user_auth.tokens_cookies)

    _, _ = schedule_client.get_schedule(1)


def test_get_scedules(api_url):
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    schedule_client = schedule.ScheduleClient(api_url, user_auth.tokens_cookies)

    _, _ = schedule_client.get_schedules()


def test_get_schedules_params(api_url):
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    schedule_client = schedule.ScheduleClient(api_url, user_auth.tokens_cookies)

    schedule_client.get_schedules(schedule.SchedulesParams(flight_number="49"))

    schedule_client.get_schedules(schedule.SchedulesParams(outbound="2017-12-05"))

    schedule_client.get_schedules(schedule.SchedulesParams(departure_airport_id=2))
    schedule_client.get_schedules(schedule.SchedulesParams(arrival_airport_id=4))
    schedule_client.get_schedules(
        schedule.SchedulesParams(departure_airport_id=2, arrival_airport_id=4)
    )

    schedule_client.get_schedules(schedule.SchedulesParams(sort_by="date_time"))

    schedule_client.get_schedules(schedule.SchedulesParams(sort_by="confirmed"))

    schedule_client.get_schedules(
        schedule.SchedulesParams(departure_airport_id=2, sort_by="ticket_price")
    )


def test_update_schedule_confirmed(api_url):
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    schedule_client = schedule.ScheduleClient(api_url, user_auth.tokens_cookies)

    _, initial_schedule = schedule_client.get_schedule(1)

    schedule_client.update_schedule_confirmed(
        schedule_id=1, set_confirmed=(not initial_schedule.confirmed)
    )

    _, udpated_schedule = schedule_client.get_schedule(1)

    assert udpated_schedule.confirmed == (not initial_schedule.confirmed)

    schedule_client.update_schedule_confirmed(
        schedule_id=1, set_confirmed=(not udpated_schedule.confirmed)
    )

    _, udpated_schedule = schedule_client.get_schedule(1)

    assert udpated_schedule.confirmed == initial_schedule.confirmed


def test_update_schedule_by_id(api_url):
    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    schedule_client = schedule.ScheduleClient(api_url, user_auth.tokens_cookies)

    schedule_client.update_schedule(
        schedule_id=1,
        params=schedule.ScheduleUpdateParams(
            date=datetime.date(2017, 10, 4),
            time=datetime.time(12, 30),
            economy_price=100,
        ),
    )

    schedule_client.update_schedule(
        schedule_id=1,
        params=schedule.ScheduleUpdateParams(
            date=datetime.date(2018, 2, 14),
            time=datetime.time(17, 10),
            economy_price=540,
        ),
    )


def delete_schedule_by_date_and_flight_number(mysql, my_date, flight_number):
    mysql.cursor().execute(
        f"""DELETE FROM `schedules`
            WHERE `Date`='{my_date}' AND FlightNumber='{flight_number}'"""
    )
    mysql.commit()


def test_schedule_upload_scv(api_url, mysql_conn):
    delete_schedule_by_date_and_flight_number(mysql_conn, "2017-09-01", "400")
    delete_schedule_by_date_and_flight_number(mysql_conn, "2017-09-01", "500")
    delete_schedule_by_date_and_flight_number(mysql_conn, "2017-09-02", "900")

    user_auth = auth.Auth(api_url, auth.test_admin)
    user_auth.login_user()

    schedule_client = schedule.ScheduleClient(api_url, user_auth.tokens_cookies)

    _, result_add = schedule_client.upload_csv("tests/data/upload_schedules_add.csv")
    assert result_add != None
    assert result_add["total_rows_cnt"] == 7
    assert result_add["successful_rows_cnt"] == 3
    assert result_add["failed_rows_cnt"] == 2
    assert result_add["missing_fields_rows_cnt"] == 1
    assert result_add["duplicated_rows_cnt"] == 1

    _, result_edit = schedule_client.upload_csv("tests/data/upload_schedules_edit.csv")
    assert result_edit != None
    assert result_edit["total_rows_cnt"] == 4
    assert result_edit["successful_rows_cnt"] == 2
    assert result_edit["failed_rows_cnt"] == 0
    assert result_edit["missing_fields_rows_cnt"] == 0
    assert result_edit["duplicated_rows_cnt"] == 2
