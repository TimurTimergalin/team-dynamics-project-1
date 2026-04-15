import time

import redis

CLIENT_OWNER_ID = "mmevent"


def acquire(
        conn: redis.Redis,
        k_lock: str,
        lock_timeout_millis: int,
        acquire_timeout_millis: int,
        owner_id: str
) -> bool:
    print("calling acquire with owner_id", owner_id)
    end = time.time() + acquire_timeout_millis / 1000

    while time.time() < end:
        prev_owner_id = conn.set(k_lock, owner_id, nx=True, get=True, px=lock_timeout_millis)
        print("prev owner id is", prev_owner_id)
        if prev_owner_id is None or prev_owner_id == owner_id:
            return True
        if not conn.ttl(k_lock):
            conn.expire(k_lock, lock_timeout_millis // 1000)
        if prev_owner_id != CLIENT_OWNER_ID:
            return False
        time.sleep(0.01)

    return False


def release(conn: redis.Redis, k_lock: str, owner_id: str) -> bool:
    print("calling release with owner id", owner_id)
    pip = conn.pipeline()
    while True:
        try:
            pip.watch(k_lock)
            res = conn.get(k_lock)
            print("current lock holder is", res)
            if res == owner_id:
                pip.multi()
                pip.delete(k_lock)
                pip.execute()
                return True

            pip.unwatch()
            return False
        except redis.WatchError:
            pass
