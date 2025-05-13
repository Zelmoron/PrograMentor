from aiogram import Router, F
from aiogram.fsm.context import FSMContext
from aiogram.filters.state import State, StatesGroup
from aiogram.types import Message
from keyboards.keyboards import get_main_keyboard
from db.db_connect import setup_database_connection
from utils.logging import logger
from utils.validation import validate_nickname
from datetime import datetime
import hashlib

router = Router()

# Устанавливаем соединение с базой данных
connection, cursor = setup_database_connection()

class RegistrationForm(StatesGroup):
    """Состояния для процесса регистрации"""
    nickname = State()
    password = State()
    change_password = State()

@router.message(F.text.startswith('/start'))
async def cmd_start(message: Message):
    """Обработка команды /start"""
    await message.answer(
        '👋 Welcome! This is the PrograMentor project bot. Use the buttons below to proceed.',
        reply_markup=get_main_keyboard()
    )

@router.message(F.text == '📝 Register')
async def register_callback(message: Message, state: FSMContext):
    """Обработка регистрации"""
    try:
        tg_user_id = message.from_user.id

        # Проверка, есть ли пользователь в базе
        cursor.execute("SELECT COUNT(*) FROM users WHERE id = %s", (tg_user_id,))
        count = cursor.fetchone()[0]

        if count != 0:
            await message.answer('⚠️ You are already registered.')
        else:
            await message.answer(
                '📝 Please enter your nickname (8-15 characters).\n'
                'Use only Latin letters (a-z, A-Z), underscores (_), and digits (0-9).\n'
                'Example: <code>user_123</code>.',
                parse_mode="HTML"
            )
            await state.set_state(RegistrationForm.nickname)

    except Exception as e:
        logger.error(f'Произошла ошибка при регистрации: {e}')

@router.message(RegistrationForm.nickname)
async def handle_nickname(message: Message, state: FSMContext):
    """Обработка ввода никнейма"""
    nickname = message.text

    if validate_nickname(nickname):
        try:
            # Проверка, существует ли уже такой никнейм
            cursor.execute("SELECT COUNT(*) FROM users WHERE username = %s", (nickname,))
            count = cursor.fetchone()[0]

            if count != 0:
                await message.reply('⚠️ A user with this nickname is already registered.')
                await state.set_state(RegistrationForm.nickname)
            else:
                await message.reply('✅ Nickname accepted.')
                await state.update_data(nickname=nickname)
                await message.answer('🔑 Please enter your password:')
                await state.set_state(RegistrationForm.password)
        except Exception as e:
            logger.error(f'Произошла ошибка: {e}')
    else:
        await message.reply('❌ Invalid nickname. Please try again.')
        await state.set_state(RegistrationForm.nickname)

@router.message(RegistrationForm.password)
async def handle_password(message: Message, state: FSMContext):
    """Обработка пароля"""
    try:
        data = await state.get_data()
        nickname = data['nickname']
        password = message.text

        # Хешируем пароль
        hashed_password = hashlib.sha256(password.encode()).hexdigest()

        # Сохраняем пользователя в базе
        cursor.execute(
            "INSERT INTO users (id, username, password, created_at) VALUES (%s, %s, %s, %s)",
            (message.from_user.id, nickname, hashed_password, datetime.now())
        )
        connection.commit()

        await message.reply(
            '🎉 Registration complete! Your password has been saved.',
            reply_markup=get_main_keyboard()
        )
        await state.clear()
    except Exception as e:
        logger.error(f'Ошибка при сохранении данных: {e}')

@router.message(F.text == '🔒 Change Password')
async def change_password_callback(message: Message, state: FSMContext):
    """Обработка запроса на смену пароля"""
    try:
        # Проверка, зарегистрирован ли пользователь
        cursor.execute("SELECT COUNT(*) FROM users WHERE id = %s", (message.from_user.id,))
        count = cursor.fetchone()[0]

        if count != 0:
            await message.answer('🔒 Please enter your new password:')
            await state.set_state(RegistrationForm.change_password)
        else:
            await message.answer('⚠️ You are not registered.', reply_markup=get_main_keyboard())
    except Exception as e:
        logger.error(f'Произошла ошибка: {e}')

@router.message(RegistrationForm.change_password)
async def handle_change_password(message: Message, state: FSMContext):
    """Обработка нового пароля"""
    try:
        new_password = message.text

        # Хешируем новый пароль
        hashed_password = hashlib.sha256(new_password.encode()).hexdigest()

        # Обновляем пароль в базе
        cursor.execute(
            "UPDATE users SET password = %s, updated_at = %s WHERE id = %s",
            (hashed_password, datetime.now(), message.from_user.id)
        )
        connection.commit()

        await message.answer(
            '✅ Your password has been changed.',
            reply_markup=get_main_keyboard()
        )
        await state.clear()
    except Exception as e:
        logger.error(f'Ошибка при смене пароля: {e}')