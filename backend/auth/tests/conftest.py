import pathlib

import pytest
import requests
import logging
import mysql.connector

# from pytest_mysql import factories

# Для вывода сообщений об ошибке со значениями
pytest.register_assert_rewrite("tests.utils")


@pytest.fixture(scope="session")
def logger():
    logging.basicConfig(
        level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
    )
    yield logging.getLogger(__name__)

### API
@pytest.fixture(scope="session")
def api_url():
    """Фикстура для отправки запросов на http://localhost:3000/."""

    base_url = "http://localhost:3000/api"
    return base_url

### mysql_connector
@pytest.fixture(scope="module")
def mysql_conn():
    # Настройки подключения
    conn = mysql.connector.connect(
        host='localhost',
        port=3306,
        user='admin',
        password='admin',
        database='airplanes'
    )
    yield conn  # Возвращаем соединение

    # Закрываем соединение после завершения тестов
    conn.close()
