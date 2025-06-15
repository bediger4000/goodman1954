# Some Practical Techniques in Serial Number Analysis

Looking at Section 5, _An Application_ of
Leo A. Goodman (1954) Some Practical Techniques in Serial
Number Analysis, Journal of the American Statistical Association, 49:265, 97-112

## Raw Data

Goodman's data is 31 serial numbers in 3 lines of text,

```
83, 135, 274, 380, 668, 895, 955, 964, 1113, 1174, 1210, 1344, 1387, 1414,
1610, 1668, 1689, 1756, 1865, 1874, 1880, 1936, 2005, 2006, 2065, 2157, 2220,
2224, 2396, 2543, 2787
```

Goodman figures out two basic things about the data:

1. How large was the initial purchase of serial numbered items?
2. How close to a uniform distribution is this group of serial numbers?

### Estimate of how many items were purchased and numbered

Goodman's formula for estimating how many items got serial numbered is:

p = d(k+1)/(k-1) - 1

- **p** is "total production", the number of items that really got serial numbered
- **d** is the numerical difference between maximum and minimum serial number
- **k** is the number of serial numbers

p = (2787 - 83)(31+1)/(31-1) - 1 = 2883.3

### How close to a uniform distribution

## Simulation
