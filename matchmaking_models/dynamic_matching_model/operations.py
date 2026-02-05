import bisect
import time

import redis

import domain
import random

from dataclasses import dataclass
from typing import Callable

MM_LOCK = "mm_lock"
HOT_POOL_KEY = "hot_user_pool"
CONNECTION_CHANNEL = "user_connect_channel"
TERMINATION_CHANNEL = "termination_channel"


def rating_key(id_):
    return f"rating:{id_}"


def node_key(id_):
    return f"node:{id_}"


def bucket_key(id_):
    return f"bucket:{id_}"


def bucket_users_key(bucket, node):
    return f"bucket_users:{node}:{bucket}"


def connect(conn: redis.Redis, user: domain.User):
    registered = conn.get(rating_key(user.id))
    if registered:
        return False

    pip = conn.pipeline()
    pip.multi()
    pip.set(rating_key(user.id), user.rating)
    pip.set(node_key(user.id), user.best_node)
    pip.rpush(HOT_POOL_KEY, user.id)
    pip.publish(CONNECTION_CHANNEL, user.id)
    pip.execute()
    return True


def store(conn: redis.Redis, user: domain.User, bucket: int):
    conn.set(bucket_key(user.id), bucket)
    conn.zadd(bucket_users_key(bucket, user.best_node), {str(user.id): user.rating})


def disconnect(conn: redis.Redis, user_id: int):
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
            node = pip.get(node_key(user_id))
            bucket = pip.get(bucket_key(user_id))
            pip.multi()
            pip.delete(rating_key(user_id))
            pip.delete(node_key(user_id))
            pip.delete(bucket_key(user_id))
            pip.zrem(bucket_users_key(bucket, node), user_id)
            pip.execute()
            return True
        except redis.WatchError:
            time.sleep(0.001)


def raw_disconnect(conn: redis.Redis, user_id: int):
    registered = conn.get(rating_key(user_id))
    if not registered:
        return False
    node = conn.get(node_key(user_id))
    bucket = conn.get(bucket_key(user_id))
    conn.delete(rating_key(user_id))
    conn.delete(node_key(user_id))
    conn.delete(bucket_key(user_id))
    conn.zrem(bucket_users_key(bucket, node), user_id)


def perform_match(conn: redis.Redis, percentiles: list[float],
                  connect_rate: float, disconnect_rate, node_distribution: dict[str, float]):
    conn.set(MM_LOCK, "1")
    user_ids = conn.lrange(HOT_POOL_KEY, 0, -1)
    bucket_count = len(percentiles) + 1
    matches = []
    for user_id in user_ids:
        user = domain.User(
            id=int(user_id),
            rating=int(conn.get(rating_key(user_id))),
            nodes={},
            best_node=conn.get(node_key(user_id))
        )

        bucket = bisect.bisect_left(percentiles, user.rating)

        found = False
        for other_bucket in range(bucket_count):
            if found:
                break
            buk = bucket_users_key(other_bucket, user.best_node)
            size = conn.zcard(buk)
            if not size:
                continue

            local_connect_rate = connect_rate / bucket_count * node_distribution[user.best_node]
            local_disconnect_rate = disconnect_rate / bucket_count * node_distribution[user.best_node]
            prob = 1 / 2 * max(1, local_disconnect_rate / local_connect_rate)
            if random.random() <= prob:
                conn.zadd(buk, {str(user.id): user.rating})
                pos = conn.zrank(buk, user.id)
                conn.zrem(buk, user.id)
                match: domain.Match
                if pos == 0:
                    (opponent_id, opponent_rating), = conn.zrange(buk, 0, 0, withscores=True)
                    match = domain.Match(
                        user1_id=user.id,
                        user2_id=int(opponent_id),
                        rating_diff=abs(user.rating - opponent_rating),
                        node=user.best_node
                    )
                else:
                    neighbors = conn.zrange(buk, pos - 1, pos, withscores=True)
                    if len(neighbors) == 1:
                        (opponent_id, opponent_rating), = neighbors
                        match = domain.Match(
                            user1_id=user.id,
                            user2_id=int(opponent_id),
                            rating_diff=abs(user.rating - opponent_rating),
                            node=user.best_node
                        )
                    else:
                        (left_id, left_rating), (right_id, right_rating) = neighbors
                        if abs(left_rating - user.rating) <= abs(right_rating - user.rating):
                            match = domain.Match(
                                user1_id=user.id,
                                user2_id=int(left_id),
                                rating_diff=abs(user.rating - left_rating),
                                node=user.best_node
                            )
                        else:
                            match = domain.Match(
                                user1_id=user.id,
                                user2_id=int(right_id),
                                rating_diff=abs(user.rating - right_rating),
                                node=user.best_node
                            )
                conn.zrem(buk, match.user1_id, match.user2_id)
                matches.append(match)
                found = True

        if not found:
            store(conn, user, bucket)

    pip = conn.pipeline()
    pip.multi()
    for match in matches:
        raw_disconnect(pip, match.user1_id)
        raw_disconnect(pip, match.user2_id)
    for _ in user_ids:
        pip.lpop(HOT_POOL_KEY)
    pip.set(MM_LOCK, "0")
    pip.execute()
    return matches


def terminate(conn: redis.Redis):
    conn.publish(TERMINATION_CHANNEL, 0)


@dataclass
class Connection:
    raw_event: dict


@dataclass
class Termination:
    raw_event: dict


def event_mapper(event):
    if event['type'] != 'message':
        return None
    if event['channel'] == CONNECTION_CHANNEL:
        return Connection(event)
    if event['channel'] == TERMINATION_CHANNEL:
        return Termination(event)

    assert False


class Subscription:
    def __init__(self, conn: redis.Redis):
        self.pubsub = conn.pubsub()
        self.subscribe()

    def subscribe(self):
        self.pubsub.subscribe(CONNECTION_CHANNEL)
        self.pubsub.subscribe(TERMINATION_CHANNEL)

    def unsubscribe(self):
        self.pubsub.unsubscribe()

    def events(self):
        return map(event_mapper, self.pubsub.listen())
