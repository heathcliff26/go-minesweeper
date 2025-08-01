[![CI](https://github.com/heathcliff26/go-minesweeper/actions/workflows/ci.yaml/badge.svg?event=push)](https://github.com/heathcliff26/go-minesweeper/actions/workflows/ci.yaml)
[![Coverage Status](https://coveralls.io/repos/github/heathcliff26/go-minesweeper/badge.svg)](https://coveralls.io/github/heathcliff26/go-minesweeper)
[![Editorconfig Check](https://github.com/heathcliff26/go-minesweeper/actions/workflows/editorconfig-check.yaml/badge.svg?event=push)](https://github.com/heathcliff26/go-minesweeper/actions/workflows/editorconfig-check.yaml)
[![Generate go test cover report](https://github.com/heathcliff26/go-minesweeper/actions/workflows/go-testcover-report.yaml/badge.svg)](https://github.com/heathcliff26/go-minesweeper/actions/workflows/go-testcover-report.yaml)
[![Renovate](https://github.com/heathcliff26/go-minesweeper/actions/workflows/renovate.yaml/badge.svg)](https://github.com/heathcliff26/go-minesweeper/actions/workflows/renovate.yaml)

# Golang minesweeper

This is an implementation of minesweeper in golang, made with the ui framework fyne.io

![](img/screenshots/difficulty-expert-dark.png#gh-dark-mode-only)
![](img/screenshots/difficulty-expert-light.png#gh-light-mode-only)

## Table of Contents

- [Golang minesweeper](#golang-minesweeper)
  - [Table of Contents](#table-of-contents)
  - [Usage](#usage)
    - [Controls](#controls)
    - [Game options](#game-options)
    - [Changing difficulty](#changing-difficulty)
  - [Installation](#installation)
    - [Download binary](#download-binary)
      - [Uninstalling](#uninstalling)
  - [Screenshots](#screenshots)
  - [Potential features](#potential-features)

## Usage

Or otherwise called how to play the Game.

Tip: The Game starts automatically when revealing the first field.

### Controls

To start a new game, you can click on the smiley between the mine count and the timer.

| Action                                        | control                             |
| --------------------------------------------- | ----------------------------------- |
| Reveal a field                                | left mouse button                   |
| Flag a field                                  | right mouse button                  |
| Reveal all none-flagged fields around a field | double click with left mouse button |

### Game options

The following options can be found under the Game menu

| Option | Description                                                                        |
| ------ | ---------------------------------------------------------------------------------- |
| New    | Start a new game, same as clicking the smiley between the mine count and the timer |
| Replay | Replay the current game                                                            |
| Quit   | Close the app                                                                      |

### Changing difficulty

To change the difficulty, select a new difficulty in the menu.

The custom option will open a dialog where you can create a custom difficulty.
It will tell you if your options don't work.

## Installation

### Download binary

1. Download the [latest release](https://github.com/heathcliff26/go-minesweeper/releases/latest)
2. Unpack the archive
3. Install the app for your user by running:
   - You can install it globally by running the script with `sudo`
```bash
./install.sh -i
```

#### Uninstalling

1. Switch to the folder where you have the installation script
2. Uninstall by running:
   - Run as `sudo` if you installed it globally
```bash
./install.sh -u
```
3. Delete the folder.

## Screenshots

![](img/screenshots/difficulty-beginner-dark.png#gh-dark-mode-only)
![](img/screenshots/difficulty-intermediate-dark.png#gh-dark-mode-only)
![](img/screenshots/difficulty-expert-dark.png#gh-dark-mode-only)
![](img/screenshots/difficulty-beginner-light.png#gh-light-mode-only)
![](img/screenshots/difficulty-intermediate-light.png#gh-light-mode-only)
![](img/screenshots/difficulty-expert-light.png#gh-light-mode-only)

## Potential features

1. Online play by saving highscores on a webserver
2. Compiling to web
3. Save a specific game to be able to play again
4. Create highscore list
