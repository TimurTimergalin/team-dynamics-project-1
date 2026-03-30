from dataclasses import dataclass, asdict
from typing import Mapping


@dataclass
class User:
    id: int
    rating: int
    nodes: Mapping[str, int]
    best_node: str

    def to_dict(self):
        return asdict(self)


@dataclass
class Match:
    user1_id: int
    user2_id: int
    rating_diff: int
    node: str

    def to_dict(self):
        return asdict(self)
