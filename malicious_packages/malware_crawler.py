import os
import csv
import requests
from git import Repo
from datetime import datetime

# Define constants
GITHUB_API_URL = "https://api.github.com"
OUTPUT_DIR = "./cloned_repos"
OUTPUT_CSV = "repo_info.csv"
REPO_URL_FILE = "repo_urls.txt"  # File containing repository URLs, one per line

# Get the GitHub token from the environment variable
TOKEN = os.environ.get("GITHUB_TOKEN")

# Ensure output directory exists
os.makedirs(OUTPUT_DIR, exist_ok=True)

def get_repositories(repo_urls):
    return repo_urls

def get_repository_info(repo_url, token):
    headers = {"Authorization": f"token {token}"}
    repo_parts = repo_url.split('/')
    if len(repo_parts) >= 5:
        owner = repo_parts[-2]
        repo_name = repo_parts[-1]
        if repo_name.endswith('.git'):
            repo_name = repo_name[:-4]  # Remove the '.git' extension
        response = requests.get(f"{GITHUB_API_URL}/repos/{owner}/{repo_name}", headers=headers)
        response.raise_for_status()
        return response.json()
    else:
        print(f"Invalid repository URL: {repo_url}")
        return None

def clone_repository(repo_url, clone_dir):
    if os.path.exists(clone_dir):
        print(f"Skipping {repo_url} as {clone_dir} already exists.")
        return
    Repo.clone_from(repo_url, clone_dir)

def get_last_commit_date(repo):
    commits = list(repo.iter_commits())
    if commits:
        return datetime.fromtimestamp(commits[0].committed_date).isoformat()
    else:
        return None

def count_lines_of_code(repo_path, extension=".go"):
    loc_count = 0
    for root, _, files in os.walk(repo_path):
        for file in files:
            if file.endswith(extension):
                with open(os.path.join(root, file), "r", encoding="utf-8") as f:
                    loc_count += sum(1 for _ in f)
    return loc_count

def main():
    with open(REPO_URL_FILE, "r") as f:
        repo_urls = [line.strip() for line in f]

    repos = get_repositories(repo_urls)

    with open(OUTPUT_CSV, mode="w", newline="", encoding="utf-8") as csvfile:
        fieldnames = ["author", "repo_name", "last_commit_date", "total_go_locs"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()

        for idx, repo_url in enumerate(repos, start=1):
            repo_info = get_repository_info(repo_url, TOKEN)
            if repo_info is None:
                continue

            repo_name = repo_info["name"]
            owner = repo_info["owner"]["login"]
            clone_dir = os.path.join(OUTPUT_DIR, f"{idx:02d}_{repo_name}")

            # Clone the repository
            print(f"Cloning {repo_name} into {clone_dir}...")
            clone_repository(repo_url, clone_dir)

	    # If the repository was skipped, continue to the next iteration
            if not os.path.exists(clone_dir):
                continue

            # Open the cloned repository with gitpython
            repo_clone = Repo(clone_dir)

            # Get the last commit date
            last_commit_date = get_last_commit_date(repo_clone)

            # Count lines of code in .go files
            total_go_locs = count_lines_of_code(clone_dir)

            # Write the repository information to CSV
            writer.writerow({
                "author": owner,
                "repo_name": repo_name,
                "last_commit_date": last_commit_date,
                "total_go_locs": total_go_locs
            })

            # Clean up by removing the cloned directory (optional)
            # shutil.rmtree(clone_dir)

if __name__ == "__main__":
    main()
