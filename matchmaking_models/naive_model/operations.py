import domain
import redis
import time
from itertools import groupby, pairwise
from log import logger


POOL_KEY = "user_pool"
MM_LOCK = "mm_lock"


def rating_key(id_):
    return f"rating:{id_}"


def node_key(id_):
    return f"node:{id_}"


def connect(conn: redis.Redis, user: domain.User):
    registered = conn.get(rating_key(user.id))
    if registered:
        return False
    pip = conn.pipeline()
    pip.multi()
    pip.set(rating_key(user.id), user.rating)
    pip.set(node_key(user.id), user.best_node)
    pip.sadd(POOL_KEY, user.id)
    pip.execute()
    return True


def raw_disconnect(conn: redis.Redis, user_id):
    registered = conn.get(rating_key(user_id))
    if not registered:
        return False

    pip = conn.pipeline()
    pip.multi()
    pip.delete(rating_key(user_id))
    pip.delete(node_key(user_id))
    pip.srem(POOL_KEY, user_id)
    pip.execute()


def disconnect(conn: redis.Redis, user_id):
    pip = conn.pipeline()
    while True:
        try:
            pip.watch(MM_LOCK)
            if pip.get(MM_LOCK) == "1":
                time.sleep(0.001)
                continue
            registered = pip.get(rating_key(user_id))
            if not registered:
                return False
            pip.multi()
            pip.delete(rating_key(user_id))
            pip.delete(node_key(user_id))
            pip.srem(POOL_KEY, user_id)
            pip.execute()
            return True
        except redis.WatchError:
            time.sleep(0.001)


def perform_match(conn: redis.Redis, max_rating_diff: int):
    conn.set(MM_LOCK, "1")
    user_ids = conn.smembers(POOL_KEY)
    print(user_ids)

    users = [
        domain.User(
            id=int(id_),
            rating=int(conn.get(rating_key(id_))),
            nodes={},
            best_node=conn.get(node_key(id_))
        )
        for id_ in user_ids
    ]

    users.sort(key=lambda u: (u.best_node, u.rating))

    res = []
    to_remove = []
    for node, users in groupby(users, lambda u: u.best_node):
        for i, (user_1, user_2) in enumerate(pairwise(users)):
            if i % 2:
                continue
            diff = abs(user_1.rating - user_2.rating)
            if diff > max_rating_diff:
                continue
            res.append(domain.Match(
                user_1.id,
                user_2.id,
                diff,
                node
            ))
            to_remove.append(user_1)
            to_remove.append(user_2)
    pip = conn.pipeline()
    pip.multi()
    for user in to_remove:
        raw_disconnect(pip, user.id)
    pip.set(MM_LOCK, "0")
    pip.execute()
    return res
