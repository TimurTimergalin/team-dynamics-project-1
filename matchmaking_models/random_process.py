import math
import random


def inverse_exp_cdf(lambda_, x):
    return -math.log(1 - x) / lambda_


def random_exp(lambda_):
    return inverse_exp_cdf(lambda_, random.random())


def generate_exponential_process(n, rate, start):
    cur = start
    for _ in range(n):
        yield cur
        cur += random_exp(rate)
