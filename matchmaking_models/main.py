from scipy.stats import norm

import dynamic_matching_model
import experiment
import naive_model
import user_generators
from itertools import chain


def naive_model_experiment(mean_rating, std_rating, connect_rate, disconnect_rate, node_count, duration, period, n):
    nodes = [f"node_{i + 1}" for i in range(node_count)]

    return experiment.Experiment(
        n=n,
        client=naive_model.Client(),
        matcher=naive_model.Matcher(
            duration=duration,
            period=period
        ),
        user_gen=user_generators.uniform_node_generator(
            mean_rating=mean_rating,
            std_rating=std_rating,
            nodes=nodes,
            node_latency_range=[30, 300]
        ),
        connect_rate=connect_rate,
        disconnect_rate=disconnect_rate,
    )


def dynamic_matching_experiment(bucket_count, mean_rating, std_rating, connect_rate, disconnect_rate, node_count,
                                n):
    percentiles = [
        float(norm.ppf(1 / bucket_count * i, mean_rating, std_rating))
        for i in range(1, bucket_count)
    ]

    nodes = [f"node_{i + 1}" for i in range(node_count)]

    return experiment.Experiment(
        n=n,
        client=dynamic_matching_model.Client(),
        matcher=dynamic_matching_model.Matcher(
            percentiles=percentiles,
            connect_rate=connect_rate,
            disconnect_rate=disconnect_rate,
            node_distribution={
                node: 0.2
                for node in nodes
            }
        ),
        user_gen=user_generators.uniform_node_generator(
            mean_rating=mean_rating,
            std_rating=std_rating,
            nodes=nodes,
            node_latency_range=[30, 300]
        ),
        connect_rate=connect_rate,
        disconnect_rate=disconnect_rate,
    )


def main():
    time_scale = 1200
    disconnect_rate = 20
    n = 1000

    i = 0
    for connect_rate in 50 * time_scale / 60, 100 * time_scale / 60, 200 * time_scale / 60:
        exps = chain(
            [
                naive_model_experiment(1500, 150, connect_rate, disconnect_rate, 5, 12, period, n)
                for period in (5 / time_scale, 10 / time_scale, 30 / time_scale, 60 / time_scale)
            ],
            [dynamic_matching_experiment(10, 1500, 150, connect_rate, disconnect_rate, 5, n)]
        )

        for exp in exps:
            logs = experiment.perform_experiment(exp)

            with open(f'exp2_{i}.jsonl', 'w') as f:
                f.write(logs)
            i += 1


if __name__ == '__main__':
    main()
