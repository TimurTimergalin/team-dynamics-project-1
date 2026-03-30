import datetime
from collections import defaultdict
from dataclasses import dataclass
from itertools import product

import redis
from scipy.optimize import linprog
from scipy.stats import norm

from log import logger
from . import operations
import json


def calculate_bucket_values(mean_rating: float, std_rating: float, bucket_count: int):
    percentiles = [
        norm.ppf(1 / bucket_count * i, mean_rating, std_rating)
        for i in range(1, bucket_count)
    ]

    res = []

    for i in range(bucket_count):
        if i == 0:
            num = -norm.pdf(percentiles[0], mean_rating, std_rating)
        elif i == bucket_count - 1:
            num = norm.pdf(percentiles[-1], mean_rating, std_rating)
        else:
            num = norm.pdf(percentiles[i - 1], mean_rating, std_rating) - norm.pdf(percentiles[i], mean_rating,
                                                                                   std_rating)
        res.append(num * bucket_count * std_rating * std_rating + mean_rating)

    return res


def calculate_reward(bucket_values, b1, b2):
    return 100 / (abs(bucket_values[b1] - bucket_values[b2]) + 1) ** 0.5


def calculate_weights(mean_rating: float, std_rating: float, connect_rate: float, disconnect_rate: float,
                      bucket_count: int, node_distribution: dict[str, float]):
    res = {}
    bucket_ratings = calculate_bucket_values(mean_rating, std_rating, bucket_count)
    for node, prob in node_distribution.items():
        node_weights = defaultdict(dict)
        local_connect_rate = connect_rate / bucket_count * prob

        minimizing_vector = [
            -calculate_reward(bucket_ratings, x, y) * local_connect_rate
            for x, y in product(range(bucket_count), repeat=2)
        ]

        bounds = [
            (0, max(1.0, local_connect_rate / disconnect_rate))
            for _ in range(bucket_count * bucket_count)
        ]

        constraints_lhs = []
        constraints_rhs = []

        for x in range(bucket_count):
            lhs_row = [0 for _ in range(bucket_count * bucket_count)]
            for y in range(bucket_count):
                lhs_row[x * bucket_count + y] += local_connect_rate
                lhs_row[y * bucket_count + x] += local_connect_rate
            constraints_lhs.append(lhs_row)
            constraints_rhs.append(local_connect_rate)

        lp_solution = linprog(
            minimizing_vector,
            A_ub=constraints_lhs,
            b_ub=constraints_rhs,
            bounds=bounds,
            method="highs"
        ).x

        for x, y in product(range(bucket_count), repeat=2):
            node_weights[x][y] = lp_solution[x * bucket_count + y]

        res[node] = dict(node_weights)

    return res


@dataclass
class DynamicMatchingMatcher:
    percentiles: list[float]
    connect_rate: float
    disconnect_rate: float
    node_distribution: dict[str, float]

    def run(self, conn: redis.Redis):
        logger.debug("Matcher started")
        subscription = operations.Subscription(conn)
        for event in subscription.events():
            if event is None:
                continue
            logger.debug(json.dumps(event.raw_event))
            if isinstance(event, operations.Connection):
                logger.debug("Matching started")
                matches = operations.perform_match(
                    conn=conn,
                    percentiles=self.percentiles,
                    connect_rate=self.connect_rate,
                    disconnect_rate=self.disconnect_rate,
                    node_distribution=self.node_distribution
                )

                logger.debug("Matching finished")

                for match in matches:
                    logger.info(json.dumps(
                        {'time': datetime.datetime.now().timestamp(), 'match': match.to_dict(), 'type': 'match'}))
            else:
                logger.debug('Terminating')
                break


if __name__ == '__main__':
    import json

    nodes = {
        f"node_{i + 1}": 0.2
        for i in range(5)
    }

    ws = calculate_weights(1500, 150, 600, 20, 15, nodes)
    print(json.dumps(ws, indent=4))
