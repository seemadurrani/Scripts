from github import Github
import sys
g = Github(sys.argv[1])
repo = g.get_user().get_repo(sys.argv[2])
contents = repo.get_contents("")
for post in contents:
    print(post.name)
pulls = repo.get_pulls(state='open')
for pr in pulls:
    print(pr.number)
