# VFS for GoLang

![Build](https://github.com/lordofscripts/vfs/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/lordofscripts/vfs?style=flat-square)](https://goreportcard.com/report/github.com/lordofscripts/vfs)
[![Coverage](https://coveralls.io/repos/github/lordofscripts/vfs/badge.svg?branch=main)](https://coveralls.io/github/lordofscripts/vfs?branch=main)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lordofscripts/vfs)

`vfs` is GO library to support *Virtual Filesystems*. It provides the basic
abstractions of filesystems and implementations, like:

* `OS` accessing the file system of the underlying OS,
* `memfs` a full filesystem in-memory, and
* `dummy` which does nothing other than outputting what file operation was
  called without actually modifiying the underlying file system.

## What's Different?

You may have noticed this is a *forked* repository. I forked it from
[3JoB/vfs](https://github.com/3JoB/vfs) which in turn is an improved fork of
the original `blang/vfs` by [Benedikt Lang)](https://github.com/blang/vfs).

I originally used BLang's version `v1.0.0`  and was satisfied with it, although I had
to write some extra code to accomplish what I needed. I realized I needed
BLang's **Dummy** File System but improved to meet my requirements. Unfortunately,
after submitting several issues to the original repository, no answer came of
it. In fact Benedikt's repository has not been updated in 9 years! But it is
still quite useful in its simplicity!

After testing my own shell object to emulate a Dummy Filesystem, I realized it
was better to simply enhance his original `DummyFS`. That's how I came across
**3JoB's** clone tagged `v1.0.0` which has some enhancements over Benedikt's version:

* Support for Symbolic Links
* Minor changes like using `any` instead of `interface{}`

Therefore, I decided to build upon this one instead. After all, 3JoB's version
was updated last year (2023).

Is this a YAUF (Yet-Another-Useless-Fork)? well, no! I plan on making certain
enhancements that would make it suitable for my application out-of-the-box
without the need for glue structures. So, Keep tuned! But to start with:

* Updated it to use `main` as branch instead of the deprecated `master`
* Added `go.mod`
* Included a GO workflow.
* Has a flexible BitBucket Filesystem `bucketfs` more suitable for testing

## Usage

```bash
$ go get github.com/lordofscripts/vfs
```
Note: Always vendor your dependencies or fix on a specific version tag.

```go
import github.com/lordofscripts/vfs
```

```go
// Create a vfs accessing the filesystem of the underlying OS
var osfs vfs.Filesystem = vfs.OS()
osfs.Mkdir("/tmp", 0777)

// Make the filesystem read-only:
osfs = vfs.ReadOnly(osfs) // Simply wrap filesystems to change its behaviour

// os.O_CREATE will fail and return vfs.ErrReadOnly
// os.O_RDWR is supported but Write(..) on the file is disabled
f, _ := osfs.OpenFile("/tmp/example.txt", os.O_RDWR, 0)

// Return vfs.ErrReadOnly
_, err := f.Write([]byte("Write on readonly fs?"))
if err != nil {
    fmt.Errorf("Filesystem is read only!\n")
}

// Create a fully writable filesystem in memory
mfs := memfs.Create()
mfs.Mkdir("/root", 0777)

// Create a vfs supporting mounts
// The root fs is accessing the filesystem of the underlying OS
fs := mountfs.Create(osfs)

// Mount a memfs inside /memfs
// /memfs may not exist
fs.Mount(mfs, "/memfs")

// This will create /testdir inside the memfs
fs.Mkdir("/memfs/testdir", 0777)

// This would create /tmp/testdir inside your OS fs
// But the rootfs `osfs` is read-only
fs.Mkdir("/tmp/testdir", 0777)

// Now use a BitBucket Filesystem in Silent mode
fsb1 := bucketfs.Create()
fsb1.Mkdir("/bucket/testdir", 0777))

// Or an extended BitBucket Filesystem
ErrCustom := errors.New("A BitBucket error"))
fsb2 := bucketfs.CreateWithError(ErrCustom).
         WithFakeDirs([]string{"/bucket/Dir1", "/bucket/Dir2"}).
         WithFakeFiles([]string{"/bucket/Dir1/test.doc", "/bucket/Dir2/test.pdf"})
entries, err := fsb2.ReadDir("/bucket")
```

Check detailed examples below. Also check the [GoDocs](http://godoc.org/github.com/lordofscripts/vfs).

## Why should I use this lib?

- Only Stdlib
- (Nearly) Fully tested (Coverage >87%)
- Easy to create your own filesystem
- Mock a full filesystem for testing (or use included `memfs` or `bucketfs`)
- Compose/Wrap Filesystems `ReadOnly(OS())` and write simple Wrappers
- Many features, see [GoDocs](http://godoc.org/github.com/lordofscripts/vfs) and examples below
- Flexible BitBucket filesystem

## Features and Examples

- [OS Filesystem support](http://godoc.org/github.com/lordofscripts/vfs#example-OsFS)
- [ReadOnly Wrapper](http://godoc.org/github.com/lordofscripts/vfs#example-RoFS)
- [DummyFS for quick mocking](http://godoc.org/github.com/lordofscripts/vfs#example-DummyFS)
- [MemFS - full in-memory filesystem](http://godoc.org/github.com/lordofscripts/vfs/memfs#example-MemFS)
- [MountFS - support mounts across filesystems](http://godoc.org/github.com/lordofscripts/vfs/mountfs#example-MountFS)

### Current state: RELEASE

While the functionality is quite stable and heavily tested, interfaces are subject to change.

    You need more/less abstraction? Let me know by creating a Issue, thank you.

### Motivation

The original author [Benedikt Lang](https://github.com/blang) wrote:

> I simply couldn't find any lib supporting this wide range of variation and adaptability.

And I (*LordOfScripts*) share his thoughts. In fact, I evaluated several similar
GO libraries but many were too bloated and included other appendages I was not
interested in. I loved this VFS version because it had no other dependencies.


### Contribution

Feel free to make a pull request. For bigger changes create a issue first to discuss about it.

### License

See [LICENSE](LICENSE) file.
