---
title: Enumerating permutation graphs
excerpt: This blog post quickly compares three different ways of enumerating the permutation graphs on \\(n\\) vertices. By using a canonical deletion approach, we can extend [OEIS A123448](https://oeis.org/A123448) from \\(n = 11\\) to \\(n = 14\\).
math: true
code: true
---
<div id="main-content" class="content" markdown="1">

This blog post is very similar to the post on [enumerating circle graphs](/blog/2020-10-04-enumerating-circle-graphs.html) from a few weeks ago.

# Enumerating permutation graphs

Given a permutation \\(\\sigma = (\\sigma_1, \\sigma_2, \\dots, \\sigma_n)\\) the corresponding *permutation graph* is the graph with vertex set \\(V = \\{1, 2, \\dots, n\\}\\) and, for \\(i < j\\), there is an edge \\(ij\\) if and only if \\(\\sigma_i > \\sigma_j \\). Much like the circle graphs that we [enumerated](/blog/2020-10-04-enumerating-circle-graphs.html) a few weeks ago, the best [list](http://www.jaist.ac.jp/~uehara/graphs/) I could find is not in the most obvious format, and the  [OEIS entry](https://oeis.org/A123448) doesn't get that far. In fact, up until 07/09/2020 the list only went up to \\(n = 10\\), and you can do \\(n = 11\\) by simply iterating over all permutations!

There are three obvious ways of enumerating circle graphs:

1. Iterating over all permutations, computing the corresponding graphs and removing all non-isomorphic graphs. 
2. Iterating over all graphs and counting how many are permutation graphs.
3. A canonical deletion approach.

### 1. Checking all permutations

Constructing the graph from a permutation can trivially be done by iterating over all pairs \\((i,j)\\) with \\(i < j\\) and checking if \\(\\sigma_i > \\sigma_j\\), so the first method takes only ~20 lines of code using `mamba`. There are *only* \\(n!\\) permutations, which is pretty good for a brute force enumeration, and this method isn't a terrible idea. Besides, I don't know if the number of permutation graphs is \\(o(n!)\\) anyway. It took only 59.7s to enumerate all the permutation graphs on 10 vertices, and  11 vertices only takes 13m 5s. That suggests that maybe 12 vertices can be done in just a few days, and maybe 13 vertices in a few weeks.  

```go
func checkAllPermutations(n int) int {
	iter := itertools.Permutations(n)

	graphs := make(map[string]struct{})
	for iter.Next() {
		v := iter.Value()
		g := graph.NewDense(n, nil)

		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if v[i] > v[j] {
					g.AddEdge(i, j)
				}
			}
		}
		ci := graph.CanonicalIsomorph(g)
		g6 := graph.Graph6Encode(graph.InducedSubgraph(g, ci))
		graphs[g6] = struct{}{}
	}

	return len(graphs)
}
```

One of the reason that permutation graphs are interesting is that they are [perfect graphs](https://en.wikipedia.org/wiki/Perfect_graph), and so some of the classical NP-hard problems are polynomial on permutation graphs. In general, the problem of graph isomorphism is not known to be NP-complete or to be solvable in polynomial time, but an \\(O(n^3)\\) algorithm for permutation graphs was found by [Colbourn in 1981](https://onlinelibrary.wiley.com/doi/epdf/10.1002/net.3230110103). While this algorithm is fast asymptotically, I don't know how it compares to a generic algorithm in practice, and the specialised algorithm may or may not speed up this method. Since it seems unlikely that this will beat a canonical deletion approach, I've not tried implementing it.

### 2. Checking all graphs

For the second and third methods we will need a recognition algorithm for permutation graphs. The best result is an \\(O(n+m)\\) algorithm from McConnell and Spinrad in 2011, but we'll stick to the simpler \\(O(n^3)\\) method given by [Pnueli, Lempel and Even in 1971](https://www.cambridge.org/core/journals/canadian-journal-of-mathematics/article/transitive-orientation-of-graphs-and-identification-of-permutation-graphs/C9E11C37FE3FCB7A88C333E74C9EED01). This is a simple algorithm which exploits the fact that a graph \\(G\\) is a permutation graph if and only if both \\(G\\) and \\(\\overline{G}\\) are transitively orientable, and it doesn't take long to implement.

We could also make a recognition algorithm using the recognition algorithm for circle graphs we used when enumerating them. A circle graph is a permutation graph if and only if it admits an *equator* i.e. an additional chord which intersects all other chords. This phantom equator would be a dominating vertex and we can equivalently ask that the graph plus a dominating vertex is a circle graph. This is incredibly quick to implement, but very slow to run.

While the number of graphs is less than \\(n!\\) for very small values of \\(n\\), this breaks down at \\(n = 10\\) and the number of graphs which we check and reject starts to massively outnumber the permutation graphs. However, checking if a small graph is a permutation graph is quite quick and even at \\(n=10\\) it seems to be quicker to check all graphs than to check all permutations. But for \\(n=11\\), there are 169x more unlabelled graphs than unlabelled permutation graphs and this is much slower than checking all permutations.

<div markdown="1" class="centering">

| \\(n\\) | Permutations | Unlabelled Graphs | Permutation Graphs |
|--------:|-------------:|------------------:|-------------------:|
|       1 |           1x |                1x |                  1 |
|       2 |           1x |                1x |                  2 |
|       3 |         1.5x |                1x |                  4 |
|       4 |         2.2x |                1x |                 11 |
|       5 |         3.6x |                1x |                 33 |
|       6 |         5.1x |              1.1x |                142 |
|       7 |         6.5x |              1.3x |                776 |
|       8 |         7.1x |              2.2x |               5699 |
|       9 |         7.2x |              5.4x |              50723 |
|      10 |         6.9x |             22.9x |             524572 |
|      11 |         6.6x |              169x |            6037518 |

</div>

### 3. Canonical Deletion

Removing any vertex from a permutation graph gives a permutation graph, so we can plug the recognition algorithm into our basic canonical deletion algorithm with no modifications! Since the complement of a permutation graph is again a permutation graph (corresponding to the reverse permutation), we could speed this up by only counting graphs with at least \\(\\frac{1}{2} \\binom{n}{2} \\) edges and then adding in their complements, and this seems to save about a third. It cuts down the time to enumerate the graphs on 10 vertices to about 4s and the graphs on 11 vertices to about 1m, which isn't bad, but isn't going to allow us to get much further. It's also worth noting that the canonical deletion currently takes a vertex of minimum degree as the vertex to remove so we have to count graphs with at least \\(\\frac{1}{2} \\binom{n}{2} \\) edges and not at most. This makes it a bit fragile (it will break if we switch to removing a vertex of max degree), and the conditions used in the canonical deletion are currently undocumented/subject to change.

### Results

The table below gives the time to run each of the three methods once. Some of the early \\(n\\) are very small and would really need to be run many times to get a more reliable estimate, so we can pretty much ignore the first half of the table. Each trial was run with a single goroutine and given one hour.

<div markdown="1" class="centering">

| \\(n\\) | Checking all permutations | Checking all graphs | Canonical Deletion |
|--------:|--------------------------:|--------------------:|-------------------:|
|       1 |                     1.2ms |               3.4ms |              916µs |
|       2 |                     826µs |               856µs |              853µs |
|       3 |                     1.0ms |               745µs |              772µs |
|       4 |                     988µs |               943µs |              778µs |
|       5 |                     2.0ms |               1.1ms |              1.1ms |
|       6 |                     7.8ms |               2.1ms |              2.2ms |
|       7 |                    52.0ms |               8.5ms |              8.0ms |
|       8 |                     428ms |              72.6ms |             55.2ms |
|       9 |                      5.1s |                1.2s |              527ms |
|      10 |                     59.7s |               39.3s |               6.0s |
|      11 |                    13m 5s |             49m 16s |             1m 28s |
|      12 |                        -- |                  -- |            27m 46s |
|      13 |                        -- |                  -- |                 -- |

</div>

By splitting the canonical deletion across 20 goroutines and leaving the code running for a few hundred hours, we can count the number of circle graphs on 12, 13 and 14 vertices, as given in the table below. It's worth noting that you can clearly see that running this code in parallel for \\(n = 12\\) added quite a bit of CPU time, but this isn't really a fair comparison. For one, the parallel code below was outputting all the graphs and writing them to file, instead of just counting the graphs. The code was also run on a different, shared computer that, as well as potentially being slower (or faster) per core, was pinned to 100% CPU usage with many other jobs. In fact, the repeated work done by every process should only be up to \\(n = 8\\) which only takes 55ms, and this should only add a second or so to the overall time. 

<div markdown="1" class="centering">

| \\(n\\) | Number of circle graphs | CPU Time |
|--------:|------------------------:|---------:|
|      12 |                75912033 |   1h 11m |
|      13 |              1029974969 |      44h |
|      14 |             14974215412 |     569h |

</div>

### Code

- <a href="/downloads/2020-10-25-perm.zip" download>Enumerate permutation graphs</a>