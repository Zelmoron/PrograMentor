import asyncio
from aiogram import Bot, Dispatcher
from config.config import BOT_TOKEN
from handlers.handlers import router


bot = Bot(token=BOT_TOKEN)
dp = Dispatcher()

dp.include_router(router)


async def main():
    await bot.delete_webhook(drop_pending_updates=True)
    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())
