import json
from dataclasses import dataclass, asdict
from typing import Mapping

from itertools import product
import numpy as np


@dataclass
class UserStat:
    connected_at: float | None = None
    left_at: float | None = None
    matched: bool | None = None

    def initialized(self):
        return self.connected_at is not None and self.left_at is not None and self.matched is not None


@dataclass
class MatchStat:
    rating_diff: int
    node: str


@dataclass
class Metrics:
    disconnected_percent: float
    connection_time: Mapping[float, float]
    rating_diff: Mapping[float, int]
    node_distribution: Mapping[str, float]


def process_logs(sorted_logs):
    user_stats = {}
    match_stats = []

    for log in sorted_logs:
        match log:
            case {'type': 'connect', 'user': {'id': id_}, 'time': time}:
                assert id_ not in user_stats, id_
                user_stats[id_] = UserStat(connected_at=time)
            case {'type': 'disconnect', 'id': id_, 'time': time}:
                assert id_ in user_stats, id_
                user_stat = user_stats[id_]
                assert user_stat.left_at is None, id_
                assert user_stat.matched is None, id_
                user_stat.left_at = time
                user_stat.matched = False
            case {'type': 'match', 'time': time,
                  'match': {'user1_id': user1_id, 'user2_id': user2_id, 'rating_diff': rating_diff, 'node': node}}:
                assert user1_id in user_stats
                assert user2_id in user_stats, user2_id
                user1_stat = user_stats[user1_id]
                user2_stat = user_stats[user2_id]
                assert user1_stat.left_at is None
                assert user1_stat.matched is None
                assert user2_stat.left_at is None
                assert user2_stat.matched is None
                user1_stat.left_at = time
                user1_stat.matched = True
                user2_stat.left_at = time
                user2_stat.matched = True

                match_stats.append(MatchStat(rating_diff, node))
            case invalid_log:
                assert False, json.dumps(invalid_log, indent=2)

    user_stats = {
        k: v
        for k, v in user_stats.items() if v.initialized()
    }

    return user_stats, match_stats


def calculate_metrics(user_stats, match_stats, time_upscale, std_rating):
    disconnected_percent = sum(1 for user in user_stats.values() if not user.matched) / len(user_stats) * 100

    percentiles = [50, 75, 95, 99, 100]

    connection_times = [user.left_at - user.connected_at for user in user_stats.values() if user.matched]
    connection_time_p = {
        p: float(np.percentile(connection_times, p)) * time_upscale
        for p in percentiles
    }

    rating_diffs = [match.rating_diff for match in match_stats]
    rating_diff_p = {
        p: int(np.percentile(rating_diffs, p)) / std_rating
        for p in percentiles
    }

    nodes = set(match.node for match in match_stats)

    node_distribution = {
        node: sum(1 for match in match_stats if match.node == node)
        for node in nodes
    }

    return Metrics(
        disconnected_percent=disconnected_percent,
        connection_time=connection_time_p,
        rating_diff=rating_diff_p,
        node_distribution=node_distribution
    )


def extract_logs(it):
    res = []
    for line in it:
        res.append(json.loads(line))

    res.sort(key=lambda d: d['time'])
    return res


def run():
    filename = "../exp1_{}.jsonl"
    print('| Alg | arg | cr | dp | t50 | t95 | t99 | r50 | r95 | r99 |')

    for i, (conn_rate, (name, arg)) in enumerate(
        product(
            (1, 10, 20),
            [
                ('Naive', 5),
                ('Naive', 10),
                ('Naive', 30),
                ('Naive', 60),
                ('LP', 1),
                ('LP', 2),
                ('LP', 5),
                ('LP', 10)
            ],
        )
    ):
        time_upscale = 1200
        std_rating = 150
        with open(filename.format(i)) as f:
            logs = extract_logs(f)

        metrics = calculate_metrics(*process_logs(logs), time_upscale, std_rating)
        print("| {} | {} | {} | {:.2f} | {:.2f} | {:.2f} | {:.2f} | {:.2f} | {:.2f} | {:.2f} |".format(
            name,
            arg,
            conn_rate,
            metrics.disconnected_percent,
            metrics.connection_time[50],
            metrics.connection_time[95],
            metrics.connection_time[99],
            metrics.rating_diff[50],
            metrics.rating_diff[95],
            metrics.rating_diff[99]
        ))
        # print(name, arg, conn_rate)
        # print(json.dumps(asdict(metrics), indent=4))


if __name__ == '__main__':
    run()
