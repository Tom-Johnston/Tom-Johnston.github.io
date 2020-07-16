---
title: The longest title on arXiv
excerpt: "A discussion about a potential title with my coauthors led to a question: What is the longest title for a maths paper?"
---
<div id="main-content" class="content" markdown="1">

# The longest title on arXiv

First, a caveat: the title is a slight lie because I'm only going to answer the question for papers where the primary category is maths. I'm also not going to do anything fancy about handling LaTeX shortcuts; I'm going to count the characters in the title field as is and go off that. Now the caveats are out of the way, I'll skip straight to the answer.

> C.F. Gauss' Pr\"azisionsmessungen terrestrischer Dreiecke und seine \"Uberlegungen zur empirischen Fundierung der Geometrie in den 1820er Jahren (C.F. Gauss' high precion measurements of terrestrial triangles and his thoughts on the empirical foundations of geometry in the 1820s)

While this is the longest title going by the title field on arXiv, really it is the title of the paper given once in German and then again in English. The second longest title doesn't suffer from this problem, although it does waste a lot characters on LaTeX.

> The consistent reduction of the differential calculus on the quantum  group $GL_{q}(2,C)$ to the differential calculi on its subgroups and  $\sigma$-models on the quantum group manifolds $SL_{q}(2,R)$,  $SL_{q}(2,R)/U_{h}(1)$, $C{q}(2\|0)$ and infinitesimal transformations

The third paper is also worthy of a mention (even if it is mathematical physics) since the title isn't long for a silly reason. It doesn't contain a translation, lots of LaTeX shortcuts or a reference. This is just a long title, plain and simple.

> De Rham-Hodge-Skrypnik theory. A survey of the spectral and differential geometric aspects of the De Rham-Hodge-Skrypnik theory related with Delsarte transmutation operators in multidimension and its applications to spectral and soliton problems. Part 1

The fourth paper is back into mathematics (`math.GM`) and absolutely littered with technical terms.

> Application of the Method of Approximation of Iterated Stochastic Ito Integrals Based on Generalized Multiple Fourier Series to the High-Order Strong Numerical Methods for Non-Commutative Semilinear Stochastic Partial Differential Equations

Obviously, as a combinatorialist, I looked up the longest in `math.CO` as well. This time the trick was giving the value of a constant, both exactly and to 23 decimal places.

> A Motivated Rendition of the Ellenberg-Gijswijt Gorgeous proof that the Largest Subset of $F_3^n$ with No Three-Term Arithmetic Progression is $O(c^n)$, with $c=\root 3 \of {(5589+891\,\sqrt{33})}/8=2.75510461302363300022127...$

The second title (given below) is better, although giving a second, light-hearted title in brackets is maybe cheating. The third longest has a thanks at the end of the title (and is otherwise quite short), the fourth has some long numbers and a decent amount of LaTeX shortcuts, the fifth seems to be mostly LaTeX, and it isn't until the eighth that the title doesn't seem too silly. 

> Sketch of a Proof of an Intriguing Conjecture of Karola Meszaros and Alejandro Morales Regarding the Volume of the $D_n$ Analog of the Chan-Robbins-Yuen Polytope (Or: The Morris-Selberg Constant Term Identity Strikes Again!)

As a bonus, we can also check the shortest titles. The shortest in all of maths is just `Q`, a single letter (and quite hard to beat). The shortest title of a paper with `math.CO` as a category is `Sects`, but this increases to `x-area` if you require combinatorics to be the primary category.

## About

I can't remember exactly how this question came about but it was in a conversation with some coauthors (Alex Scott, Carla Groenland and Jane Tan) when we were starting to look at a problem for planar graphs and graphs embeddable on higher-dimensional surfaces. Planar graphs were the obvious starting point and I remarked that they certainly make the catchier title, or at least the shorter title. Somehow this lead to the question of how long the titles of papers get and the most obvious (and easily accessible) dataset is arXiv.

Pulling all the metadata from arXiv is a relatively straightforward process as thye provide an [API](https://arxiv.org/help/oa/index) explicitly designed for bulk data access. The standard is pretty well-documented and it didn't take long to write a basic harvester (code at the bottom of the page) to pull and parse the data in the `arXiv` format. There are a few things worth noting about the data though. 

 - ArXiv must have used a different scheme for naming the categories before and there are a few categories which appear to have been renamed. There are also a couple of aliases which I mapped to their `math.XX` version.
    - `alg-geom: math.AG`
    - `dg-ga: math.DG`
    - `funct-an: math.FA`
    - `q-alg: math.QA`
    - `math-ph: math.MP`
    - `cs.IT: math.IT`
 - ArXiv seems to enforce a maximum line length in the titles (and other fields like abstacts)and will insert line breaks in the titles (explicitly it seems to replace a space with a new line followed by two spaces). I was very confused when I first plotted the distribution of title lengths and found no titles with lengths 73 and 74. To get round this I've simply replaced all instances of `\n  ` with ` ` and I'm hoping this hasn't hit any false positives.

## Distribution of title lengths

Since I had spent the time downloading the metadata to find the longest title, I decided I might as well also look at the distribution. There are far too many categories to nicely display the distributions for all of them on a graph, so I settled on four categories: all of maths, combintorics, a selection of pure areas and a selection of applied areas. This is only a blog post so I didn't invest much time pigeonholing the areas and just went down the list assigning the area to either pure, applied or neither. 

Pure: `math.AC, math.AG, math.AT, math.CT, math.FA, math.GN, math.GR, math.GT, math.KT, math.LO, math.NT, math.RA, math.RT`

Applied: `math.CA, math.DS, math.MP, math.NA, math.OC, math.ST`

{% include 2020-07-16-the-longest-title-on-arxiv/titles.svg %}

(I like this style of graph even if I should actually be using a different one.)

The titles are definitely longer in general for the applied ares compared to pure ones, and combinatorics seems to be even more concise. I've not read many applied papers but I generally think of them as having more descriptive titles compared to other areas, so this was pretty much what I expected.

## Number of authors

One of the advantages of the `arXiv` format is that the authors are listed separately and this makes it easy to check how many authors are on each paper. The record is 60 on the paper `Problems on invariants of knots and 3-manifolds`. Second place (`Proceedings of Workshop AEW10: Concepts in Information Theory and Communications`) is only just over half of that with 34. In general it looks like information theory seems to be taking most of the top spots with 5 of the top 7 papers and 11 of the top 25. It is also worth mentioning that combinatorics gets fourth place with `Supercharacters, symmetric functions in noncommuting variables, and related Hopf algebras` which has 28 authors.

{% include 2020-07-16-the-longest-title-on-arxiv/authors.svg %}

Unsurprisingly, the vast majority of papers have only a few authors and almost all of them have less than 10. At least some areas in applied maths order the authors by contribution, so there is less incentive to keep the number of authors lower. I've also heard that in some ares the supervisor is always listed even if they didn't contribute much directly to the paper, and this problem contributes a bit to the applied papers having more authors in general.

## Number of papers

It is a pretty terrible measure of the output of a mathematician (maybe even worse than citations), but we can also look at the number of papers the average mathematician churns out. The results won't necessarily be very accurate because different papers might have different ways of formatting the same name. For example, there are separate authors `Béla Bollobás`, `Bela Bollobas`, `B. Bollobas` and `B. Bollobás`, which are presumably all the same person.  It is also silly to include people who have never published in an area, so the graph below only includes authors with at least one paper of the appropriate type.

{% include 2020-07-16-the-longest-title-on-arxiv/papers.svg %}

Somewhat surprisingly, over half of the people who have published a combinatorics paper have only published 1. I expected it to be quite high from people who have strayed into combinatorics from their primary areas, but I expected to see more people with 2 or 3 papers. The distribution does have quite a long tail though, so there are actually quite a lot of people who have more than say 5 papers (nearly 17%). A more instructive graph is the following, showing the proportion of authors who have a given number of papers or more. There are too many factors which could be effecting these proportions so I'm not going to draw any conclusions from this graph. It still makes for a pretty much though.

{% include 2020-07-16-the-longest-title-on-arxiv/geq-papers.svg %}

## The most prolific authors on arXiv

Finally, we can end with a leaderboard based on the number of papers arXiv.

<div markdown="1" class="centering">


|Name| Number of papers| Primary area|
|---|---|---|
|Saharon Shelah|740|Logic|
|H. Vincent Poor| 438| Information Theory|
|Delfim F. M. Torres| 334|Optimization and Control|
|Indranil Biswas|312|Algebraic Geometry|
|Terence Tao|269|--|
|Rui Zhang|242|Information Theory|
|Yuval Peres| 230| Probability|
|Benny Sudakov|204| Combinatorics|
| Xueliang Li| 199| Combinatorics|
|Mohamed-Slim Alouini|195| Information Theory|

</div>

And for combinatorics.

<div markdown="1" class="centering">

|Name | Number of combinatorics papers|
|---|---|
|Benny Sudakov| 201|
|Xueliang Li| 198|
|Doron Zeilberger| 144|
|William Y. C. Chen| 131|
|David R. Wood| 121|
|Alan Frieze| 105|
|Jacob Fox|104|
|Michael Krivelevich| 104| 
|Bojan Mohar| 88|
|Toufik Mansour| 87|

</div>

<link rel="stylesheet" type="text/css" media="screen" href="/css/solarize.css">

## Code

You can download the code I used to pull the data from arXiv and generate the statisitics below. The graphs are saved from some poorly hacked together d3.js code that is a collage of different snippets I found online and I'm not prepared to put my name on it.

<a href="/downloads/2020-07-16-the-longest-title-on-arxiv.zip" download>Download the code</a>

</div>
 