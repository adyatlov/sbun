# SBun
Tool for analyzing [DC/OS service diagnostics bundle](https://support.d2iq.com/s/article/create-service-diag-bundle)

## Usage

```
$ sbun [-p <service diagnostics bundle directory>] <command>
```

Launch the following command to see the list of commands:

```
$ sbun help
```

## Features

* Writes service task list to the standard output or file in the CSV format. 
* Checks for updates and updates itself.

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

