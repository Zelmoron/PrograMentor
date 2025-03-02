from config.config import DB_USER, DB_PASSWORD, DB_NAME, DB_HOST
import mysql.connector as con
import logging


def setup_database_connection():
    """Подключение к базе данных"""
    try:
        connection = con.connect(
            host=DB_HOST,
            user=DB_USER,
            password=DB_PASSWORD,
            database=DB_NAME
        )
        cursor = connection.cursor()
        logging.info('Соединение с базой данных установлено.')

        return connection, cursor
    except con.Error as e:
        logging.error(f'Ошибка подключения к базе данных: {e}')
        return None, None
