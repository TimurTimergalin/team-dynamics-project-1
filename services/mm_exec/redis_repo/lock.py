import redis
import time
import uuid


def acquire(conn: redis.Redis, k_lock: str, lock_timeout_millis: int, acquire_timeout_millis: int) -> str | None:
    owner_id = str(uuid.uuid4())
    end = time.time() + acquire_timeout_millis / 1000

    while time.time() < end:
        if conn.setnx(k_lock, owner_id):
            conn.expire(k_lock, lock_timeout_millis // 1000)
            return owner_id
        elif not conn.ttl(k_lock):
            conn.expire(k_lock, lock_timeout_millis // 1000)
        time.sleep(0.01)

    return None


def release(conn: redis.Redis, k_lock: str, owner_id: str) -> bool:
    pip = conn.pipeline()
    while True:
        try:
            pip.watch(k_lock)
            res = conn.get(k_lock)
            if res == owner_id:
                pip.multi()
                pip.delete(k_lock)
                pip.execute()
                return True

            pip.unwatch()
            return False
        except redis.WatchError:
            pass
