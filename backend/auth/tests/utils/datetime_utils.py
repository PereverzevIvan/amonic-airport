from datetime import datetime
import pytz

# date_time_str = "2024-09-15T21:28:35+03:00"
golang_format_str = "%Y-%m-%dT%H:%M:%S%z"
local_tz = pytz.timezone('Europe/Moscow')

def parse(timestamp_str: str, format_str: str = golang_format_str) -> datetime:
    if timestamp_str == None:
        return None
    date_time_obj = datetime.strptime(timestamp_str, format_str)
    return date_time_obj

def now(use_microseconds = False) -> datetime:
    now_with_tz = datetime.now(local_tz)
    if not use_microseconds:
        now_with_tz = now_with_tz.replace(microsecond = 0)
    else:
        now_with_tz = now_with_tz.replace(second= now_with_tz.second+1, microsecond = 0) 
    return now_with_tz
    