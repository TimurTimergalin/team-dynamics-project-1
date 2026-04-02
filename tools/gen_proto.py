import os
import subprocess
import sys

__PROJECT_ROOT__ = os.path.split(os.path.split(__file__)[0])[0]


def generate_proto(input_dir: str, output_dir: str) -> None:
    os.makedirs(output_dir, exist_ok=True)
    entries = os.listdir(input_dir)

    for entry in entries:
        sub_path = os.path.join(input_dir, entry)
        if not os.path.isdir(sub_path):
            continue

        target_subdir = os.path.join(output_dir, entry)
        os.makedirs(target_subdir, exist_ok=True)

        cmd = [
            "protoc",
            f"--proto_path={input_dir}",
            f"--go_out={output_dir}",
            f"--go_grpc_out={output_dir}",
            f"--go_opt=paths=source_relative",
            f"--go_grpc_opt=paths=source_relative",
            f"{sub_path}/*.proto"
        ]

        # Run the command
        print(f"Running protoc for {sub_path}...")
        subprocess.run(cmd, check=True)


def init_mod(module_path):
    cmd = [
        "go",
        "mod",
        "init",
        "team_dynamics/api/proto"
    ]
    print("initializing go module")
    subprocess.run(cmd, check=False, cwd=module_path)

    cmd = [
        "go",
        "mod",
        "tidy"
    ]
    print("fixing go.mod")
    subprocess.run(cmd, check=True, cwd=module_path)

    cmd = [
        "go",
        "build",
        "./..."
    ]
    print("building")
    subprocess.run(cmd, check=True, cwd=module_path)


def main():
    input_dir = os.path.join(__PROJECT_ROOT__, "api", "proto")
    output_dir = os.path.join(__PROJECT_ROOT__, "gen", "proto", "go")
    try:
        generate_proto(input_dir, output_dir)
    except subprocess.CalledProcessError:
        print("Errors found when generating proto", file=sys.stderr)

    try:
        init_mod(output_dir)
    except subprocess.CalledProcessError:
        print("Errors found when initializing go mod", file=sys.stderr)


if __name__ == '__main__':
    main()
