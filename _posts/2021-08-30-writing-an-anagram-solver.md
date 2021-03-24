---
title: Writing a PWA to solve anagrams and more
excerpt: A couple of years ago I got fed up with looking at ads when doing something as simple as solving an anagram, so I decided to write my own [solver](/crossword-assistant).
---
<div id="main-content" class="content" markdown="1">

{% include figure.html src="/images/2021-08-30-writing-an-anagram-solver/icon-circle.svg" caption="The icon of the PWA — I'm actually pretty happy with this one." %}

# An anagram solver PWA built with Go

After doing most of the work a couple of years ago, I finally got round to actually finishing my PWA to solve anagrams and find words with missing letters in May. Check it out here: [Crossword assistant](/crossword-assistant).

Why has it taken so long for this blog post to be written? I actually wrote 90% of this post in May but the app stopped working. I hadn't changed any of the code and it was still working fine on my desktop, but it would crash on Chrome for Android. Weirdly, it wouldn't crash if I reloaded the page while USB debugging was active, but I did manage to get an error message out and the problem seemed to be casting an `int` inside the generated WASM. I put this on the back burner for a bit and concentrated on the thesis, next thing you know it's the end of August.

## The Why

Three things: 

1. **Completely free**: You should be able to look for anagrams and find words which match a pattern without having to look at a webpage/app full of ads.
2. **Powerful**: You should be able to look for anagrams with checking letters.
3. **Fun**: It might be kind of fun to make and should be a good chance to learn some useful skills.

The first of these is fixed fairly easily by using an adblocker or choosing the right website/app, but nothing solving the second problem showed up with a quick Google. While the three reasons above motivated me to make the solver, there are a few other things that would be nice to have (and other solvers already do):

1. **Fast**: Goes without saying.
2. **Cross-platform**: It doesn't need to work across every platform, but I'd like to use it on both my phone and my desktop.
3. **Work offline**: Now I'm no longer in the office, most of my crossword solving probably happens on a train where the WiFi is poor and the 4G coverage is variable.

## The How

There are three main components to an anagram solver:

1. A list of valid words.
2. The ability to search a list of words for matches.
3. A UI for inputting the queries to solve and displaying the results.

First things first, we need to work out a basic plan of action and the stack we want to use. This is a simple app which I'm only going to use occasionally (and chances are no one else will ever use it), so we want something that works well enough but can be thrown together pretty quickly. Given we also want it to be cross-platform, the easiest solution is to build something for the web and use modern browser features like caching to make it more like an app. The so-called progressive web apps (PWA) *should* work on almost all modern platforms and the small subset of features we want are widely supported. Of course, it's never that simple and the app doesn't currently work on iOS (and debugging iOS Safari requires either signing up to a service online or both an iPhone and a Mac). This also means that I can write the UI in HTML/CSS, the only UI "toolkit" I've used before, and I don't need to learn a toolkit I'm unlikely to use much again. Although I've written some of it before, I'm not a huge fan of Javascript and I definitely didn't want to have to write any complicated logic in it unless I had to. Fortunately, Go (my language of choice) now has WASM support and we can use this for the bulk of the searching. Of course, we can't avoid Javascript entirely and we'll use some Javascript for interacting with the UI and all the bits that make this a PWA.

- **UI**: HTML/CSS
- **Searcher**: Go compiled to WASM
- **Interactivity**: Javascript

### A list of valid words

This turns out to be a bit harder than you might expect. English is a somewhat poorly defined language and there are a lot of "words" that  may or may not be valid depending on who you ask. Since this is aimed mainly at solving cryptic crosswords, I decided it was generally safer to err on the side of false positives and there are plenty of words in there that I wouldn't really consider as valid. Its main use so far has been working on/failing at the Listener crosswords which seem to be mostly filled with words I have never heard of, and some are so obscure they have to specify an alternative dictionary as a reference. There are plenty of lists of words available online and you may even a decent sized list already on your computer (try `/usr/share/dict/words`), but I settled on the one provided by [SCOWL](http://wordlist.aspell.net/). 

SCOWL has conveniently split the words into different files based on how obscure the word is and whether it comes from American, Australian, British or Canadian English. It also separates out abbreviations, contractions, words which contain an upper-case letter that you might still expect to see in a dictionary and the other proper names. The version of English to consider is pretty easy: I'm exclusively going to be solving crosswords in British English and probably the only one to actually use the solver, so it makes the most sense to go with British English for now. While we need a large enough list to include the more obscure words in a crossword, we don't want to overload the user with too many invalid words and hide the actual solutions, so I settled on including the words up to 80 on the size/obscurity scale and ignoring abbreviations. Abbreviations don't seem particularly important in a crossword solver as it seems unlikely you want an anagram to give you an abbreviation, and it wouldn't surprise me if almost every sequence of (say) 3 letters or less is an abbreviation. 

With this many choices, one might want to let the user customise the lists they use, and the complete list of words is only around 8MB uncompressed -- not too costly. Unfortunately, it takes over a minute to process the word lists (on my laptop) and build the data structure we use to search them, too long to expect a user to wait, even if we only do it once. We instead ship pre-compiled versions of the data structures and load them at runtime. It would actually be quite a bit of effort to let the user choose, there's building the data structures, making sure the client only pulls the necessary copies, loading the correct copies, making a settings UI etc. and the big app stores tell me users don't want choice anyway -- for now the user will just have to put up with the selection I've made.

Many simple "English" words have accents, but these are generally ignored when solving a crossword. How should an anagram solver handle these? The obvious solution is to simply ignore them. Unsurprisingly, I'm not a huge fan of giving answers that I know are incorrect, so I decided to keep a record of which words were stripped of accents and apostrophes so I could add them back in later. This does however mean there are a lot of words ending in `S` which are there both with and without an apostrophe. At  some point I might decide to only display one of the versions, but this will do for now.

### Searching a word list

The simplest solution to checking a word list is to simply iterate through all the words and check if they match the criteria we have set, and realistically this is probably quick enough for most use cases. But this is far too boring and we'll try something better. One might try sorting the letters in each word ahead of time, sorting the input and looking up the words where they match, but this doesn't play nicely with blanks or with searching for patterns. Instead, we'll use everyone's favourite search method, depth-first search (DFS). 

To apply DFS to searching a list of words, we first need to view our list as some kind of tree structure. We'll start with a root node with the empty string and at each step into the tree, we add a valid letter to our string, or a special terminating character to signify the end of a word. The words in our list are then the leaves of our tree. This now presents our first problem: How do we know the valid letters we can add at each stage? There are two options: either we explicitly build our tree structure ahead of time or we somehow jump around the list of words and work out the edges from a vertex when we reach that vertex. 

Working out the edges from a vertex when we need them is certainly possible, but it isn't particularly clean and doesn't sound particularly performant either. On the other hand, building our tree structure ahead of time makes it quick and simple to implement a DFS, but the tree is going to be fairly large. We can a bit better than just using a tree: we can easily keep track of the edges we used to reach a given node and so we (mostly) only care about where we can go from this node. If we only had the (contradictory) words `FUN` and `RUN` in our dictionary, we wouldn't need to have separate nodes for `RU` and `FU` as in both cases we can only add the letter `N`. This leads us to a [Directed Acyclic Word Graph](https://en.wikipedia.org/wiki/Deterministic_acyclic_finite_state_automaton), a more efficient way of storing a list of words. As mentioned above, I stripped out the accent information when building the DAWG and so this needs to be added back in at some point going to want. Since there are multiple ways of reaching a single node, we can't simply assign an ID to a node, but it isn't hard to get round this. With a little bit of bookkeeping, we can keep track of the index of the node if we expanded the DAWG to a tree.

Although having implemented a DAWG and DFS, I didn't particularly want to not use it, it is probably worth checking that this a reasonable approach. I cobbled together a quick benchmark on a list of 113,809 words which were already normalised into ASCII and checked how long it took to search for a few strings. Of course, these aren't necessarily the kind of searches (hopefully no one is checking for anagrams of `STOP`) one might make in the real-world, but they'll do for this quick check. It seems that the DFS is quicker in most instances, but slows down considerably in words with blanks. In fact, it is quicker to check every word than to use the DFS when searching for `ALERT???` (where the character `?` is a blank which matches all characters). The code iterating over all words does a quick check that the length of the word matches the length of the anagram and removing this means it takes increases the time taken to 4529192ns i.e. it makes a big difference. One disadvantage of the DAWG approach is that all the words are in the same DAWG and the code might spend a long time in the early few layers of the tree using the three blanks. In the actual website, I split the words into DAWGS based on the length of the word. This is really to make it easier to implement some of the other features later, but it should improve the speed of these searches as well.

<div markdown="1" class="centering">

|    Letters | Number of solutions | Iterate over all words | Depth-first search |
|-----------:|--------------------:|-----------------------:|-------------------:|
|     `STOP` |                   6 |              130543 ns |            5696 ns |
| `ALERTING` |                   6 |              885194 ns |          100567 ns |
|       `??` |                  85 |               65202 ns |           51927 ns |
|     `?TOP` |                  17 |              142090 ns |           39108 ns |
| `ALERT???` |                 424 |             1519711 ns |         2669077 ns |

</div>

### The UI

The UI is fairly standard and written in simple HTML and CSS without any fancy framework, much like the rest of this website. Although it seems like a query might be the perfect use of a custom element, they are so simple I decided to create the structure directly in Javascript every time you press the relevant button. I did, however, make use of a couple of the other modern features added to the browser to make it into a proper PWA. Although I'd heard of all of these features, I'd never actually used them before and this was a good learning exercise.

- **Web workers**: Searching the DAWG takes milliseconds so we could probably get away with running it on the main thread, but the initialisation takes a little longer and very noticeably slows the page down. Instead, we will use a web worker to run the computations on a second thread, send the queries across to it and receive back the results.
- **History API**: We only need to use a single page for this simple app, but being able to navigate between the searches in the history is useful. It might also be useful to send someone a link with the fields already filled in.
- **Service workers**: For such a simple app, making it work offline is as simple as using a service worker to cache every resource and checking the cache when a resource is requested.


## The future

Although I hope that this app is simple enough to be "finished" at some point, I decided that it was worth getting the MVP out the door and leaving a few things on the roadmap. 

- Store the DAWGs etc. as separate files which are passed into the WASM module.
  - This might allow for some limited options on which words are included.
- Add more constraints such as subwords, superwords and contains.
- Add a convenient way to look up definitions, even if this is just a button to do a Google search for you.
- Work out why the PWA doesn't work on iOS/Safari; it is just a little difficult without an iPhone and a Mac.

</div>
