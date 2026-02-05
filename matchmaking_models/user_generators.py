import domain
import random


def uniform_node_generator(mean_rating, std_rating, nodes, node_latency_range):
    def gen(n):
        for i in range(n):
            node_latency = {
                node: random.randint(*node_latency_range)
                for node in nodes
            }
            best_node = min(nodes, key=node_latency.__getitem__)
            yield domain.User(
                id=i,
                rating=round(random.gauss(mean_rating, std_rating)),
                nodes=node_latency,
                best_node=best_node
            )

    return gen
