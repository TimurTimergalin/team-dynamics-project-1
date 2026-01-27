import experiment
import rate_calculation
import user_generators
import json


def calculate_probabilities(connect_rate: float, disconnect_rate: float,
                            bucket_count: int, node_distribution: dict[str, float]):
    connect_res = {}
    disconnect_res = {}

    for node, prob in node_distribution.items():
        for bucket in range(bucket_count):
            connect_res[str((bucket, node))] = connect_rate / bucket_count * prob
            disconnect_res[str((bucket, node))] = disconnect_rate

    return connect_res, disconnect_res


def run():
    client = rate_calculation.Client(
        mean_rating=1500,
        std_rating=150,
        bucket_count=10
    )

    nodes = [f"node_{i + 1}" for i in range(5)]
    exp = experiment.Experiment(
        n=1000,
        client=client,
        matcher=rate_calculation.Matcher(),
        user_gen=user_generators.uniform_node_generator(
            mean_rating=1500,
            std_rating=150,
            nodes=nodes,
            node_latency_range=[30, 300]
        ),
        connect_rate=600,
        disconnect_rate=5
    )

    experiment.perform_experiment(exp)
    print("Actual")
    connect_res, disconnect_res = client.calculate_rates()
    print(json.dumps(connect_res, indent=4))
    print(json.dumps(disconnect_res, indent=4))
    print("Theoretical")
    connect_res, disconnect_res = calculate_probabilities(
        connect_rate=600,
        disconnect_rate=20,
        bucket_count=10,
        node_distribution={
            node: 0.2
            for node in nodes
        }
    )
    print(json.dumps(connect_res, indent=4))
    print(json.dumps(disconnect_res, indent=4))


if __name__ == '__main__':
    run()
