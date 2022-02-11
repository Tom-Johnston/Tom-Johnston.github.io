---
title: Optimal Wordle strategies
excerpt: "Wordle has been taking the world by storm, but what is the best way to play? This blog post looks at a few different strategies and the optimum play for a few different versions of \"best\""
math: true
---
<div id="main-content" class="content" markdown="1">

# Optimal Wordle strategies

Unless you have been living under a rock for the past few weeks, you've probably heard of [Wordle](https://www.powerlanguage.co.uk/wordle/) ([Wikipedia](https://en.wikipedia.org/wiki/Wordle)), the simple word game that has been taking the world by storm. I'm by no means the first to "solve" Wordle (and it isn't helped by the multi-week lag time between doing something and writing about it), but I fancied writing down what I did anyway and maybe there are a few new bits here.

## Greedy strategies

### Mutual information

Since Wordle is all about gathering information, the obvious first approach is to greedily choose the word which maximises the *[mutual information](https://en.wikipedia.org/wiki/Mutual_information)*. To quote Wikipedia "it quantifies the "amount of information" obtained about one random variable by observing the other random variable", i.e. it quantifies how much we learn about the hidden word from observing a guess. Let \\(X\\) be the hidden word, which we assume is chosen uniformly at random from the list of 2315 possible answers. When we guess a word \\(w\\), we observe a pattern, which we will call \\(Y_{w}\\), and which is a function of \\(X\\). The mutual information between \\(Y_{w}\\) and \\(X\\) quantifies the amount we learn about the hidden \\(X\\) from the pattern \\(Y_{w}\\) observed by guessing \\(w\\), and we want to choose \\(w\\) which maximises this quantity. 

The mutual information can be expressed in terms of entropy (and conditional entropy) with the formula
\\(I(X;Y_{w}) = H(X) - H(X|Y_w)\\). The entropy of the uniform distribution over the set of answers \\(\\mathcal{A}\\) is easy: \\(\\log_2(|\\mathcal{A}|)\\). Given we observe a pattern \\(P\\) from guessing a word \\(w\\), \\(X\\) is uniformly distributed over the subset of answers which give the pattern \\(P\\), which we denote by \\(\\mathcal{A}_P\\). The mutual information is therefore

$$I(X;Y) = \log_2(|\mathcal{A}|) - \sum_{P : |\mathcal{A}_P| \neq 0} \log_2(|\mathcal{A}_P|)$$

where the sum is over the possible patterns \\(P\\) for which \\(\\mathcal{A}_P\\) is non-empty. Computing this gives the following top 10 answers.

<div markdown="1" class="centering">

| Rank | Word  | Mutual information |
|------|-------|--------------------|
| 1    | SOARE | 5.89               |
| 2    | ROATE | 5.88               |
| 3    | RAISE | 5.88               |
| 4    | RAILE | 5.87               |
| 5    | REAST | 5.87               |
| 6    | SLATE | 5.86               |
| 7    | CRATE | 5.83               |
| 8    | SALET | 5.83               |
| 9    | IRATE | 5.83               |
| 10   | TRACE | 5.83               |

</div>

This strategy can be used at every step. We can compute the current list of possible answers and the answer is uniform over that list. We can then compute the word which maximises the mutual information with a uniform random variable over the new list, and guess that. (There is one edge case here: we might reach the point where we know the word but haven't guessed it. In this case, there is no entropy and every word has 0 mutual information with the answer, so we need to force the player to guess the answer.) This strategy guesses all words in at most 6 guesses and takes a total of 8418 guesses over all the words for an average solve in 3.64 guesses. Not too bad, but we can do better.

### Number of greens

Before we do better, let me mention another strategy I was curious to try, but always knew would do badly. This actually came after finding the optimum, but was motivated by actually playing the game. While working out that letters are not in the word is useful and the best strategies probably make extensive use of that, I find it easier to come up with the answer if I know the letters it contains, especially if I know their location. As a very naive strategy one could look to maximise the number of greens gained in the first word (with ties broken by the number of yellows). This obviously ignores that some greens are more useful than others and getting the first letter is often particularly useful, but this is not going to be a great strategy anyway. The mutual information seems like a reasonable ranking metric, even if it isn't perfect, and I've included it here just as an indication of which words are good. It's not perfect though and the "best" word from this list is probably SLANE. 

<div markdown="1" class="centering">


| Rank | Word  | Number of greens | Number of yellows | Mutual information rank |
|------|-------|------------------|-------------------|-------------------------|
| 1    | SAREE | 1575             | 2596              | 371                     |
| 2    | SOOEY | 1571             | 1696              | 5925                    |
| 3    | SOREE | 1550             | 2385              | 923                     |
| 4    | SAINE | 1542             | 2238              | 26                      |
| 5    | SOARE | 1528             | 2565              | 1                       |
| 6    | SAICE | 1512             | 2166              | 45                      |
| 7    | SEASE | 1510             | 2287              | 6582                    |
| 8    | SEARE | 1491             | 2614              | 404                     |
| 9    | SLANE | 1480             | 2301              | 20                      |
| 10   | SEINE | 1480             | 2076              | 1808                    |

</div>

It looks like lots of word begin with S and end with E then. We can use this method on every guess as well and it takes a total of 9006 guesses to guess all the words for an average of 3.89 guesses per word. This is again not bad, but there is a problem. It fails to guess 26 words within 6 guess and this would lose occasionally. 

<div markdown="1" class="centering">

<table>
        <tr>
          <th>Number of guesses required</th>
          <td>1</td>
          <td>2</td>
          <td>3</td>
          <td>4</td>
          <td>5</td>
          <td>6</td>
          <td>7</td>
          <td>8</td>
          <td>9</td>
        </tr>
        <tr>
          <th>Number of words</th>
          <td>0</td>
          <td>94</td>
          <td>703</td>
          <td>1023</td>
          <td>386</td>
          <td>83</td>
          <td>20</td>
          <td>5</td>
          <td>1</td>
        </tr>
    </table>

<!-- | Number of guesses required | Number of words |
|----------------------------|-----------------|
| 1                          | 0               |
| 2                          | 94              |
| 3                          | 703             |
| 4                          | 1023            |
| 5                          | 386             |
| 6                          | 83              |
| 7                          | 20              |
| 8                          | 5               |
| 9                          | 1               | -->

</div>

## The optimal strategy

The first step to working out an optimal strategy is deciding what you mean by optimal. How much better is guessing the answer in 2 guesses than in 3 guesses? Would you rather get it in a guaranteed 3 guesses, or a 50/50 chance of getting it in 2 or 4? We'll take the obvious choice here. An *optimal strategy* is a (deterministic) strategy which always guesses the answer within 6 guesses and takes the minimum number of guesses on average. Let \\(N_{\\mathcal{A}}\\) be the expected number of guess to when the answer is chosen randomly from the set \\(\\mathcal{A}\\) when playing the optimum strategy (for that set of answers). Using the law of total expectation, we can write down the expectation conditional on the pattern \\(Y_{w}\\) we see when guessing \\(w\\).

$$ \mathbb{E}[N_{\mathcal{A}}] = \sum_{P} \mathbb{E}[N_{\mathcal{A}} | Y_{w} = P ] \mathbb{P} ( Y_{w} = P) = \sum_{P} \mathbb{E}[N_{\mathcal{A}_P}] \frac{|\mathcal{A}_p|}{|\mathcal{A}|}$$

Assuming we pick a word which isn't completely useless, the set \\(\\mathcal{A}_P\\) is a strict subset of the set \\(\\mathcal{A}\\). The case where there is a single word is clearly trivial and so, if we know the optimum for all smaller sets, we can find the optimum for all the set \\(\\mathcal{A}\\) by checking each word. In practice, we don't know the optimal solution for all smaller sets and we calculate it on demand. To get the calculation done in a reasonable time there are a few tricks worth employing such as caching the minimum number of guesses required for each set and calculating lower bounds on the number of guesses required, but I'm not going to delve into it. It's worth noting that I didn't have to include the restriction that the strategy always guesses it in at most six guesses as the strategy I found does this for free. Since other people have already done this calculation, I'm not wasting the computer time doing it again and I'll just copy down the top 10 from [Alex Selby](http://sonorouschocolate.com/notes/index.php?title=The_best_strategies_for_Wordle).

<div markdown="1" class="centering">

| Rank | Word  | Total | Avg   | Mutual Information Rank |
|------|-------|-------|-------|-------------------------|
| 1    | SALET | 7920  | 3.422 | 8                       |
| 2    | REAST | 7923  | 3.422 | 5                       |
| 3    | CRATE | 7926  | 3.424 | 7                       |
| 3    | TRACE | 7926  | 3.424 | 10                      |
| 5    | SLATE | 7928  | 3.425 | 6                       |
| 6    | CRANE | 7930  | 3.425 | 31                      |
| 7    | CARLE | 7937  | 3.429 | 19                      |
| 8    | SLANE | 7943  | 3.431 | 20                      |
| 9    | CARTE | 7949  | 3.434 | 14                      |
| 10   | TORSE | 7950  | 3.434 | 37                      |

</div>

It is generally pretty quick to work out the optimum for a "good" word as it partitions the list of answers into many sets which are mostly fairly small and it is quick to solve the problem for small lists of words. The slow part is checking that the "bad" words can't somehow be done quickly. The first guess might only separate out a handful of words, but lower bounds are hard to come by and there are no guarantees this doesn't make the rest of the list easy. Instead, you'll end up wanting to compute the next best word, which isn't really that different to computing the best word in the first place. However, we can clearly see that the best words also do well at maximising the mutual information so you can get a pretty good idea of the best words by working down the list of words in order of decreasing mutual information.

### Hard mode

Wordle also has a hard mode in which you are restricted in the words you can guess. If you find a green letter, you can only guess words which have that correct letter in that place, and if you find a yellow letter, all future guesses must include the yellow letter. Note that you can repeat the position of a yellow letter and you can reuse grey letters, so it isn't quite as restrictive as it could be. Although the list of possible guesses changes at each level, you can still use a recurrence relation as above to find the optimum. Disappointingly, the optimum is still SALET, but interestingly you need to add the restriction that the strategy guesses everything within six guesses, and you don't get it for free. If you are allowed to take 7 guesses to guess the word, the minimum total number of guesses across all words is 8116 (for an average of 3.51 guesses per word), but if we restrict to winning strategies, this goes up to 8122. To add the condition that the strategy must always win, we can still use a recursive formula but we need to keep track of the number of of guesses already used and add some kind of penalty when it takes too many guesses. In fact, we might as well be more general and give each guess a cost. To make sure the winning strategy doesn't use 7 or more guesses, we can give the 7th guess a large cost of say 10,000. Since the best strategy takes 8122 guesses, this retroactively justifies that 10,000 was a large enough number. Interestingly, this seems to make the code quite a bit slower. This is probably because it is harder to pick words which are close to the lower bounds on the number of guesses.

<div markdown="1" class="centering">

| Rank | Word  | Total | Avg   | Mutual Information Rank |
|------|-------|-------|-------|-------------------------|
| 1    | SALET | 8122  | 3.508 | 8                       |
| 2    | LEAST | 8127  | 3.511 | 29                      |
| 3    | REAST | 8134  | 3.514 | 5                       |
| 4    | CRATE | 8143  | 3.517 | 7                       |
| 5    | TRAPE | 8144  | 3.518 | 61                      |
| 6    | SLANE | 8149  | 3.520 | 20                      |
| 7    | PRATE | 8151  | 3.521 | 42                      |
| 8    | CRANE | 8155  | 3.523 | 31                      |
| 9    | TEALS | 8160  | 3.525 | 70                      |
| 9    | TRAIN | 8160  | 3.525 | 151                     |

</div>

The lists are quite similar, but even good words in normal mode and can be bad in hard mode. For example, SLATE is joint 5th in normal mode with 7928 guesses, but it can't even win every time in hard mode. Suppose that the T and E are green while the S and A are yellow. There are 6 possible answers (BASTE, CASTE, HASTE, PASTE, TASTE, WASTE), and there isn't a guess that distinguishes between any of them unless it is correct. One of these must take 7 guesses in our strategy.


## Other definitions of best

We've taken the obvious definition of best above: a strategy which always wins and minimises the average number of guesses, but there are other sensible options. The strategy I found for SALET always uses at most 5 guesses, and you might wonder if that is optimal. Is it possible to always get the answer in 4 guesses, even if the average is worse. I believe the answer is no, but it raises the interesting question of what happens if you try to maximise the chance of guessing the answer in some small number of guesses (amongst winning strategies and breaking any ties by the total number of guesses). With 1 guess, we obviously need to guess something from the list of possible answers and then we are just playing a normal game. Looking at the table above, the best is TRACE (or CRATE) in normal mode and LEAST in hard mode. What about for other numbers?

<div markdown="1" class="centering">

| Target number of guesses | Best word | Number meeting target | Total | Best word (Hard mode) | Number meeting target (Hard mode) | Total (Hard mode) |
|--------------------------|-----------|-----------------------|-------|-----------------------|-----------------------------------|-------------------|
| 1                        | TRACE     | 1                     | 7926  | LEAST                 | 1                                 | 8127              |
| 2                        | TRACE     | 150                   | 8078  | SALET                 | 147                               | 8136              |
| 3                        | TRACE     | 1388                  | 7968  | REAST                 | 1223                              | 8167              |
| 4                        | TENOR(?)  | 2298                  | 8101  | ???                   | ???                               | ???               |
| 5                        | SALET     | 2315                  | 7920  | PALET                 | 2315                              | 8206               |
| 6                        | SALET     | 2315                  | 7920  | SALET                 | 2315                              | 8122              |

</div>

It is possible to guarantee a win in hard mode in at most 5 guesses, but you have to change from the word SALET to do it and you have to spend quite a few more guesses. I've run the code overnight for 4 guesses in normal mode and checked the top ~650 words, finding the word TENOR can guess all but 17 or them in at most 4 guesses.

## Good pairs of words

In general, whatever your starting word, you probably want to choose your second word based on the feedback from the first word. Instead, many people (myself included) like to start with two standard words unless the first word does particularly well. What is the best way to start if we always start with the same two words? I don't know the answer to this question and there are probably too many pairs of word which need checking to work out the optimum answer, but we can definitely come up with some good words. With 12,972 possible guesses and therefore 84,129,906 unordered pairs (the only time the order matters is when one is a possible answer and one isn't, in which case guess the possible answer first), I don't particularly even want to check the mutual information of all of these pairs. Instead, I've only considered the pairs which start with one of the [top 105 words](http://sonorouschocolate.com/notes/index.php?title=The_best_strategies_for_Wordle), ranked those pairs and checked the top 1000 pairs. This has given me the following top 11.

<div markdown="1" class="centering">

| Rank | Word 1 | Word 2 | Total | Avg   |
|------|--------|--------|-------|-------|
| 1    | PARSE  | CLINT  | 8325  | 3.596 |
| 2    | CRANE  | SPILT  | 8336  | 3.600 |
| 3    | CRINE  | SPALT  | 8339  | 3.602 |
| 3    | SPALT  | CRINE  | 8339  | 3.602 |
| 5    | SLANT  | PRICE  | 8344  | 3.604 |
| 6    | CRANE  | SLIPT  | 8346  | 3.605 |
| 6    | SLANT  | CRIPE  | 8346  | 3.605 |
| 8    | CRISE  | PLANT  | 8350  | 3.607 |
| 9    | PRASE  | CLINT  | 8351  | 3.607 |
| 9    | SOARE  | CLINT  | 8351  | 3.607 |
| 9    | CRINE  | PLAST  | 8351  | 3.607 |

</div>