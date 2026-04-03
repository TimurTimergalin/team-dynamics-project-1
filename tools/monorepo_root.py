import os


def get_root():
    start = os.path.split(__file__)[0]
    target = ".monorepo"

    while not os.path.isfile(os.path.join(start, target)):
        start, end = os.path.split(start)
        if not end:
            raise RuntimeError("Not in monorepo - add .monorepo file to root")

    return start


if __name__ == '__main__':
    print(get_root())
