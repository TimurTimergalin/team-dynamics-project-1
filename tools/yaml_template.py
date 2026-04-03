import os
import secrets
import uuid
import monorepo_root
from copy import deepcopy
import dataclasses
import yaml
import sys

__REPO_ROOT__ = monorepo_root.get_root()

template_resolvers = {}


def template_resolver(name):
    def dec(f):
        template_resolvers[name] = f
        return f

    return dec


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


def generate_secret():
    rnd = secrets.token_bytes(16)
    return str(uuid.UUID(bytes=rnd))


@dataclasses.dataclass
class TemplateResolution:
    version: str
    data: dict


@template_resolver("merge")
def merge_resolver(base, data):
    return merge_dicts(base, data)


@template_resolver("secrets")
def secretes_resolver(base, data):
    assert isinstance(data, list)
    string_data = base.get("stringData", {})
    for sec in data:
        if sec in string_data:
            raise RuntimeError(f"conflict: {sec} already in stringData")
        string_data[sec] = generate_secret()
    base["stringData"] = string_data
    return base


def resolve_template(data):
    base = data.get("base", {})
    res = []
    for version, version_template in data.get("versions", {}).items():
        final_item = deepcopy(base)
        for resolver, resolver_data in version_template.items():
            final_item = template_resolvers[resolver](final_item, resolver_data)
        res.append(TemplateResolution(version, final_item))
    return res


def resolve_in_folder(dir_):
    for current_dir, _, filenames in os.walk(dir_):
        for filename in filenames:
            if not filename.endswith(".yaml.template"):
                continue
            with open(os.path.join(current_dir, filename)) as f:
                data = yaml.safe_load(f)
            resolutions = resolve_template(data)
            for resolution in resolutions:
                new_filename = resolution.version + "_" + filename[:-len(".template")]
                with open(os.path.join(current_dir, new_filename), "w") as f:
                    yaml.dump(resolution.data, f)


def main(argv):
    for i, arg in enumerate(argv):
        if not i:
            continue
        resolve_in_folder(os.path.join(__REPO_ROOT__, arg))


if __name__ == '__main__':
    main(sys.argv)
