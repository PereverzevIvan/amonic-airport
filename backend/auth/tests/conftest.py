import pathlib

import pytest
import requests
import logging

from pytest_mysql import factories



@pytest.fixture(scope="session")
def logger():
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
    yield logging.getLogger(__name__)

# Создание фикстуры для процесса MySQL
mysql_proc = factories.mysql_proc(port=3306)

# Создание клиентской фикстуры для подключения к MySQL
mysql_client = factories.mysql(mysql_proc)

@pytest.fixture(scope="session")
def mysql_clear():
    pass

@pytest.fixture(scope="session")
def mysql_connection(mysql_client):
    """Фикстура для подключения к базе данных MySQL."""
    # Здесь можно добавить код для инициализации базы данных, если это необходимо
    yield mysql_client
    # Код для очистки или закрытия соединения, если требуется


# Заполнение БД тестовыми данными из mysql_data.sql 
@pytest.fixture(scope="session")
def mysql_fill_data(mysql_connection):
    """Заполнение БД тестовыми данными из mysql_data.sql."""
    
    mysql_connection.execute_file('mysql_data.sql')



### API

@pytest.fixture(scope="session")
def api_url():
    """Фикстура для отправки запросов на http://localhost:3000/."""
    
    base_url = "http://localhost:3000/api"
    return base_url


# from testsuite.databases.pgsql import discover

# pytest_plugins = ['pytest_userver.plugins.postgresql']


# @pytest.fixture(scope='session')
# def service_source_dir():
#     """Path to root directory service."""
#     return pathlib.Path(__file__).parent.parent


# @pytest.fixture(scope='session')
# def initial_data_path(service_source_dir):
#     """Path for find files with data"""
#     return [
#         service_source_dir / 'postgresql/data',
#     ]


# @pytest.fixture(scope='session')
# def pgsql_local(service_source_dir, pgsql_local_create):
#     """Create schemas databases for tests"""
#     databases = discover.find_schemas(
#         'url_shortener',  # service name that goes to the DB connection
#         [service_source_dir.joinpath('postgresql/schemas')],
#     )
#     return pgsql_local_create(list(databases.values()))
