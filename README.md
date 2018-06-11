# AWS Parameter Store Manager

Easily list or export your AWS Parameter Store values.

> **IMPORTANT:** Secrets are decrypted and displayed on the console in clear text.

<details>
<summary>
<b>Table of Contents</b>
</summary>

TBD

</details>

## Building

This step only needs to happen on first use, or when you have updated the code from this repository.

Install the [Glide] package manager, install the dependencies, then build the app for the current OS.

```bash
glide install
make build
```

<asciicast>

## Running

### Configure Credentials

The [AWS CLI Tools] and the various AWS SDKs already have a well-established pattern for managing credentials, so this software piggybacks on those existing patterns.

If you've already set-up the AWS CLI Tools, you can set default credentials that way.

```bash
aws configure
```

You can also create secondary profiles.

```bash
aws configure --profile production
```

### Configure Profile

By default, `pstore` will use the _default_ credentials from `aws configure`. If you are using multiple profiles, you can use a specific profile by passing the `--profile` flag.

```bash
pstore list --profile production
```

### Manual Credentials

If you would prefer to pass credentials manually, use the environment variables.

```bash
AWS_ACCESS_KEY_ID=AKIAJHJBEXAMPLEW7VJA \
AWS_SECRET_ACCESS_KEY=QIBAlwMNrzExampleG9iRf1ttflI0PDooExample \
AWS_DEFAULT_REGION=us-west-2 \
pstore list
```

See [AWS CLI Configuration Variables](https://docs.aws.amazon.com/cli/latest/topic/config-vars.html) for more information.

## Parameters

For all commands, you can pass an _argument_ to the command line. This argument is the _path_ of the Parameter Store value that you want to see. The default value is `/` (root/all).

### `list`

```bash
pstore --profile=staging list /awesome-app
```

An example response looks like:

```plain
+---------------------------------------+--------------+
| Key                                   | Value        |
+---------------------------------------+--------------+
| /awesome-app/staging/environment-name | staging      |
| /awesome-app/staging/project-name     | awesome-app  |
| /awesome-app/staging/vpc-id           | vpc-abcd1234 |
+---------------------------------------+--------------+

3 results.
```

### `cli`

```bash
pstore --profile=staging cli /awesome-app
```

An example response looks like:

```bash
aws ssm put-parameter \
    --profile staging \
    --name "/awesome-app/staging/environment-name" \
    --type SecureString \
    --value "staging" \
;

aws ssm put-parameter \
    --profile staging \
    --name "/awesome-app/staging/project-name" \
    --type SecureString \
    --value "awesome-app" \
;

aws ssm put-parameter \
    --profile staging \
    --name "/awesome-app/staging/vpc-id" \
    --type SecureString \
    --value "vpc-abcd1234" \
;
```

### Filtering

Parameter Store’s built-in searching is pretty weak. If you have a key named `/awesome-app/staging/vpc-id`, passing `/awesome-app/` or `/awesome-app/staging/` as the `path` parameter will return results. However, passing `/awesome-app/staging/vpc` will return zero results. This is because Parameter Store can only match on path fragments delimited with a `/` character.

`pstore` goes beyond this with support for _filtering_ the results that come back from AWS on the client-side. It supports two filtering modes: _substring_ and _regular expression_.

* **substring** — Filtering is _case-insensitive_, and attempts to match partial strings in both the key name as well as the value.

  ```bash
  --filter vpc-id
  ```

* **regular expression** — [RE2 expressions](https://github.com/google/re2/wiki/Syntax) are supported (similar to PCRE). Attempts to match both the key name as well as the value.

  ```bash
  --regex "(?i)vpc-[0-9a-f]{8}"
  ```

## Developing

### Linter

First, install the linting dependencies.

```bash
make lint
```

## Troubleshooting

### error parsing regexp

> error parsing regexp: <error message>

Check your regular expression. The RE2 engine in Go is very similar to the PCRE expression engine supported by PHP, JavaScript, Ruby, and Perl. You can test your regular expression with [Regex Tester - Golang](https://regex-golang.appspot.com).

### ValidationException: Parameter path: can't be prefixed with "aws" or "ssm"

> ValidationException: Parameter path: can't be prefixed with "aws" or "ssm" (case-insensitive) except global parameter name path prefixed with "aws" (case-sensitive). It should always begin with / symbol.It consists of sub-paths divided by slash symbol; each sub-path can be formed as a mix of letters, numbers and the following 3 symbols .-_

You may have tried to use a `*` or some other wildcard character for your `<path>` argument. The _path_ argument does not support wildcards. Use `--filter` or `--regex` instead.

  [AWS CLI Tools]: https://aws.amazon.com/cli/
  [Glide]: https://glide.sh
