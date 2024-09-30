# IP Address Unique Counter

This project is a solution to a technical test task for a job application, aimed at efficiently counting the number of unique IPv4 addresses in a large file.

## Problem Description

You are given a simple text file that contains IPv4 addresses, with each address on a new line, like so:

```
145.67.23.4
8.34.5.23
89.54.3.124
89.54.3.124
3.45.71.5
...
```

The file size is **unlimited**, and it can be **hundreds of gigabytes**. Your task is to calculate the number of **unique IP addresses** in this file while minimizing both **memory usage** and **execution time**.

The naive solution—reading the file line by line and storing each line in a hash set—is not efficient enough for large files. This implementation provides a more optimized approach to handle large files using multiple threads and memory-efficient techniques.

## Solution

The implementation supports two modes:
1. **Simple mode**: Single-threaded processing with a straightforward approach.
2. **Multi-threaded mode** (default): The file is divided into chunks and processed concurrently using multiple threads.

### Features:
- **Multi-threaded processing**: Faster execution by dividing the file into chunks and processing them in parallel.
- **Efficient memory usage**: IP addresses are stored in a bit map to reduce memory overhead.
- **Customizable chunk size**: You can specify how much memory each thread uses.
- **Dynamic thread management**: The number of threads can be manually set, or it can default to the number of available CPU cores.

## Usage

You can run the program with the following command-line arguments:

### Basic Usage:

Use this to run in a simple way multithreaded processing

```
--file-path=D:\\ip_addresses\\ip_addresses
```

### Extended Usage:

```
--file-path=D:\\ip_addresses\\ip_addresses --mode=multi --chunk-size-mb=5 --threads-number=20
```

### there are two modes:
simple - In this mode, the file is processed sequentially using a single thread, and a bit map is used to count unique IP addresses.
```
--mode=simple
```
multi (active by default) - This mode divides the file into chunks, processes each chunk in parallel, and aggregates the results.
```
--mode=multi
--chunk-size-mb=5 - Defines the amount of memory (in MB) allocated per thread.
--threads-number=20 - Specifies the number of threads to use. If the number exceeds the available CPU cores, it defaults to the number of cores.
```

### Performance:

On a machine with the following specifications:

- Processor: Intel Core i9
- RAM: 16 GB
- Storage: SSD
- Operating System: Windows 11

The program processes a 120 GB uncompressed file in ~1 minute 30 seconds. (2GB RAM is used)