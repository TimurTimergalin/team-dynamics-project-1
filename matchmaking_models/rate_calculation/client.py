import bisect
import datetime
from collections import defaultdict
from dataclasses import dataclass

from scipy import stats

import domain


@dataclass
class Event:
    bucket: int
    node: str
    time: float


class RateCalculationClient:
    def __init__(self, mean_rating, std_rating, bucket_count):
        self.connects = []
        self.disconnects = []
        self.start = None
        self.end = None
        self.percentiles = [
            stats.norm.ppf(1 / bucket_count * i, mean_rating, std_rating)
            for i in range(1, bucket_count)
        ]
        self.users = {}
        self.user_connects = {}
        self.user_disconnects = {}

    def get_bucket(self, rating):
        return bisect.bisect_left(self.percentiles, rating)

    def connect(self, _, user: domain.User):
        now = datetime.datetime.now().timestamp()
        if self.start is None:
            self.start = now
        self.connects.append(Event(
            bucket=self.get_bucket(user.rating),
            node=user.best_node,
            time=now
        ))

        self.users[user.id] = user
        self.user_connects[user.id] = now
        return True

    def disconnect(self, _, user_id: int):
        assert user_id in self.users
        now = datetime.datetime.now().timestamp()
        user = self.users[user_id]
        self.disconnects.append(Event(
            bucket=self.get_bucket(user.rating),
            node=user.best_node,
            time=now
        ))

        self.user_disconnects[user_id] = now
        return True

    def terminate(self, _):
        self.end = datetime.datetime.now().timestamp()

    def calculate_rates(self):
        total_time_span = self.end - self.start
        connect_res = defaultdict(int)
        disconnect_res = defaultdict(int)
        for event in self.connects:
            connect_res[str((event.bucket, event.node))] += 1
        for event in self.disconnects:
            disconnect_res[str((event.bucket, event.node))] += 1

        for k in connect_res:
            connect_res[k] /= total_time_span
        for k in disconnect_res:
            disconnect_res[k] /= total_time_span

        return dict(connect_res), dict(disconnect_res)
