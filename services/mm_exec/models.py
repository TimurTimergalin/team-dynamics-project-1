from dataclasses import dataclass


@dataclass
class Player:
    id: int
    rating: float
    fleet: str
    name: str
    displayed_rating: int
    reg_id: str


@dataclass
class Match:
    fleet: str
    player1: Player
    player2: Player
