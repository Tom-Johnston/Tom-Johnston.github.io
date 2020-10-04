---
title: Enumerating circle graphs
excerpt: "There are 22576188846 unlabelled circle graphs on 13 vertices but how do we know that? This blog post looks at a few different ways of generating small circle graphs from the simplest ways up to a canonical deletion method."
math: true
code: true
---
<div id="main-content" class="content" markdown="1">

<figure markdown="1">
{% include circle-graphs.svg %}
<figcaption markdown="1">
A circle representation of the graph \\(P_5\\). A corresponding double occurrence word is `ACBEDECDAB`.
</figcaption>
</figure>

# Enumerating circle graphs

After first hearing about circle graphs not that long ago, I looked them up and I wasn't able to find a simple enumeration of them. The best I found was an enumeration in the [Database of circle graphs](http://www.ii.uib.no/~larsed/circle/) which gives one graph in each LC orbit, but I didn't know what an LC orbit was nor how easy/hard it is to generate a complete list of graphs from there. I already have some code which generates all non-isomorphic graphs using the [canonical deletion method](https://computationalcombinatorics.wordpress.com/2012/08/13/canonical-deletion/) and (like the rest of my code) it could do with an example or two, so I set about enumerating circle graphs. Unfortunately, I currently host this blog on Github Pages and the lists of graphs are too big to host here.

However, this isn't going to be a complete waste of time; the [OEIS entry](http://oeis.org/A156809) only goes up to \\(n = 12\\) and we can go one further.

Although we'll be aiming to enumerate circle graphs with the canonical deletion method, we'll start with the simplest method we can think of. In this case, it will be by iterating over double occurrence words. This doesn't require the ability to recognise circle graphs so we can build it out of pre-existing functions very quickly, and it can compute the first 6 numbers in a little over a minute. There are a couple of easy modifications which we can use to speed it up and we can push the enumeration to \\(n = 10\\) without much effort. Given a (non-trivial) function `isCircleGraph` which recognises circle graphs, we will try the obvious way of counting circle graphs, simply iterate over all graphs and check if they are a circle graph. Finally, we'll try using the canonical deletion method.

The table below shows the amount of time it took to compute the number of circle graphs with each method. Everything was done with a maximum of one goroutine and each call was given an hour to complete. Although the methods do get considerably quicker as we progress, we quite quickly hit the point where they can do \\(n = 10\\) in well under an hour, but none of them can do \\(n = 11\\).

<div markdown="1" class="centering">

| \\(n\\) | Double occurrence words | Topological Sorts |  Orbits |  Check all | Canonical Deletion |
|--------:|------------------------:|------------------:|--------:|-------:|-------------------:|
|       1 |                   121µs |            32.7µs |    68µs | 50.5µs |              134µs |
|       2 |                   154µs |            76.1µs |   112µs |  175µs |             87.1µs |
|       3 |                   947µs |             250µs |  84.8µs |  212µs |              112µs |
|       4 |                  19.9ms |             678µs |   501µs |  198µs |              211µs |
|       5 |                   1.01s |            13.0ms |  2.47ms | 1.05ms |              461µs |
|       6 |                  1m 11s |             119ms |  13.4ms | 2.08ms |             2.17ms |
|       7 |                      -- |             1.85s |   207ms | 15.9ms |             16.1ms |
|       8 |                      -- |             28.2s |   2.85s |  340ms |              289ms |
|       9 |                      -- |             8m 4s |   44.0s |  11.0s |              6.85s |
|      10 |                      -- |                -- | 12m 31s | 10m 8s |             2m 55s |
|      11 |                      -- |                -- |      -- |     -- |                 -- |

</div>

## Circle graphs

A *circle graph* is the intersection graph of a set of chords on the circle. Given a circle with chords labelled \\(1, \dots, n\\), the corresponding graph \\(G\\) has vertices \\(\{1, \dots, n\}\\) and there is an edge \\(ij\\) if and only if the chord \\(i\\) and the chord \\(j\\) intersect. Another way to view circle graphs is via *double occurrence words* i.e. words where each letter in the given alphabet appears exactly twice. Given a double occurrence word on the alphabet \\(\{1, \dots, n\}\\), let the corresponding circle graph be the graph with vertices \\(\{1, \dots, n\}\\) and an edge \\(ij\\) if and only if the \\(i\\)s and \\(j\\)s are interlaced i.e. form the substring \\(ijij\\) or \\(jiji\\). As simple discrete structures, double occurrence words are much easier to work with on a computer and they provide a convenient certificate for being a circle graph.

### A naive enumeration

Given the functions we already have in [`mamba`](https://github.com/Tom-Johnston/mamba), the easiest way of enumerating all circle graphs is to try all double occurrence words, convert each one into a graph, compute the canonical isomorph of each graph and only keep the unique graphs. There is a lot of wasted computational work this way, but it isn't a bad way to get the smallest graphs under our belt straightaway.

The [`itertools`](https://github.com/Tom-Johnston/mamba/tree/master/itertools) package has an iterator for permutations of a multiset and it can be used to trivially generate all double occurrence words in just 8 lines of code!

```go
freq := make([]int, n)
for i := range freq {
	freq[i] = 2
}
iter := itertools.MultisetPermutations(freq)
for iter.Next() {
	v := iter.Value()
	//Do something with the double occurrence word.
}
```

We need to add a little bit of code to convert the double occurrence word into a graph, to find the canonical isomorph of the graph and to remove duplicates. Converting the double occurrence word into a graph is a relatively simple function, but is still more code than I want to drop into a blog post. The completed code can be found attached to the bottom of the blog post so I'll just include a simplified version here. It takes in `n`, the number of vertices in the circle graphs we want to enumerate, and returns the number of non-isomorphic circle graphs. It wouldn't be hard to extract the list of graphs from it, but we'll stick to the number for now.

```go
func checkDOW(n int) int {
	//Iterate over all double occurrence words.
	freq := make([]int, n)
	for i := range freq {
		freq[i] = 2
	}
	iter := itertools.MultisetPermutations(freq)

	//Use a map to act as a set.
	graphs := make(map[string]struct{})

	for iter.Next() {
		//Compute the graph corresponding to v.
		v := iter.Value()
		g := graphFromDOW(v)

		//Compute the canonical isomorph and encode it in graph6 format.
		ci := graph.CanonicalIsomorph(g)
		g6 := graph.Graph6Encode(graph.InducedSubgraph(g, ci))
		//Add the encoded version of the graph to the set.
		graphs[g6] = struct{}{}
	}
	return len(graphs)
}
```

Some of the details are hidden in the function `graphFromDOW` which creates a graph from a double occurrence word, but that's a relatively simple piece of code and the whole method didn't take long to implement. It provides an easy way of counting the graphs when \\(n\\) is very small, but it is already taking over an hour when \\(n = 7\\). To hit our goal of counting the graphs on 13 vertices we'll need something better. 

Of course, we can use the results from this code to help check/debug the more complicated methods later on, so it wasn't a complete waste of time.

### Improving the naive enumeration

The above method is very quick to implement and it gives us the counts for very small values of \\(n\\), but it doesn't make any attempt to avoid repeating graphs. Since we are only interested in unlabelled graphs, we are free to relabel the vertices however we want. The simplest way of relabelling the vertices in the double occurrence word is such that the letters appear (for the first time) in ascending order. 
This simple observation saves us a massive \\(n!\\) in the number of graphs we have to check!

While we can check if the letters appear in the correct order for each double occurrence word and skip the double occurrence words where this isn't the case, it is much better to enforce the condition in the iterator. An iterator that did exactly this would be very specialised, but a small change will allow us to use one already in `itertools`. Instead of \\(1\\) appearing twice in the double occurrence word, we label the first occurrence with \\(1\\) and the second occurrence with \\(n + 1\\) and consider permutations of \\(\\{1, \dots, 2n\\}\\). Because of this we need to enforce that \\(i \\) comes before \\(i + n\\) for all \\(i \leq n\\), and ensuring the letters are in the correct order means ensuring \\(i\\) comes before \\(j\\) whenever \\(i < j \leq n\\). This is called a topological sort and The Art of Computer Programming gives the algorithm we after. 

The code to count circle graphs using topological sorts is essentially the same as the code which uses all double occurrence words, and we only need to change our definition of `iter` to the following.

```go
less := func(i, j int) bool {
	if i < j && j < n {
		return true
	}
	if i%n == j%n {
		return i < j
	}
	return false
}

iter := itertools.TopologicalSorts(2*n, less)
```

While switching to topological sorts already gives a massive speedup, there is another easy trick we can try. Cyclically shifting the double occurrence word doesn't change the corresponding (unlabelled) graph, and we only need to consider one of the words from this orbit. Given a double occurrence word \\(W\\), we can try each cyclic shift, relabel the letters so they appear in order and check if it is less than \\(W\\) lexicographically. The code to do this is less than 30 lines, but I don't have much to say about it so I'll omit it for now.

The table below shows how many words we check relative to the number of circle graphs, and you can clearly see the huge amount of progress these simple changes have made. We've gone from trying a billion times more words than graphs to just 15 times more! In fact, if we were given just an hour of CPU time to compute graphs, we'd already have reached as far as we'd get with the other two methods below. This is something worth bearing in mind with exhaustive searches like this: sometimes increasing \\(n\\) by 1 is still infeasible even if you make your code 10 times faster.

There are a couple of drawbacks to this method which only become more apparent as you go higher. The first is the memory usage. We are currently storing the entire list of graphs in a map in memory, not great when you are talking about \\(2^{30}\\) graphs. We could start storing the list of graphs on the disk, but that is going to slow things down massively and will still require many gigabytes of space. The other problem is running the code in parallel. There should be some way of splitting up the topological sorts into a small number of parallel cases, but it isn't obvious how balanced they will be and we still have to merge the lists at the end. Ideally, we want a way of counting the graphs without having to store the previous graphs...

<div markdown="1" class="centering">

| \\(n\\) | Double occurrence words | Topological sorts | Orbits | Circle graphs|
|---------|-----------------------:|------------------:|-------:|-------------:|
|        1|                      1x|                 1x|      1x|             1|
|        2|                      3x|               1.5x|      1x|             2|
|        3|                   22.5x|              3.75x|   1.25x|             4|
|        4|                    229x|               9.5x|    1.6x|            11|
|        5|                   3335x|                28x|    3.1x|            34|
|        6|                  48600x|                68x|    5.9x|           154|
|        7|                 696000x|               138x|     10x|           978|
|        8|                8610000x|               213x|     13x|          9497|
|        9|               97700000x|               269x|     15x|        127954|
|       10|             1100000000x|               302x|     15x|       2165291|

</div>


### The opposite approach

The above approach generates a load of circle graphs and then removes the non-isomorphic graphs, but we can also do the opposite. Let's try starting with a list of non-isomorphic graphs and check which are circle graphs. It's easy to generate a list of graphs using the `graph/search` package so we only need a way of recognising circle graphs, not that this is entirely trivial. Circle graphs can be [recognised in \\(O(n^2)\\) time](https://www.sciencedirect.com/science/article/abs/pii/S0196677484710121), although the algorithm is split across two papers and too complicated to be worth implementing just for a blog post. Of course, there are no guarantees that the algorithm is the quickest in practice and, for the small graphs we consider, the asymptotics don't matter much anyway. Instead we will use the [first polynomial time algorithm](https://www.sciencedirect.com/science/article/pii/0012365X85901177), which runs in \\(O(n^7)\\) but is considerably easier to implement.

First note that a disconnected graph is a circle graph if and only if each connected component is a circle graph, and we only need a test for connected graphs. A connected graph with edge set \\(E\\) is a circle graph if and only if the following system of equations has a solution over \\(\mathbb{F}\_2\\). The \\(n(n-1)/2\\) variables are \\(x_{u,v}\\) where \\(u\\) and \\(v\\) are vertices.

<p>
\[
\begin{aligned}
x_{uv} + x_{vu} &= 1 \text{ for } uv \in E \\
x_{uv} + x_{uw} &= 0 \text{ when } uv, uw \not \in E \text{ and } vw \in E\\
x_{uv} + x_{uw} + x_{vw} + x_{wv}  &= 1 \text{ when } uv, uw  \in E \text{ and } vw \not\in E\\
\end{aligned}
\]
</p>

It's easy to solve a system of linear equations using Gaussian elimination and, since we are working over \\(\mathbb{F}\_2\\), we don't even need to worry about numerical errors! We can use bitwise operations to do the arithmetic so a naive implementation isn't even that slow for small values. The starting matrix is very sparse which my implementation doesn't use in any way, and doing something better here should be one way of speeding up the code. Even my simple implementation takes nearly 150 lines of code so I won't paste the complete solution here and instead I'll use the function `isCircleGraph` as a black box. In fact, I'll cheat a little bit further by not reusing the data structures between calls of `isCircleGraph` and simplifying the function signature as well.

```go
func checkAllGraphs(n int) int {
	c := make(chan *graph.DenseGraph)
	go search.All(n, c, 0, 1)
	counter := 0
	for g := range c {
		if isCircleGraph(g) {
			counter++
		}
	}
	return counter
}
```

This method doesn't require us to store any of the previous graphs and, at least for \\(n \\leq 10\\), this method is quicker than the others considered so far. But the time looks to be increasing quickly as \\(n\\) increases, very quickly. By the time we hit \\(n = 11\\) there are more graphs than orbits, and who knows which computation is quicker, it will depend on how long it takes to recognise circle graphs compared to compute canonical isomorphs.

### A canonical deletion approach

Fortunately, circle graphs satisfy a useful property: removing any vertex from a circle graph leaves another (smaller) circle graph. When enumerating the non-isomorphic graphs we build them up by adding vertices one-by-one, and our useful property means we can never reach a circle graph from a non-circle graph. This means we can test the smaller graphs as we build up to \\(n\\) vertices and skip most of them. Think of building up all the graphs by traversing a tree with the \\(n\\) vertex graphs as the leaves \\(n\\) levels deep into the tree. To enumerate all the graphs we traverse the tree 
according to a DFS, and if we only want to enumerate circle graphs, we can prune any branch which starts from a non-circle graph.

Even better, all the hard work in implementing a canonical deletion algorithm with pruning is already done and we need just a few lines of code. We won't do any pre-pruning so we pass a function which always returns `false` first. We want to prune a graph whenever it is not a circle graph, the opposite of what our `isCircleGraph` returns, so we need a wrapper to negate it. With just 9 lines of code, some whitespace and a function that recognises circle graphs, we can count all the way up to \\(n = 10\\) in just a few minutes!

```go
func countCanonicalDeletion(n int) int{
	c := make(chan *graph.DenseGraph)

	go search.WithPruning(n, c, 0, 1, func(g *graph.DenseGraph) bool { return false }, func(g *graph.DenseGraph) bool { return !isCircleGraph(g) })

	counter := 0
	for range c{
		counter++
	}
	return counter
}
```

I've already mentioned that the function `isCircleGraph` is potentially far from optimal, and there might be even more inefficiency here. We have already checked the subgraph on \\(\\{1, \\dots, n-1\\}\\) is a circle graph, but we've completely ignored this later afterwards. We've already done work to put the matrix into row reduced echelon form and then we do this again afterwards. But maybe we are missing something even better. Maybe we can extract a double occurrence word for the subgraph and check if \\(G\\) can be formed by adding in an extra letter. Unfortunately (and somewhat predictably), this may not always be enough. For example, the double occurrence word `ACBEDECDAB` gives the graph \\(P_5\\), but there isn't any way of adding `F` to get the graph \\(C_6\\). It's not hard to see from the picture at the start of this post, there is no way to intersect the line `A` and the line `E` without intersecting the line `C`.

As well being quicker for single-threaded performance, this method also works well in parallel. We can prune every other branch at a given depth to (hopefully) split the work into roughly two equal parts. Since the majority of the work is done when \\(n\\) is large, the repeated work at the lower depths doesn't make much difference. This method again doesn't require us to store any graphs, so there aren't any worries about memory usage and we could leave it counting away for weeks at a time...


### Results

The maths department has a few sizeable computers which I can use for free and this can just about count as work, so I spread the work across 20 goroutines and left it running for a little over a week. For some reason the first time I tried to run it, it was unable to allocate memory and crashed. On the successful run, the maximum resident set size was less than 57 MiB, and even if this was per goroutine, it pales in comparison with the 755GiB RAM the computer has, so I'm going to blame another user and leave it at that. A bit of hope and some 4314 core hours (on an Intel Xeon Gold 6140) later, the second run completed with the answer: `22576188846`.

The code worked well enough, but this exercise highlighted a big issue: there isn't any way of saving and then resuming an enumeration. I will probably look into fixing this by creating some kind of `GraphIterator` struct which iterates over all graphs. This would allow a second goroutine to pause the calculation at regular intervals and save the state to a file in case the computer crashes. The graph could also be accessed from the struct instead of being sent over a channel, and we could avoid a load of unnecessary copying when all we want to do is count the graphs or output the Graph6. The name `graph/search` is left over from when I was actually searching for certain graphs and should probably be updated to something else, breaking this blog post immediately.

### Code

I've split the code into two parts, `circle-comp` contains all the different methods and is designed to compare the time taken between the different methods, and `circle` contains just the canonical deletion approach and outputs the graphs to `Stdout`. 

- <a href="/downloads/2020-10-04-circle-comp.zip" download>Compare the different methods</a>
- <a href="/downloads/2020-10-04-circle.zip" download>Enumerate circle graphs</a>

</div>
