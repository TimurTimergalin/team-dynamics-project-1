import os
import subprocess
import sys
import monorepo_root

__REPO_ROOT__ = monorepo_root.get_root()


def generate_proto(input_dir: str, output_dir_go: str, output_dir_python: str) -> None:
    for output_dir in output_dir_go, output_dir_python:
        os.makedirs(output_dir, exist_ok=True)
    entries = os.listdir(input_dir)

    for entry in entries:
        sub_path = os.path.join(input_dir, entry)
        if not os.path.isdir(sub_path):
            continue

        for output_dir in output_dir_go, output_dir_python:
            target_subdir = os.path.join(output_dir, entry)
        os.makedirs(target_subdir, exist_ok=True)

        go_cmd = [
            "protoc",
            f"--proto_path={input_dir}",
            f"--go_out={output_dir_go}",
            f"--go_grpc_out={output_dir_go}",
            f"--go_opt=paths=source_relative",
            f"--go_grpc_opt=paths=source_relative",
            f"{sub_path}/*.proto"
        ]

        python_cmd = [
            "python",
            "-m",
            "grpc_tools.protoc",
            f"--proto_path={input_dir}",
            f"--python_out={output_dir_python}",
            f"--grpc_python_out={output_dir_python}",
            f"{sub_path}\\*.proto",
        ]

        # Run the command
        print(f"Running protoc for {sub_path}...")
        subprocess.run(go_cmd, check=True)
        subprocess.run(python_cmd, check=True)


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
    input_dir = os.path.join(__REPO_ROOT__, "api", "proto")
    output_dir_go = os.path.join(__REPO_ROOT__, "gen", "proto", "go")
    output_dir_python = os.path.join(__REPO_ROOT__, "gen", "proto", "python")
    try:
        generate_proto(input_dir, output_dir_go, output_dir_python)
    except subprocess.CalledProcessError:
        print("Errors found when generating proto", file=sys.stderr)

    try:
        init_mod(output_dir_go)
    except subprocess.CalledProcessError:
        print("Errors found when initializing go mod", file=sys.stderr)


if __name__ == '__main__':
    main()
