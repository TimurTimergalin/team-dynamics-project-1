import config as config_py
import os
import redis
from redis.retry import Retry
from redis.backoff import ExponentialBackoff
import service as service_py
import scheduler as scheduler_py
import redis_repo
import gen.proto.python.match_service.match_service_pb2_grpc as ms_grpc_pb2
import grpc
import sys
from contextlib import redirect_stdout


def get_mm_exec_config() -> config_py.MMExecConfig:
    try:
        lock_timeout = int(os.environ["MM_LOCK_TIMEOUT_MILLIS"])
        lock_acquire_timeout = int(os.environ["MM_LOCK_ACQUIRE_TIMEOUT_MILLIS"])
        scheduling_period = int(os.environ["MM_SCHEDULING_PERIOD_MILLIS"])
        match_service_address = os.environ["MM_MATCH_SERVICE_ADDRESS"]
    except KeyError as e:
        raise RuntimeError(f"Missing required environment variable: {e.args[0]}") from None
    except ValueError as e:
        raise RuntimeError(f"Invalid integer value for environment variable: {e}") from None

    return config_py.MMExecConfig(
        lock_timeout_millis=lock_timeout,
        lock_acquire_timeout_millis=lock_acquire_timeout,
        scheduling_period_millis=scheduling_period,
        match_service_address=match_service_address,
    )


def get_redis_config() -> config_py.RedisConfig:
    try:
        host = os.environ["REDIS_HOST"]
        port = int(os.environ["REDIS_PORT"])
        password = os.environ["REDIS_PASSWORD"]
        db = int(os.environ["REDIS_DB"])
        max_connections = int(os.environ["REDIS_MAX_CONNECTIONS"])
        socket_connect_timeout = float(os.environ["REDIS_SOCKET_CONNECT_TIMEOUT"])
        socket_timeout = float(os.environ["REDIS_SOCKET_TIMEOUT"])
        retry_number = int(os.environ["REDIS_RETRY_NUMBER"])
    except KeyError as e:
        raise RuntimeError(f"Missing required environment variable: {e.args[0]}") from None
    except ValueError as e:
        raise RuntimeError(f"Invalid value for environment variable: {e}") from None

    return config_py.RedisConfig(
        host=host,
        port=port,
        db=db,
        password=password,
        max_connections=max_connections,
        socket_connect_timeout=socket_connect_timeout,
        socket_timeout=socket_timeout,
        retry_number=retry_number,
    )


def main():
    mme_cfg = get_mm_exec_config()
    redis_cfg = get_redis_config()
    retry_number = redis_cfg["retry_number"]
    redis_cfg = dict(redis_cfg)
    del redis_cfg["retry_number"]
    rdb = redis.Redis(
        **redis_cfg,
        decode_responses=True,
        encoding='utf-8',
        retry=Retry(backoff=ExponentialBackoff(), retries=retry_number),
    )
    channel = grpc.insecure_channel(mme_cfg.match_service_address)
    stub = ms_grpc_pb2.MatchServiceStub(channel)

    service = service_py.MMExecService(
        redis_repo.Ops(rdb, mme_cfg),
        stub
    )
    job = service.execute
    scheduler = scheduler_py.Scheduler(job, mme_cfg.scheduling_period_millis)
    with redirect_stdout(sys.stderr):
        print("STARTING")
        scheduler.run()


if __name__ == '__main__':
    main()
