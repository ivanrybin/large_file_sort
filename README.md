# How to build

## Generator

```bash
$ ./build_generator
```

Creates binary file `./file_gen`

## Sorter

```bash
$ ./build_sorter
```

Creates binary file `./file_sort`

# How to run

## Generator

```
-c, --count <count>   - random strings count
-l, --max-len <len>   - generated strings max length
-o, --output <file>   - output file
-a, --alpha           - not random reversed alphabetic strings (without --max-len flag)
```

```bash
# random strings
$ ./file_gen -c 1000 -l 2000 -o random.txt

# alphabetic strings
$ ./file_gen -a -c 1000 -o alpha.txt
```

Random strings generates from `[A-Za-z_0-9]` alphabet.

## Sorter

```
-i, --input <path>   - input file to sort
-o, --output <path>  - output file (can be made the same as --input)
```

```bash
# sort large file
$ ./file_sort -i input.txt -o output.txt
```

## Algorithm

Merge sort with 4 file buffers and maximum 2 strings in RAM in one moment. 

![](https://github.com/ivanrybin/large_file_sort/tree/algo_pic.png)