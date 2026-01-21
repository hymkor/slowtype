# slowtype

`slowtype` is a small command-line utility that behaves like `cat`, but intentionally outputs data **slowly**, inserting delays between chunks of input.

It is mainly designed for testing tools that read from **stdin**, especially to verify that they behave correctly when input arrives gradually or temporarily stops.

## Motivation

When developing interactive or streaming tools, it is easy to accidentally assume that input arrives immediately and continuously.

`slowtype` helps test scenarios such as:

* Input arriving little by little
* Temporary pauses in input
* Programs that must not block UI or key input while waiting for stdin
* Behavior under slow pipes or delayed producers

This tool was originally created to test stdin handling of tools like:

* CSV editor ([Csvi])
* Binary viewer ([binview])

…but it can be useful for any program that reads from standard input.

## Usage

```bash
slowtype [options] [FILE]
```

If `FILE` is omitted, input is read from **stdin**, just like `cat`.

## Options

```
-b uint
    Insert delay after reading the specified number of bytes

-kb uint
    Same as -b, but in kilobytes

-mb uint
    Same as -b, but in megabytes

-ms uint
    Sleep duration in milliseconds (default: 100)

-hang
    Sleep after all input has been output
```

## Examples

### Read from a file (cat-like behavior)

```bash
slowtype large.txt
```

Outputs the contents of `large.txt`, inserting a delay after each line (default behavior).
### Use as a slow pipe

```bash
slowtype large.txt | your-application
```

This is useful to test how an application behaves when input arrives slowly.

### Delay output every 1 KB

```bash
slowtype -kb 1 large.txt
```

### Control delay duration

```bash
slowtype -kb 4 -ms 500 large.txt
```

Outputs data in 4 KB chunks, waiting 500 ms between each chunk.

### Simulate stalled or slow input for another tool

```bash
slowtype -b 10 large.csv | csvi
```

This is useful for testing tools that read from stdin and must remain responsive even when input is delayed.

### Keep the process alive after output

```bash
slowtype -kb 1 -hang large.txt
```

Useful for testing behavior when input ends but the process does not immediately exit.

## Typical Use Cases

* Testing stdin handling of CLI tools
* Simulating slow or chunked input
* Verifying non-blocking or asynchronous input handling
* Debugging interactive terminal applications
* Reproducing edge cases that are hard to trigger with `cat`

[Csvi]: https://github.com/hymkor/csvi
[binview]: https://github.com/hymkor/binview
