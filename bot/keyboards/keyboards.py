from aiogram.types import ReplyKeyboardMarkup, KeyboardButton


def get_main_keyboard():
    """Создание клавиатуры из трёх кнопок: регисрация, смена никнейма, смена пароля"""
    register_button = KeyboardButton(text='Зарегистрироваться')
    change_nickname_button = KeyboardButton(text='Изменить никнейм')
    change_password_button = KeyboardButton(text='Изменить пароль')

    keyboard = ReplyKeyboardMarkup(
        keyboard=[
            [register_button],
            [change_nickname_button],
            [change_password_button]
        ],
        resize_keyboard=True
    )
    return keyboard
