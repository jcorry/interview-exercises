# S3 Compression

## The Problem

We receive large (20MB+) code files that have to be stored in S3 for record-keeping.  To minimize costs, we would like to store them in a compressed format.  To further minimize costs, we would like to offload this process onto low-memory hardware.  We get these files regularly and need the software that processes them to be expedient.  For simplicity, we have decided to use the gzip compression format as it offers the balance between speed/compression that we need.  Please write code that takes uncompressed input and writes compressed output and test(s) that verify its efficacy.  The interface requirements are:

- The upload manager to S3 takes an io.Reader as its argument (output from your code)
- The uncompressed data is provided to your code as an io.ReadCloser (input to your code)

You are encouraged to mock out these inputs and outputs to simplify your solution

## Solution

Provide a Compressor interface that will provide a method that takes an `io.Reader` as its single argument and returns an `io.Reader` containing gzip compressed data.
