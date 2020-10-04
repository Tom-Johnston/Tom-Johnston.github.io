---
excerpt: A graph \\(G\\) is said to be *induced saturated for \\(H\\)* if \\(G\\) contains no induced copy of \\(H\\), but adding or removing any edge creates an induced copy. We give some graphs which are induced saturated for \\(P_5\\).
math: true
---

<div class="content" markdown="1">

<div class="flex-share">
    <img src="/images/2020-05-22-induced-saturation-for-paths/K%7Di%5BBLZR%60lfY.svg" alt="K}i[BLZR`lfY" title="K}i[BLZR`lfY">
    <img src="/images/2020-05-22-induced-saturation-for-paths/circulant(19,1,2,3,5,7,8).svg" alt="R~]y|^VYxjrVrVXjuYyrVjL^UY~UYw" title="R~]y|^VYxjrVrVXjuYyrVjL^UY~UYw">
</div>

# Induced Saturation for \\(P_5\\)

<span class="hint">Joint work with Marthe Bonamy, Carla Groenland, Tash Morrison and Alex Scott.</span>

A graph \\(G\\) is said to be *induced saturated for \\(H\\)* if \\(G\\) contains no induced copy of \\(H\\), but adding or removing any edge creates an induced copy. While this definition holds for any graph \\(H\\), research has so far focused on paths and, to a lesser extent, cycles.

Paths are nearly completely understood: the cases \\(P_2\\) and \\(P_3\\) are easily handled by the empty and complete graphs respectively, the case \\(P_4\\) was shown to be impossible by Martin and Smith in 2012, and the case of \\(P_n\\) where \\(n \geq 6\\) has recently been solved by Vojtěch Dvořák. This leaves just the case \\(P_5\\) which we close here.

Armed with some code which reads in graphs from `stdin` and some which computes the number of paths of a given length, it didn't take long to create a simple program to check if graphs are induced saturated for paths. In fact, the first version was cobbled together essentially in the time it took my collaborators to check the Petersen graph is induced saturated for \\(P_6\\). To find an induced saturated graph for \\(P_5\\) (and many other small paths), we first need a list of graphs to check. All small graphs is a reasonable starting point, and it yields two graphs for \\(P_6\\), the Petersen graph and a graph on 11 vertices, but running this code for all graphs on 12 vertices is already infeasible. The next list we tried was all connected vertex-transitive graphs, which we could do up to 31 vertices. This provides plenty of examples for small graphs, including 5 for \\(P_5\\). 

- ``K}i[BLZR\`lfY``
- `R}~t}VF{TdtefRUnkzzZu]Zg^{xnmo`
- `T~~~}rbxS^q}mNqdulInvVNxlnre~X^nsz|{`
- `Z}v\tmn}^r}u|jzs|xmFdozQ{Q|wYr~d|~l]E^~wl~~@z~{]X~b\m^dlzrrW`
- `]v~vVN~~~~^y~rvoVwN@fwJ^_ff?v^?zi^?^V~nx~}~bN~{Dn~xC~~w@n~{W\~~I^oNz|~@^vo`

The first two examples are pictured above. Only the second one seems to have any sort of special structure; it is the circulant graph on 19 vertices and jumps 1, 2, 3, 5, 7 and 8. I don't know of any nice argument why these graphs are induced saturated for \\(P_5\\), so checking that these work is left as an exercise for the reader.

Despite providing plenty of examples, checking all vertex transitive graphs didn't reveal many obvious patterns, especially with no simple way of working out the structure of the examples (if any). The next thing to try was families of graphs. This gave us some families which are probably induced saturated such as the Kneser graphs \\(K(n,2)\\) for \\(P_6\\) (as proved by Cho, Choi and Park) and the bipartite Kneser graphs \\(H(n, 2)\\) for \\(P_{12}\\), but it was slow to come up with, generate and check each family of graphs. With a bit of effort, Mathematica can be made to output its list of graphs in Graph6 along with all the names associated with each graph. This offers a quick way to check many different families and is definitely something worth trying again in the future, even if it didn't give what we were looking for this time. 

We did find a couple of constructions which were induced saturated for infinitely many \\(n\\), both of which are cubic, but that's for another time. Finally, by running our code for nearly 144 hours, we can say that the folded cube of order 7 isn't \\(P_k\\)-induced saturated for any \\(k\\). Both of these results answer questions of Cho, Choi and Park.
</div>
