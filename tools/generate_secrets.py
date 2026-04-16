import yaml
import monorepo_root
import secrets
import uuid
import copy
import os


__REPO_ROOT__ = monorepo_root.get_root()


def generate_secret():
    rnd = secrets.token_bytes(16)
    return str(uuid.UUID(bytes=rnd))


def merge_dicts(a, b):
    result = {}
    all_keys = set(a.keys()) | set(b.keys())

    for key in all_keys:
        if key not in a:
            result[key] = b[key]
        elif key not in b:
            result[key] = a[key]
        else:
            val_a = a[key]
            val_b = b[key]

            if isinstance(val_a, dict) and isinstance(val_b, dict):
                result[key] = merge_dicts(val_a, val_b)
            else:
                if val_a == val_b:
                    result[key] = val_a
                else:
                    raise ValueError(f"Conflict for key '{key}': {val_a!r} vs {val_b!r}")

    return result


def resolve_template(tmp):
    base = tmp["base"]
    base["metadata"] = base.get("metadata", {})
    base["metadata"]["labels"] = base["metadata"].get("labels", {})
    if "metadata" not in base:
        base["metadata"] = {}

    metadata = base["metadata"]
    if "labels" not in metadata:
        metadata["labels"] = {}

    to_return = []
    for version, config in tmp["versions"].items():
        result = copy.deepcopy(base)
        result["metadata"]["labels"] = merge_dicts(result["metadata"]["labels"], config.get("labels", {}))

        result["stringData"] = result.get("stringData", {})
        for key in config.get("secrets", []):
            result["stringData"][key] = generate_secret()

        to_return.append((version, result))

    return to_return


def main():
    target_dir = os.path.join(__REPO_ROOT__, "deployment", "secrets")

    for entry in os.listdir(target_dir):
        filename = os.path.join(target_dir, entry)
        input_suffix = "_template.yaml"
        if not entry.endswith(input_suffix):
            continue
        if not os.path.isfile(filename):
            continue

        with open(filename) as f:
            yml = yaml.safe_load(f)

        res = resolve_template(yml)
        for version, data in res:
            res_filename = version + "_" + entry[:-len(input_suffix)]
            if not res_filename.endswith("_secret"):
                res_filename += "_secret"
            res_filename += ".yaml"
            with open(os.path.join(target_dir, res_filename), "w") as f:
                yaml.dump(data, f)


if __name__ == '__main__':
    main()
