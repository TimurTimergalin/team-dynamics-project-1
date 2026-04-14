from dataclasses import dataclass
from typing import TypedDict


@dataclass
class MMExecConfig:
    lock_timeout_millis: int
    lock_acquire_timeout_millis: int
    scheduling_period_millis: int
    match_service_address: str


class RedisConfig(TypedDict):
    host: str
    port: int
    password: str
    db: int
    max_connections: int
    socket_connect_timeout: float
    socket_timeout: float
    retry_number: int
