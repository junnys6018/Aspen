const examples = [
    {
        name: 'Hello, World!',
        code: 'print "Hello, 世界!";',
    },
    {
        name: 'Recursion',
        code: `fn factorial(n i64) i64 {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
}

for (let i i64 = 0; i <= 10; i = i + 1) {
    print itoa(i) + "! = " + itoa(factorial(i));
}
`,
    },
    {
        name: 'Fibonacci',
        code: `let a i64 = 0;
let b i64 = 1;

while (a < 10000) {
    print a;
    let temp i64 = a;
    a = a + b;
    b = temp;
}
`,
    },
    {
        name: 'Closures',
        code: `fn makeCounter() fn()void {
    let i i64 = 0;
    fn count() void {
        i = i + 1;
        print i;
    }
    return count;
}

let counter fn()void = makeCounter();
counter();
counter();
counter();
`,
    },
    {
        name: 'First Class Functions',
        code: `fn list(generator fn(i64)i64, n i64) void {
    for (let i i64 = 0; i <= n; i = i + 1) {
        print generator(i);
    }
}

fn square(n i64) i64 {
    return n * n;
}

list(square, 10);
`,
    },
    {
        name: 'Fizzbuzz',
        code: `for (let i i64 = 1; i < 100; i = i + 1) {
    if (i % 15 == 0) {
        print "fizzbuzz";
    } else if (i % 3 == 0) {
        print "fizz";
    } else if (i % 5 == 0) {
        print "buzz";
    } else {
        print i;
    }
}
`,
    },
];

export default examples;
