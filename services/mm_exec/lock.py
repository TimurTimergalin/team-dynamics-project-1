import redis
import time
import uuid


def acquire(conn: redis.Redis, k_lock: str, lock_timeout: int, acquire_timeout: float) -> str | None:
    owner_id = str(uuid.uuid4())
    end = time.time() + acquire_timeout

    while time.time() < end:
        if conn.setnx(k_lock, owner_id):
            conn.expire(k_lock, lock_timeout)
            return owner_id
        elif not conn.ttl(k_lock):
            conn.expire(k_lock, lock_timeout)
        time.sleep(0.01)

    return None


def release(conn: redis.Redis, k_lock: str, owner_id: str) -> bool:
    pip = conn.pipeline()

    while True:
        try:
            pip.watch(k_lock)
            if pip.get(k_lock) == owner_id:
                pip.multi()
                pip.delete(k_lock)
                pip.execute()
                return True

            pip.unwatch()
            return False
        except redis.WatchError:
            pass
