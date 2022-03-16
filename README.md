# Goto

> Shortcut and speedy tool to all your daily links.

`goto` let you create shortcuts to access webpages. 

For example running `goto alias wat
https://www.destroyallsoftware.com/talks/wat` will let you open the
website just by `goto wat`.

And you can also use templates! `goto alias jira
https://jira.company.com/projects/AB/issues/{{ . }}` allows you to
quickly open issues `goto jira/AB-123`.

## Installation

## Usage

```console
# map shortcuts to urls
goto alias gh https://github.com
# potentially with a template
goto alias jira https://jira.company.com/projects/AB/issues/{{ . }}
# then whatever provided after `jira/` will replace `{{ . }}` before
opening the URL

# open an aliased website
goto <alias>

# list existing aliases
goto ls
```

## Development

## Release

Releases are powered by [Goreleaser](https://goreleaser.com/).

### Features

- [ ] Export/Import (something to rebuild the DB if needed)
- [ ] Related: create sharable configurations? Like common service integrations
- [ ] Edit existing entries

### Todo

- [ ] Refactor this mess
- [ ] Add more CLI options, like DB path
- [ ] Proper repository setup and binary release
- [ ] Support profiles (`goto gh` as pro or personal may not take you to the same place)
- [ ] Enforce DB entries uniqueness
- [ ] Make the id shown and row id the same

<br />
<br />

<p align="center">
  <img src="https://raw.github.com/hivetech/hivetech.github.io/master/images/pilotgopher.jpg" alt="gopher" width="200px"/>
</p>


[GoDoc]: https://godoc.org/github.com/hackliff/cliper
[walker]: http://gowalker.org/github.com/hackliff/cliper
[GoDoc Widget]: https://godoc.org/hackliff/cliper?status.svg
[releases]: https://github.com/hackliff/cliper/releases

[semver]: http://semver.org
[commit]: https://chris.beams.io/posts/git-commit/
