from . import operations

import redis
import domain


class DynamicMatchingClient:
    def connect(self, conn: redis.Redis,  user: domain.User):
        return operations.connect(conn, user)

    def disconnect(self, conn: redis.Redis, user_id: int):
        return operations.disconnect(conn, user_id)

    def terminate(self, conn: redis.Redis):
        return operations.terminate(conn)
