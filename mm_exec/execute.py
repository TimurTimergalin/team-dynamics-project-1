import redis
import keys
import lock
from dataclasses import dataclass
from itertools import groupby, pairwise
import uuid


@dataclass
class User:
    id: int
    rating: int
    node: str


@dataclass
class Match:
    node: str
    user1: User
    user2: User


def remove_user_from_pool(mm_pool_conn: redis.Redis, user_id: int):
    mm_pool_conn.delete(keys.rating(user_id))
    mm_pool_conn.delete(keys.node(user_id))
    mm_pool_conn.srem(keys.MM_POOL, user_id)


def gather_matches(mm_pool_conn: redis.Redis) -> list[Match]:
    user_ids = mm_pool_conn.lrange(keys.MM_POOL, 0, -1)
    users = [
        User(id_, int(mm_pool_conn.get(keys.rating(id_))), mm_pool_conn.get(keys.node(id_)))
        for id_ in map(int, user_ids)
    ]

    users.sort(key=lambda user: (user.node, user.rating))
    matches = []
    for node, local_users in groupby(users, key=lambda user: user.node):
        for i, (user1, user2) in enumerate(local_users):
            if i % 2:
                continue
            matches.append(Match(node, user1, user2))

    return matches


def cleanup_pool(mm_pool_conn: redis.Redis, matches: list[Match]):
    for match in matches:
        remove_user_from_pool(mm_pool_conn, match.user1.id)
        remove_user_from_pool(mm_pool_conn, match.user2.id)



def execute():
    pass
