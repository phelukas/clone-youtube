from kombu import Connection
from kombu.exceptions import OperationalError


def create_rabbitmq_connection() -> Connection:
    return Connection("amqp://guest:guest@host.docker.internal:5672//")
