import os
import re
import sys
import subprocess
import json
import time
import argparse

def change_directory(folder):
    try:
        os.chdir(folder)
    except FileNotFoundError:
        print(f"Folder '{folder}' not found.")
    except PermissionError:
        print(f"Permission denied to access '{folder}'.")


def is_go_package(dir_path):
    go_files = [f for f in os.listdir(dir_path) if f.endswith('.go')]
    for go_file in go_files:
        file_path = os.path.join(dir_path, go_file)
        with open(file_path, 'r', encoding='utf-8') as file:
            content = file.read()
            match = re.search(r'\bpackage\s+(\w+)\b', content)
            if match:
                package_name = match.group(1)
                return True, package_name
    return False, None


def can_build_go_package(dir_path):
    try:
        result = subprocess.run(["go", "build"], cwd=dir_path, capture_output=True, text=True, timeout=60)
        if result.returncode == 0:
            # Compilation successful
            return True, ""
        else:
            # Compilation failed, return error output
            return False, result.stderr.strip()
    except Exception as e:
        # Exception occurred during compilation
        return False, str(e)


def explore_go_packages(base_dir):
    package_list = []
    total_subdirs = sum([len(dirs) for _, dirs, _ in os.walk(base_dir)])
    processed_subdirs = 0

    def print_progress(processed_subdirs, total_subdirs):
        progress = processed_subdirs / total_subdirs * 100
        sys.stdout.write('\r')
        sys.stdout.write("[%-50s] %.2f%%  (%d/%d)" % ('=' * int(progress / 2), progress, processed_subdirs, total_subdirs))
        sys.stdout.flush()


    # Check if the parent folder is a package
    is_go, package_name = is_go_package(base_dir)
    if is_go:
        can_build, output = can_build_go_package(base_dir)
        if can_build:
            package_list.append({"name": package_name, "path": base_dir})
            #print(f"Directory Path: {base_dir}, Package Name: {package_name}")
        #else:
            #print(f"Directory Path: {base_dir}, Compilation Error: {output}")

    # Recursively check each subfolder
    for root, dirs, files in os.walk(base_dir):
        for directory in dirs:
            dir_path = os.path.join(root, directory)
            is_go, package_name = is_go_package(dir_path)
            if is_go:
                can_build, output = can_build_go_package(dir_path)
                if can_build:
                    package_list.append({"name": package_name, "path": dir_path})
                    #print(f"\nDirectory Path: {dir_path}, Package Name: {package_name}")
                #else:
                    #print(f"\nDirectory Path: {dir_path}, Compilation Error: {output}")
            processed_subdirs += 1
            print_progress(processed_subdirs, total_subdirs)

    return package_list


if __name__ == "__main__":

    print("""
Welcome to RUNCAPSLOCK - The Go Module Analyzer
RUNCAPSLOCK is a tool for searching all the packages in a Golang module,
and executing capability analysis using Capslock.
    """)

    parser = argparse.ArgumentParser()
    parser.add_argument("-output", help="Specify the output format for the capabilities list ('v' for verbose, 'json' for JSON, 'compare' for compare feature)", required=True)
    parser.add_argument("-module", help="Specify the module path (relative or absolute)", required=True)
    parser.add_argument("-packages", help="Specify the packages output file name", required=True)
    parser.add_argument("-capabilities", help="Specify the capabilities output file name", required=True)
    parser.add_argument("-comparewith", help="Specify the capabilities file to compare with", required=False)

    args = parser.parse_args()
    output_format = args.output
    module_folder = args.module
    pkgs_output_file = args.packages
    caps_output_file = args.capabilities
    compare_file = args.comparewith

    # Check output format
    if output_format not in ('v', 'json', 'compare'):
        print("Invalid output format. Use '-output=v' for verbose, '-output=json' for JSON format, '-output=compare' for compare feature.")
        system.exit(1)

    # Handling relative paths
    starting_dir = os.getcwd()
    pkgs_output_file = os.path.join(starting_dir, pkgs_output_file)
    caps_output_file = os.path.join(starting_dir, caps_output_file)
    module_folder = os.path.normpath(module_folder)
    change_directory(module_folder)
    module_path = os.getcwd()

    # Retrieve the list of packages paths in the module
    print("""Starting to retrieve package list for the module: """, module_path)
    print("\n")
    pkgs_list = explore_go_packages(module_path)
    pkgs_paths = [package['path'] for package in pkgs_list]

    with open(pkgs_output_file, 'w') as pkgs_file:
        json.dump(pkgs_list, pkgs_file, indent=4)
    print("\nPackages list saved to", pkgs_output_file)


    # Execute Capslock on the explored packages
    if output_format == 'v':
        command = ["capslock", "-packages", ','.join(pkgs_paths), "-output=v"]
    elif output_format == 'json':
        command = ["capslock", "-packages", ','.join(pkgs_paths), "-output=json"]
    if output_format == 'compare':
        compare_file = os.path.normpath(compare_file)
        command = ["capslock", "-packages", ','.join(pkgs_paths), "-output=compare", compare_file]

    print("""Starting capability analysis. Please wait...""")

    try:
        with open(caps_output_file, "w") as cap_file:
            subprocess.run(command, check=True, stdout=cap_file)
            print("Capslock output redirected to", caps_output_file)
    except subprocess.CalledProcessError as e:
        print("Error running command:", e)
