from dataclasses import dataclass
import domain
import redis
import datetime
import random_process
import heapq
from log import logger, clear_in_memory_logs, get_in_memory_logs
from itertools import cycle
import json
from typing import Protocol
from typing import Callable, Iterable
import threading
import random


class Client(Protocol):
    def connect(self, conn: redis.Redis,  user: domain.User):
        pass

    def disconnect(self, conn: redis.Redis, user_id: int):
        pass

    def terminate(self, conn: redis.Redis):
        pass


class Matcher(Protocol):
    def run(self, conn: redis.Redis):
        pass


@dataclass
class Connect:
    user: domain.User

    def __lt__(self, other):
        return False


@dataclass
class Disconnect:
    user_id: int

    def __lt__(self, other):
        return False


@dataclass
class Experiment:
    n: int
    client: Client
    matcher: Matcher
    user_gen: Callable[[int], Iterable[domain.User]]
    connect_rate: float
    disconnect_rate: float


def run_client(conn: redis.Redis, n: int, client, user_gen, connect_rate, disconnect_rate):
    start = datetime.datetime.now().timestamp() + 1

    events = list(
        zip(
            random_process.generate_exponential_process(n, connect_rate, start),
            map(lambda x: Connect(x), user_gen(n))
        )
    )
    heapq.heapify(events)

    while events:
        event_time, event = heapq.heappop(events)
        now = datetime.datetime.now().timestamp()
        while now < event_time:
            now = datetime.datetime.now().timestamp()
            pass

        match event:
            case Connect(user):
                if client.connect(conn, user):
                    logger.info(json.dumps({'type': 'connect', 'time': now, 'user': user.to_dict()}))
                    heapq.heappush(events, (now + random_process.random_exp(disconnect_rate), Disconnect(user.id)))

            case Disconnect(user_id):
                if client.disconnect(conn, user_id):
                    logger.info(json.dumps({'type': 'disconnect', 'time': now, 'id': user_id}))

    client.terminate(conn)


def run_matcher(conn: redis.Redis, matcher):
    matcher.run(conn)


def perform_experiment(experiment: Experiment, host='localhost', port=6379):
    conn = redis.Redis(host, port)
    conn.flushall()
    clear_in_memory_logs()
    conn.close()
    client_t = threading.Thread(target=run_client, args=(
        redis.Redis(host, port, decode_responses=True),
        experiment.n,
        experiment.client,
        experiment.user_gen,
        experiment.connect_rate,
        experiment.disconnect_rate
    ))

    matcher_t = threading.Thread(target=run_matcher, args=(
        redis.Redis(host, port, decode_responses=True),
        experiment.matcher
    ))

    client_t.start()
    matcher_t.start()
    client_t.join()
    matcher_t.join()
    return get_in_memory_logs()

