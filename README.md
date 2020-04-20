# sbun
Tool for analyzing [DC/OS service diagnostics bundle](https://support.d2iq.com/s/article/create-service-diag-bundle)

## Usage

```
$ cd <service diagnostics bundle directory>
$ sbun
```

## Features

Writes service task list to the standard output in the CSV format. The order of columns is:

1. task name
1. starting timestamp
1. running timestamp
1. killed timestamp
1. task ID
1. path to the task directory

## How to release

1. Install [GoReleaser](https://goreleaser.com/install/).
2. Create [Github personal access token](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line)
    with the `repo` scope and export it as an environment variable called `GITHUB_TOKEN`:

  	```bash
  	$ export GITHUB_TOKEN=<your personal GitHub access token>
  	```

    Please find more information about this step [here](https://goreleaser.com/environment/).
3. Create a Git tag which adheres to [semantic versioning](https://semver.org/) and
    push it to GitHub:

    ```bash
    $ git tag -a v1.9.8 -m "Release v1.9.8"
    $ git push origin v1.9.8
    ```

    If you made a mistake on this step, you can delete the tag remotely and locally:

    ```bash
    $ git push origin :refs/tags/v1.9.8
    $ git tag --delete v1.9.8
    ```

4. Test that the build works with the following command:

    ```bash
    $ goreleaser release --skip-publish --rm-dist
    ```

5. If everything is fine publish the build with the following command:

    ```bash
	$ goreleaser release --rm-dist
    ```

