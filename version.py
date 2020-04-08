#!/usr/bin/env python3
import subprocess
import time
import sys
import argparse

if sys.version_info < (3, 0):
    sys.stderr.write(
        "Sorry, this script could not run with Python {}.{}, "
        "Python 3.0 or higher is required\n".format(*sys.version_info))
    sys.exit(1)


def cmd(cmdline):
    return subprocess.check_output(
            cmdline.split(),
            stderr=subprocess.STDOUT,
    ).decode("utf-8")


def get_commit_info():
    hash = cmd("git log --format=format:%h -n 1").strip()
    branch = cmd("git rev-parse --abbrev-ref HEAD").strip()
    timestamp = time.gmtime(int(
            cmd("git log --format=format:%ct -n 1").strip()
    ))
    isClean = cmd("git status --porcelain").strip() == ""
    try:
        tag = cmd("git describe --exact-match --tags").strip()
    except subprocess.CalledProcessError:
        tag = None

    return dict(
        tag=tag,
        hash=hash,
        branch=branch,
        date=time.strftime("%Y.%m.%d", timestamp),
        time=time.strftime("%H:%M", timestamp),
        isClean=isClean,
    )


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--clean-tag-only",
        help="Print only tag or error if no clean tag is currenly checked out",
        action="store_true"
    )
    args = parser.parse_args()
    commitInfo = get_commit_info()
    if args.clean_tag_only:
        if not commitInfo["tag"]:
            print("Error: no tag currenly checked out", file=sys.stderr)
            exit(1)
        if not commitInfo["isClean"]:
            print("Error: directory is not clean", file=sys.stderr)
            exit(2)
        print(commitInfo["tag"])
        return

    if commitInfo["tag"]:
        versionFormat = "{tag}-{hash}-{date}"
    elif commitInfo["branch"] != "HEAD":
        versionFormat = "{branch}-{hash}-{date}-{time}"
    else:
        versionFormat = "{hash}-{date}-{time}"

    version = versionFormat.format(**commitInfo)
    if not commitInfo["isClean"]:
        version += '-notclean'

    print(version)


if __name__ == "__main__":
    main()
