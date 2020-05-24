# CS 3345 Project 5 Submission

## Student details

**Name:** Dien Tran  
**NetID:** DKT180000  
**Section:** CS 3345.003

## Project structure

```bash
<extracted directory>
|-- QuickSorter.java   # QuickSorter class
|-- Main.java          # Main driver class
`-- README.txt         # this file

0 directories, 3 files
```

## Dev environment used

**Code editor**: Visual Studio Code with Java Extension Pack  
**JDK version**: 14.0.1  

## Build instructions

```bash
# cd into extracted directory
$ javac *.java
$ java Main ARRAY_SIZE REPORT_FILE UNSORTED_FILE SORTED_FILE
```

where `ARRAY_SIZE` is the array input size, `REPORT_FILE` is the report file's filename, `UNSORTED_FILE` is the filename of the file to contain the randomly generated list, and `SORTED_FILE` is where to save the output of the quicksort. If already exists, these will get overwritten.

Should work in *both* Mac and Windows. Use Terminal or Command Prompt, respectively.

## Notes

The `timedQuickSort` method is implemented *iteratively*. The reason is that were we to use recursion, since each recursive call takes an `ArrayList`, we would be forced to pass *a copy of the original list*. Obviously this will backfire, so the author went with iteration and a stack as the method to divide and conquer.

## Investigating behavior of different inputs

### Randomly generated inputs

First, we will give, as input, randomly generated lists of length 500, 1000, 5000, and 1000000. This is to see how the algorithm's processing time scales when the inputs scale linearly and quadratically.

We present the sample reports here:

- For input size 500:

```none
Array Size = 500

FIRST_ELEMENT = PT0.001554S

RANDOM_ELEMENT = PT0.0006651S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.0045763S

MEDIAN_OF_THREE_ELEMENTS = PT0.0009264S
```

- For input size 1000:

```none
Array Size = 1000

FIRST_ELEMENT = PT0.0023885S

RANDOM_ELEMENT = PT0.0013496S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.0055136S

MEDIAN_OF_THREE_ELEMENTS = PT0.0017659S
```

- For input size 5000:

```none
Array Size = 5000

FIRST_ELEMENT = PT0.0082463S

RANDOM_ELEMENT = PT0.0051131S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.0094541S

MEDIAN_OF_THREE_ELEMENTS = PT0.0041973S
```

- For input size 1000000:

```none
Array Size = 1000000

FIRST_ELEMENT = PT0.5627743S

RANDOM_ELEMENT = PT0.5047064S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.5205783S

MEDIAN_OF_THREE_ELEMENTS = PT0.4893934S
```

For input size 500, `RANDOM_ELEMENT` wins because of how likely a random element turns out to be a good pivot; otherwise, there isn't a clear winner overall. We inspect the median-based pivoting methods, and the worse performing one is `MEDIAN_OF_THREE_RANDOM_ELEMENTS`. This is understandable since it is costly to generate 3 random numbers on each partition. `MEDIAN_OF_THREE_ELEMENTS` mitigate this by not using random numbers and uses fixed indices for values to compute the median out of.

The first element is very unlikely to be a good pivot as input size scales; this is why `FIRST_ELEMENT` is the worst performing one at size 1000000. `RANDOM_ELEMENT` fares not much better because the probability of randomly picking a good pivot does not scale favorably. This leaves us with `MEDIAN_OF_THREE_ELEMENTS` as our best overall method.

### Already sorted list

Now we generate lists of length 500, 1000, 5000, and 1000000. This is to see how the algorithm's processing time scales when the inputs scale linearly and quadratically.

- For input size 500:
  
```none
Array Size = 500

FIRST_ELEMENT = PT0.0069101S

RANDOM_ELEMENT = PT0.000545S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.0048529S

MEDIAN_OF_THREE_ELEMENTS = PT0.0006274S
```

- For input size 1000:

```none
Array Size = 1000

FIRST_ELEMENT = PT0.0190851S

RANDOM_ELEMENT = PT0.0006705S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.0050196S

MEDIAN_OF_THREE_ELEMENTS = PT0.0007879S
```

- For input size 5000:

```none
Array Size = 5000

FIRST_ELEMENT = PT0.0499277S

RANDOM_ELEMENT = PT0.0078954S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.0110975S

MEDIAN_OF_THREE_ELEMENTS = PT0.0076694S
```

- For input size 1000000:

> ***suppressed because `FIRST_ELEMENT` takes too long***

It is clear `FIRST_ELEMENT` is the bottleneck because it always partitions the list in the most uneven way possible (all other elements are greater than or equal to the pivot). For other strategies however, `RANDOM_ELEMENT` and `MEDIAN_OF_THREE_ELEMENTS` tie over the position of best pivot strategy for sorted lists.

### Almost sorted list (last 90% of elements are sorted)

We sort the last 90% of the input list only. Same input size sequence: 500, 1000, 5000, 1000000.

- For input size 500:

```none
Array Size = 500

FIRST_ELEMENT = PT0.0016214S

RANDOM_ELEMENT = PT0.0006493S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.0058214S

MEDIAN_OF_THREE_ELEMENTS = PT0.0013367S
```

- For input size 1000:

```none
Array Size = 1000

FIRST_ELEMENT = PT0.0017257S

RANDOM_ELEMENT = PT0.0011765S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.0071442S

MEDIAN_OF_THREE_ELEMENTS = PT0.0019744S
```

- For input size 5000:

```none
Array Size = 5000

FIRST_ELEMENT = PT0.0153S

RANDOM_ELEMENT = PT0.0055388S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.0095033S

MEDIAN_OF_THREE_ELEMENTS = PT0.0042522S
```

- For input size 1000000:

```none
Array Size = 1000000

FIRST_ELEMENT = PT0.5824554S

RANDOM_ELEMENT = PT0.3702574S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT0.3826917S

MEDIAN_OF_THREE_ELEMENTS = PT0.2331649S
```

Again `FIRST_ELEMENT` performs worst (except for input size 500), but because as far as picking the first element as pivot is concerned that this case is the same as a randomly generated list, this holds true as well; and as input size grows `MEDIAN_OF_THREE_ELEMENTS` becomes the most sustainable pivot strategy. We see if our prediction holds true by observing what happens when input size is 5000000:

```none
Array Size = 5000000

FIRST_ELEMENT = PT3.3723378S

RANDOM_ELEMENT = PT2.9135165S

MEDIAN_OF_THREE_RANDOM_ELEMENTS = PT2.7908045S

MEDIAN_OF_THREE_ELEMENTS = PT2.6848917S
```

Our prediction holds true.
