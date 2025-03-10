from config.config import DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT
import psycopg2
import logging

def setup_database_connection():
    """
    Подключение к базе данных PostgreSQL.
    """
    try:
        # Установка соединения с базой данных
        connection = psycopg2.connect(
            host=DB_HOST,
            port=DB_PORT,
            user=DB_USER,
            password=DB_PASSWORD,
            database=DB_NAME
        )
        cursor = connection.cursor()
        logging.info('Соединение с базой данных установлено.')

        return connection, cursor
    except psycopg2.Error as e:
        logging.error(f'Ошибка подключения к базе данных: {e}')
        return None, None