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

```bash
make build
```

<asciicast>

## Running

```bash
pstore
```

<asciicast>

By default, `pstore` leverages the credentials from the [AWS CLI Tools].

### Pass the Key and Secret manually

Using the `list` command as an example, and assuming you've installed and configured the [AWS CLI Tools](https://aws.amazon.com/cli/) at some point:

```bash
pstore list \
    --key AKIAJHJBEXAMPLEW7VJA \
    --secret QIBAlwMNrzExampleG9iRf1ttflI0PDooExample
```

If you are using multiple profiles, you can use a specific profile by passing the `--profile` flag.

```bash
pstore list --profile production
```

## Parameters

For all commands, you can pass an _argument_ to the command line. This argument is the _path_ of the Parameter Store value that you want to see. The default value is `/`.

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

## Troubleshooting

### Regexp Error: error parsing regexp

> ```plain
> panic: Regexp Error: error parsing regexp: <error message>
> 
> goroutine 1 [running]:
> github.com/skyzyx/pstore/cmd.glob..func2.2(0xc4203e7b80, 0x2, 0x2, 0x199c468)
>   /path/to/gocode/src/github.com/skyzyx/pstore/cmd/list.go:59 +0x1ac
> github.com/skyzyx/pstore/cmd.arrayFilter(0xc4200f0800, 0xe4, 0x100, 0x15d3928, 0x81, 0xc4200f0800, 0x80)
>   /path/to/gocode/src/github.com/skyzyx/pstore/cmd/root.go:167 +0xe1
> github.com/skyzyx/pstore/cmd.glob..func2(0x1979980, 0xc420122e80, 0x0, 0x2)
>   /path/to/gocode/src/github.com/skyzyx/pstore/cmd/list.go:54 +0x498
> github.com/spf13/cobra.(*Command).execute(0x1979980, 0xc420122e20, 0x2, 0x2, 0x1979980, 0xc420122e20)
>   /path/to/gocode/src/github.com/spf13/cobra/command.go:766 +0x2c1
> github.com/spf13/cobra.(*Command).ExecuteC(0x1979be0, 0x1979e40, 0x1979d30, 0xc42015bf48)
>   /path/to/gocode/src/github.com/spf13/cobra/command.go:852 +0x30a
> github.com/spf13/cobra.(*Command).Execute(0x1979be0, 0xc42002c178, 0x0)
>   /path/to/gocode/src/github.com/spf13/cobra/command.go:800 +0x2b
> github.com/skyzyx/pstore/cmd.Execute()
>   /path/to/gocode/src/github.com/skyzyx/pstore/cmd/root.go:66 +0x2d
> main.main()
>   /path/to/gocode/src/github.com/skyzyx/pstore/main.go:21 +0x20
> ```

Check your regular expression. The RE2 engine in Go is very similar to the PCRE expression engine supported by PHP, JavaScript, Ruby, and Perl. You can test your regular expression with [Regex Tester - Golang](https://regex-golang.appspot.com).

  [AWS CLI Tools]: https://aws.amazon.com/cli/
