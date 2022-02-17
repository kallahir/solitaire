# Solitaire

## Overview
Solitaire game written in Golang using [SDL2 binding for Go](https://github.com/veandco/go-sdl2).

## Current Status 

<p align="center">
<img src="resources/examples/example01.gif" align="center" hspace="5" vspace="5">
</p>

## Running Unit Tests

Run unit tests for the entire project and output code coverage details:
```
go test -v -coverpkg=./... -cover -coverprofile=c.out -covermode=atomic ./...
```

Generate HTML visualization from code coverate output:
```
go tool cover -html=c.out -o coverage.html
```

Detailed list of code coverage per method and total coverage:
```
go tool cover -func c.out
```

## Copyright Attribution

[cardset-trumps.zip](https://opengameart.org/sites/default/files/cardset-trumps.zip) card set is used under this project in the `resources/cards` folder is based on cards by **[Nora Shishi](http://noragames.com/)**, along with additions by **rh0**, **usr_share**, and **Kumpu** licensed according to the terms of the Creative Commons 3.0 Attribution license [CC-BY 3.0](http://creativecommons.org/licenses/by/3.0/): http://opengameart.org/content/dice-trumps / Not all images from [cardset-trumps.zip](https://opengameart.org/sites/default/files/cardset-trumps.zip) have been used and **shade.gif** file has been renamed to **empty.gif**.
