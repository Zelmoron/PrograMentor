import re


def validate_nickname(nickname: str) -> bool:
    """
    Проверяет, соответствует ли никнейм заданным условиям.

    :param nickname: Никнейм для проверки.
    :return: True, если никнейм валиден, False иначе.
    """
    pattern = r'^[a-zA-Z0-9_]{8,15}$'
    return bool(re.match(pattern, nickname))
