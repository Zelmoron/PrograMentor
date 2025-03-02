import logging

logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    filename='reg_bot_logs.log',
    level=logging.INFO,
    filemode='a'
)

logger = logging.getLogger(__name__)
