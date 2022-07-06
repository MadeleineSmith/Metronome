# Metronome

https://user-images.githubusercontent.com/16705876/177578458-711f4562-ecc1-4951-91d3-c7084bf8eb9b.mov


## Background
* I wanted to build out a **metronome** on the command line so I didn't have to rely on Google's metronome
* The metronome logic is taken from [Rosetta Code's metronome solution](https://rosettacode.org/wiki/Metronome#Go), and I used the Go [beep](https://github.com/faiface/beep) package to output sound
* I created a [Homebrew tap](https://github.com/MadeleineSmith/homebrew-metronome) for [easy installation](#first-install-instructions) of the metronome package

---

## First install instructions:
* ` brew tap madeleinesmith/metronome && brew install metronome `

---

## How to release new version and update local package:
Releasing new version:
* Tag code using ` git tag -a v0.0.2 -m "version 0.2.0" `
* Push tag with ` git push origin v0.2.0 `
* Create a new release on GitHub for that tag
* Copy the link of the `tar.gz` file on GitHub
* Change the `url` line of `homebrew-metronome` [repo](https://github.com/MadeleineSmith/homebrew-metronome/blob/4661e8c8d8ef9dcafb2a46e645d57550990ba31b/metronome.rb#L7) to be this
* And also update the `sha256` line by running `shasum -a 256 xxxxxxx.tar.gz` on the downloaded tar file (above)  
* Commit and push the `homebrew-metronome` repo with these edits

Updating the local package:
* Run ``` brew update && brew upgrade metronome ```

---

## Tutorials referenced whilst building:
Creating a Homebrew tap:
* https://betterprogramming.pub/a-step-by-step-guide-to-create-homebrew-taps-from-github-repos-f33d3755ba74
* https://flowerinthenight.com/blog/2019/07/30/homebrew-golang

Using go:embed:
* https://blog.jetbrains.com/go/2021/06/09/how-to-use-go-embed-in-go-1-16/
