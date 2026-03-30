import datetime
from dataclasses import dataclass
from . import operations
import json
import redis
from log import logger
import time


@dataclass
class NaiveModelMatcher:
    duration: float
    period: float
    max_rating_diff: int = 10000000000000000

    def run(self, conn: redis.Redis):
        now = datetime.datetime.now().timestamp()
        end_time = now + self.duration

        while now < end_time:
            logger.debug("Starting matching")
            matches = operations.perform_match(conn, self.max_rating_diff)
            logger.debug("Ending matching")
            for match in matches:
                logger.info(json.dumps({'time': now, 'match': match.to_dict(), "type": "match"}))
            time.sleep(self.period)
            now = datetime.datetime.now().timestamp()
