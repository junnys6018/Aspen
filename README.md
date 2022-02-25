# Aspen
![ci badge](https://github.com/junnys6018/aspen/actions/workflows/CI.yml/badge.svg)

Aspen is a statically typed, multi-paradigm, interpreted toy programming language created for learning purposes.

```
fn list(generator fn(i64)i64, n i64) void {
    for (let i i64 = 0; i <= n; i = i + 1) {
        print generator(i);
    }
}

fn square(n i64) i64 {
    return n * n;
}

list(square, 10);
```

## Building

1. Install [go](https://go.dev/dl/)
2. cd into `aspen`
3. run `go build`

## Downloading

Aspen is automatically built with every commit as a part of a GitHub action.

You can download the latest binary [here](https://nightly.link/junnys6018/aspen/workflows/CI/master/aspen.zip)
