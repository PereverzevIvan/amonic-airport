from datetime import datetime
import pytz

# date_time_str = "2024-09-15T21:28:35+03:00"
golang_format_str = "%Y-%m-%dT%H:%M:%SZ%z"
golang_format_no_tz_str = "%Y-%m-%dT%H:%M:%SZ"
local_tz = pytz.timezone("Europe/Moscow")


def parse(timestamp_str: str, format_str: str = golang_format_str) -> datetime:
    if timestamp_str is None:
        return None
    # date_time_obj = datetime.strptime(timestamp_str, format_str)
    date_time_obj = datetime.fromisoformat(timestamp_str)
    return date_time_obj

def to_str(date_time_obj: datetime, format_str: str = golang_format_str) -> str:
    if date_time_obj is None:
        return None
    return date_time_obj.strftime(format_str)

def to_str_no_tz(date_time_obj: datetime, format_str: str = golang_format_no_tz_str) -> str:
    if date_time_obj is None:
        return None
    return date_time_obj.strftime(format_str)

def now(use_microseconds=False) -> datetime:
    now_with_tz = datetime.now(local_tz)
    # if not use_microseconds:
    #     now_with_tz = now_with_tz.replace(microsecond = 0)
    # else:
    #     now_with_tz = now_with_tz.replace(second= now_with_tz.second+1, microsecond = 0)
    return now_with_tz

