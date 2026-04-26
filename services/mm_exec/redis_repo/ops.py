import time
import uuid
from contextlib import contextmanager
from itertools import pairwise, groupby
from typing import Iterable

import config as config_py
import models
import redis

from . import keys
from . import lock

MM_LOCK = "mmlock"
MM_POOL = "mmpool"
LAST_UPDATE = "mmtime"


class NoAcquireException(Exception):
    pass


def get_time():
    return time.time_ns() // 1_000_000


class Ops:
    def __init__(self, pool: redis.Redis, config: config_py.MMExecConfig):
        self.pool: redis.Redis = pool
        self.config = config
        self.lock_owner_id = None

    def acquire(self) -> bool:
        owner_id = str(uuid.uuid4())
        if self.lock_owner_id:
            return False
        if not lock.acquire(
            self.pool,
            MM_LOCK,
            self.config.lock_timeout_millis,
            self.config.lock_acquire_timeout_millis,
            owner_id
        ):
            return False
        self.lock_owner_id = owner_id
        return True

    def verify_time(self) -> bool:
        last_update = self.pool.get(LAST_UPDATE)
        if last_update is None:
            return True
        if get_time() - int(last_update) < self.config.scheduling_period_millis:
            return False
        return True

    def update_time(self):
        self.pool.set(LAST_UPDATE, get_time())

    def read_player(self, player_id: int | str) -> models.Player | None:
        key_set = keys.PlayerKeySet(player_id)
        rating = self.pool.get(key_set.rating)
        if rating is None:
            return None
        return models.Player(
            int(player_id),
            float(rating),
            self.pool.get(key_set.fleet),
            self.pool.get(key_set.name),
            int(self.pool.get(key_set.displayed_rating)),
            self.pool.get(key_set.reg_id)
        )

    def get_players(self) -> list[models.Player]:
        player_ids = set(self.pool.lrange(MM_POOL, 0, -1))
        players = (self.read_player(id_) for id_ in player_ids)
        players = list(p for p in players if p)
        return players

    @staticmethod
    def gather_matches(players: list[models.Player]) -> list[models.Match]:
        players.sort(key=lambda p: (p.fleet, p.rating))
        res = []
        for fleet, pls in groupby(players, key=lambda p: p.fleet):
            for i, (p1, p2) in enumerate(pairwise(pls)):
                if i % 2:
                    continue
                res.append(models.Match(fleet, p1, p2))
        return res

    def remove_players(self, player_ids: Iterable[int | str]):
        for id_ in player_ids:
            key_set = keys.PlayerKeySet(id_)
            self.pool.delete(*key_set.keys())

    def return_players(self, players_taken: int, to_return: Iterable[int | str]):
        to_return = list(to_return)
        self.pool.lpop(MM_POOL, players_taken)
        if len(to_return) > 0:
            self.pool.lpush(MM_POOL, *to_return)

    def release(self):
        if not lock.release(self.pool, MM_LOCK, self.lock_owner_id):
            raise RuntimeError("Cannot unlock the lock")

    @contextmanager
    def lock(self):
        if not self.acquire():
            raise NoAcquireException()
        yield
        self.release()

    def reset(self):
        self.lock_owner_id = None
