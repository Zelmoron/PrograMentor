from aiogram import Router, F
from aiogram.fsm.context import FSMContext
from aiogram.filters.state import State, StatesGroup
from aiogram.types import Message
from keyboards.keyboards import get_main_keyboard
from db.db_connect import setup_database_connection
from utils.logging import logger
from utils.validation import validate_nickname
import hashlib


router = Router()

# Установка соединения с базой данных
connection, cursor = setup_database_connection()


class RegistrationForm(StatesGroup):
    """Форма состояний при регистрации"""
    nickname = State()
    password = State()
    change_nickname = State()
    change_password = State()


@router.message(F.text.startswith('/start'))
async def cmd_start(message: Message):
    """Обработка команды /start"""
    await message.answer(
        'Добро пожаловать! Это бот проекта PrograMentor. Для действий используйте кнопки ниже.',
        reply_markup=get_main_keyboard()
    )


@router.message(F.text == 'Зарегистрироваться')
async def register_callback(message: Message, state: FSMContext):
    """Обработка регистрации"""
    try:
        tg_user_id = message.from_user.id

        # Проверка наличия пользователя в базе данных
        cursor.execute("SELECT COUNT(*) FROM tg_users WHERE tg_id = %s", (tg_user_id,))
        count = cursor.fetchone()[0]

        if count != 0:
            await message.answer('Вы уже зарегистрированы.')
        else:
            await message.answer(
                'Введите ваш никнейм (8-15 символов).\nИспользуйте только латинские буквы (a-z, A-Z), '
                'нижнее подчёркивание (_) и цифры (0-9).\nПример: <code>user_123</code>.',
                parse_mode="HTML"
            )
            await state.set_state(RegistrationForm.nickname)

    except Exception as e:
        logger.error(f'При регистрации произошла ошибка: {e}')


@router.message(RegistrationForm.nickname)
async def handle_nickname(message: Message, state: FSMContext):
    """Обработка ввода никнейма"""
    nickname = message.text

    if validate_nickname(nickname):
        try:
            # Проверка наличия никнейма в базе данных
            cursor.execute("SELECT COUNT(*) FROM users_info WHERE nickname = %s", (nickname,))
            count = cursor.fetchone()[0]

            if count != 0:
                await message.reply('Пользователь с таким никнеймом уже зарегистрирован.')
                await state.set_state(RegistrationForm.nickname)
            else:
                await message.reply('Никнейм принят.')
                await state.update_data(nickname=message.text)
                await message.answer('Введите ваш пароль:')
                await state.set_state(RegistrationForm.password)
        except Exception as e:
            logger.error(f'Произошла ошибка: {e}')
    else:
        await message.reply('Никнейм не валиден. Пожалуйста, попробуйте снова.')
        await state.set_state(RegistrationForm.nickname)


@router.message(RegistrationForm.password)
async def handle_password(message: Message, state: FSMContext):
    """Обработка ввода пароля"""
    try:
        data = await state.get_data()
        nickname = data['nickname']
        password = message.text

        # Хеширование пароля
        hashed_password = hashlib.sha256(password.encode()).hexdigest()

        # Добавление пользователя в базу данных
        cursor.execute("INSERT INTO tg_users (tg_id) VALUES (%s)", (message.from_user.id,))
        cursor.execute("INSERT INTO users_info (tg_id, nickname, password) VALUES ((SELECT id FROM tg_users WHERE tg_id = %s), %s, %s)",
                       (message.from_user.id, nickname, hashed_password))
        connection.commit()

        await message.reply(
            f'Регистрация завершена! Ваш пароль сохранен.',
            reply_markup=get_main_keyboard()
        )
        await state.clear()
    except Exception as e:
        logger.error(f'Произошла ошибка при сохранении данных: {e}')


@router.message(F.text == 'Изменить никнейм')
async def change_nickname_callback(message: Message, state: FSMContext):
    """Обработка изменения никнейма"""
    try:
        # Проверка наличия пользователя в базе данных
        cursor.execute("SELECT COUNT(*) FROM tg_users WHERE tg_id = %s", (message.from_user.id,))
        count = cursor.fetchone()[0]

        if count != 0:
            await message.answer('Введите новый никнейм:')
            await state.set_state(RegistrationForm.change_nickname)
        else:
            await message.answer('Вы не зарегистрированы.', reply_markup=get_main_keyboard())
    except Exception as e:
        logger.error(f'Произошла ошибка: {e}')


@router.message(RegistrationForm.change_nickname)
async def handle_change_nickname(message: Message, state: FSMContext):
    """Обработка нового никнейма"""
    try:
        new_nickname = message.text

        if validate_nickname(new_nickname):
            # Проверка наличия нового никнейма в базе данных
            cursor.execute("SELECT COUNT(*) FROM users_info WHERE nickname = %s", (new_nickname,))
            count = cursor.fetchone()[0]

            if count != 0:
                await message.reply('Пользователь с таким никнеймом уже зарегистрирован.')
                await state.set_state(RegistrationForm.change_nickname)
            else:
                # Обновление никнейма в базе данных
                cursor.execute("UPDATE users_info SET nickname = %s WHERE tg_id = (SELECT id FROM tg_users WHERE tg_id = %s)",
                               (new_nickname, message.from_user.id))
                connection.commit()
                await message.answer(
                    f'Ваш никнейм изменён на: {new_nickname}',
                    reply_markup=get_main_keyboard()
                )
                await state.clear()
        else:
            await message.reply('Никнейм не валиден. Пожалуйста, попробуйте снова.')
            await state.set_state(RegistrationForm.change_nickname)
    except Exception as e:
        logger.error(f'Произошла ошибка при изменении никнейма: {e}')


@router.message(F.text == 'Изменить пароль')
async def change_password_callback(message: Message, state: FSMContext):
    """Обработка изменения пароля"""
    try:
        # Проверка наличия пользователя в базе данных
        cursor.execute("SELECT COUNT(*) FROM tg_users WHERE tg_id = %s", (message.from_user.id,))
        count = cursor.fetchone()[0]

        if count != 0:
            await message.answer('Введите новый пароль:')
            await state.set_state(RegistrationForm.change_password)
        else:
            await message.answer('Вы не зарегистрированы.', reply_markup=get_main_keyboard())
    except Exception as e:
        logger.error(f'Произошла ошибка: {e}')


@router.message(RegistrationForm.change_password)
async def handle_change_password(message: Message, state: FSMContext):
    """Обработка нового пароля"""
    try:
        new_password = message.text

        # Хеширование пароля
        hashed_password = hashlib.sha256(new_password.encode()).hexdigest()

        # Обновление пароля в базе данных
        cursor.execute("UPDATE users_info SET password = %s WHERE tg_id = (SELECT id FROM tg_users WHERE tg_id = %s)",
                       (hashed_password, message.from_user.id))
        connection.commit()
        await message.answer(
            f'Ваш пароль изменён.',
            reply_markup=get_main_keyboard()
        )
        await state.clear()
    except Exception as e:
        logger.error(f'Произошла ошибка при изменении пароля: {e}')