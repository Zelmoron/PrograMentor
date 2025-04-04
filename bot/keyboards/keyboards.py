from aiogram.types import ReplyKeyboardMarkup, KeyboardButton


def get_main_keyboard():
    """Создание клавиатуры из трёх кнопок: регисрация, смена пароля"""
    register_button = KeyboardButton(text='📝 Зарегистрироваться')
    change_password_button = KeyboardButton(text='🔒 Изменить пароль')

    keyboard = ReplyKeyboardMarkup(
        keyboard=[
            [register_button],
            [change_password_button]
        ],
        resize_keyboard=True
    )
    return keyboard
