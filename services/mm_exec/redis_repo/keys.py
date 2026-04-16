from dataclasses import dataclass


@dataclass
class PlayerKeySet:
    id: str | int

    @property
    def rating(self):
        return f"rating:{self.id}"

    @property
    def fleet(self):
        return f"fleet:{self.id}"

    @property
    def name(self):
        return f"name:{self.id}"

    @property
    def displayed_rating(self):
        return f"displayed_rating:{self.id}"

    @property
    def reg_id(self):
        return f"reg_id:{self.id}"

    def keys(self):
        return [self.rating, self.fleet, self.name, self.displayed_rating, self.reg_id]
