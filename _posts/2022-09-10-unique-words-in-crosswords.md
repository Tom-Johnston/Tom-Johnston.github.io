---
title: Unique words in Guardian cryptic crosswords
excerpt: "How often do you get new words in crosswords? How many clues should you expect to get? What's the longest clue? I recently downloaded over 7,000 Guardian cryptic crosswords to answer questions just like these. "
math: true
---
<div id="main-content" class="content" markdown="1">

{% include figure.html src="/images/2022-09-10-unique-words-in-crosswords/average-grid.svg" caption="The mean 15x15 grid" %}

# Unique words in Guardian cryptic crosswords

How often do you get new words in crosswords? How many clues should you expect to get? What's the longest clue? I recently downloaded over 7,000 Guardian cryptic crosswords to answer questions just like these. My original motivation was actually the first one. When clicking around on my website a friend stumbled across the [anagram solver/crossword solver](/blog/2021-08-30-writing-an-anagram-solver) I wrote. Of course, despite having a list of over 310,000 words, the first (and only) word they tried was not there, although this might be because `BARCODE` should be two words. We then got talking about using the answers from previous crosswords as a list of words. But how often do you come across a word that hasn't been clued before? As someone who doesn't do that many crosswords, it wasn't obvious to me what the answer should. Of course, the words you might want to use a crossword solver for are exactly the rare words which are unlikely to have cropped up before, but it is still an interesting question. (For the record, `BAR CODE` appears six times, `BARCODES` appears once, but `BARCODE` never appears as an answer.)

The Guardian has a nice interactive solver where you can fill in the crossword and check your answers. This means that there must be data there for the scraping and fortunately, it is very handily encoded in a JSON object already. I have mostly done the Quiptic crosswords (so that I actually stand a chance of getting a decent proportion of the answers in a reasonable amount of time), but there are less than 2000 of them available and I decided to stick to the actual cryptics for now. The Guardian has been publishing cryptic crosswords for over 100 years and has published nearly 30,000 of them, but the first one I could find online is [21,620](https://www.theguardian.com/crosswords/cryptic/21620). This gives us over 7,000 crosswords spanning over 20 years to use as our data.

## Gathering the data

The URLs for the crosswords follow a nice pattern: if you want to look at cryptic number 1234, then the URL is either `cryptic/1234` or, if it is a prize cryptic on a Saturday, `prize/1234`. As much as it would be nice to work out which URL we need to query for each crossword, it is much easier to just try `/cryptic/1234`, check if it returns a 404 and then try `prize/1234` if necessary. Once you are on the page, the data is conveniently stored in a JSON string in the (only) `data-` tag. In an effort to be polite, I scraped one crossword a second and then stored both the original HTML (in case I messed up and needed the original data) and the JSON (to actually work with). The HTML comes in at a hefty `5.2GB`, but the JSON is a fairly lightweight `48MB`. 

Some of the crosswords have little games associated with them and we'll have to skip these. For example, in [23,047](https://www.theguardian.com/crosswords/prize/23047) all but three of the down clues are put in backwards to celebrate Australia day. Maybe this is more fun than it sounds, but it is awkward for the purposes of checking if words have been used before and we'll skip it. Some of them weren't adapted to the online format and were available only as a PDF, and it is simply too much effort to work with these/type them up. In some of them the clues aren't linked to the answers, which causes it's own set of problems. Some of them just do crazy things and you don't even input actual words (see e.g. [21,783](https://www.theguardian.com/crosswords/cryptic/21783)). I've had a look through the list of special instructions and hopefully removed all the crosswords we don't want.

*Garbage in, garbage out.* Fortunately, this isn't quite going to the be case here, but there are definitely some problems with the data (even on the crosswords without any silly games) and everything here should be taken with a grain of salt. Sandwiched between [Cryptic 21,830](https://www.theguardian.com/crosswords/cryptic/21830) on Fri 25th Feb 2000 and [Cryptic 21,832](https://www.theguardian.com/crosswords/cryptic/21832) on Mon 28th Feb 2000 is a Saturday prize crossword as one would expect, but it is not down as Prize 21,831 as one would expect, but as [Prize 21,381](https://www.theguardian.com/crosswords/prize/21831). However, the main problem is when a single clued answer is split across several places/answers in the grid. To handle this, each answer in the grid is in a group and a group represents the answer to a clue. Unfortunately, this data is far from accurate. The most common issue is when a clued answer is split across three or more grid answers and instead of being in a single group, there are multiple groups each consisting of the first grid answer and one of the remaining ones. It's also worth mentioning that the separator data is also dodgy. While there is separator data in the JSON, it is inconsistent on whether it lists a separator where an answer is split across different grid answers. There is usually some kind of separator here, but I've found examples where this isn't the case. I decided the simplest way to split answers into words was to parse the splitting given in the clue and use that to split the answer. It also helpfully highlights problems where the lengths don't match up and these can then be fixed or, more likely, ignored. With over 1000 errors, I'm not fixing them all manually and I'll remove those clues/answers from the data, except for a few interesting ones like the longest answer.

## The numbers

Let's start with the headline graph.

{% include 2022-09-10-unique-words-in-crosswords/prop-words-new.svg %}

This graph may not be the clearest, but it certainly shows that there are new words fairly often and that the words that have come up before is not a good list of words for an anagram/crossword solver. 

Across the 7155 crosswords I analysed there were 200,992 answers, although only 71,910 distinct answers. About half of these appeared once (37,362) while the most frequent answer came up 63 times. 

### Most popular answers

{% include 2022-09-10-unique-words-in-crosswords/answer-freq.svg %}

We can start asking lots of questions. For example, we might be interested in the answers that come up the most often. Some of these words are expected; `STYE` feels like it crops up a lot in crosswords (but almost never in real life), but some of these definitely surprised me.

<div markdown="1" class="centering">


| Rank | Answer | Frequency | Rank | Answer  | Frequency |
|------|--------|-----------|------|---------|-----------|
| 1    | EXTRA  | 63        | 13   | IDEAL   | 40        |
| 2    | ISLE   | 61        | 13   | ACHE    | 40        |
| 3    | USED   | 49        | 13   | BLUE    | 40        |
| 3    | STUD   | 49        | 18   | NIECE   | 39        |
| 3    | ECHO   | 49        | 18   | ISSUE   | 39        |
| 6    | ERATO  | 48        | 18   | ARENA   | 39        |
| 7    | STYE   | 47        | 21   | ADONIS  | 38        |
| 8    | EDGE   | 46        | 22   | ORANGE  | 37        |
| 9    | ARCH   | 43        | 22   | ADDRESS | 37        |
| 9    | IDLE   | 43        | 22   | STAY    | 37        |
| 11   | ANON   | 42        | 25   | STIR    | 36        |
| 11   | STUN   | 42        | 25   | UNIT    | 36        |
| 13   | ADIEU  | 40        | 25   | ITEM    | 36        |
| 13   | ONSET  | 40        | 25   | EVENT   | 36        |

</div>


We can also check the most popular words, but it isn't particularly enlightening. The small connecting words crop up a lot, although maybe `UP` being so high is a little surprising.

<div markdown="1" class="centering">

| Rank | Word | Frequency |
|------|------|-----------|
| 1    | THE  | 2129      |
| 2    | IN   | 1361      |
| 3    | UP   | 1208      |
| 4    | A    | 999       |
| 5    | ON   | 987       |
| 6    | OF   | 967       |
| 7    | TO   | 693       |
| 8    | AND  | 662       |
| 9    | OUT  | 605       |
| 10   | OFF  | 524       |

</div>

### Number of clues

{% include 2022-09-10-unique-words-in-crosswords/num-clues.svg %}

Over the 6,726 crosswords where I haven't encountered any errors reading the clues, there is an average of 28.09 clues per crossword. The number of clues is very clearly concentrated around 28, although there are a few outliers. The lowest number is 19 in [27,142](https://www.theguardian.com/crosswords/prize/27142) (where one might feel quite short-changed), while there are two with a huge 68 clues, [28,010](https://www.theguardian.com/crosswords/prize/28010) and [28,632](https://www.theguardian.com/crosswords/prize/28632). Although unsurprisingly these are both 23 x 23 grids instead of the usual 15 x 15. Trimming the graph at 36 clues makes it much more readable and only misses 29 of the crosswords.

### The longest answers 

{% include 2022-09-10-unique-words-in-crosswords/answer-lengths.svg %}

The longest answer is a phrase from Shakespeare's *Julius Caesar* and comes in at a huge 82 letters long.

> Araucaria in May, reportedly tied in with dynamite in Welsh river, displays sign about fellow that volunteers knowledge, having chewed the fat with John Duke: Shed holds a medal for musical composition: let's go! (5,2,1,4,2,3,7,2,3,5,5,2,3,5,5,2,2,7)<br>
> `THERE IS A TIDE IN THE AFFAIRS OF MEN WHICH TAKEN AT THE FLOOD LEADS ON TO FORTUNE`

It's a bit a of a drop to the second one at 71, but this is still incredibly long.

> Quentin Crisp's note retrieved high up rank to low, where he'd pay even less - compete? Just ridiculous! (5,4,2,4,3,7,4,4,4,2,4,5,3,7)<br>
> `NEVER KEEP UP WITH THE JONESES DRAG THEM DOWN TO YOUR LEVEL ITS CHEAPER`

There are also a handful at 60+, although quite a few of them are incorrect in the JSON.

### The longest clues

{% include 2022-09-10-unique-words-in-crosswords/clue-lengths.svg %}

This is one where you can clearly see the errors in the data. There are 517 clues of length 4, all of which consist of a space followed by `(d)` for some digit `d`. These clearly aren't proper clues, but I'm again not investing the time to fix all of these mistakes and they will just have to pollute the graph. Overall there is still a nice bell curve though.

The longest clue with 252 characters is:

> "We turn for a short time from the topics of the day to commemorate, in all love and reverence, the genius and virtues of John Milton, the 12, the 10, the 14, the g24 of ___ 11, the 16 and the 19 of ___ 28" (Macaulay, making the most of his length) (7)

The next best is not far behind with 250 (and the longest answer):

> Araucaria in May, reportedly tied in with dynamite in Welsh river, displays sign about fellow that volunteers knowledge, having chewed the fat with John Duke: Shed holds a medal for musical composition: let's go! (5,2,1,4,2,3,7,2,3,5,5,2,3,5,5,2,2,7)

And the third (with 175 characters):

> Harem without a doubt (first and last and centre): come fly to chaos, keeping killer now endlessly in French and English church with cult work of the ’70s (3,3,3,3,2,10,11)


### The average crossword grid

{% include figure.html src="/images/2022-09-10-unique-words-in-crosswords/average-grid-21.svg" caption="The mean 21x21 grid" %}


Most cryptic crosswords are based on 15 x 15 grids and this holds true for the Guardian as well. Of the 7155 I've got the data for, there are 27 on 21 x 21 grids, 7 on 23 x 23 grids and a single 20 x 20 grid (which breaks the convention that the grids have odd dimensions). To create the average grids I started with a black square and then overlaid white squares. The opacity of the overlaid squares is equal to the number of times each square appears as a white square divided by the total number of grids. 

I had to fix a couple of crosswords to get their grids to actually work, but there was one I wasn't sure how to fix. Cryptic [22482](https://www.theguardian.com/crosswords/cryptic/22482/print) seems to have a very odd shape and the clue lengths don't add up. There don't seem to be any obvious fixes either, but maybe solving the crossword will reveal the correct structure.

{% include figure.html src="/images/2022-09-10-unique-words-in-crosswords/average-grid-23.svg" caption="The mean 23x23 grid" %}
</div>
