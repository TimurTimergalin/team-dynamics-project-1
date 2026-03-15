MM_LOCK = "mmlock"
MM_POOL = "mmpool"


def rating(id_: str | int) -> str:
    return f"rating:{id_}"


def node(id_: int | str) -> str:
    return f"node:{id_}"
