// Copyright (C)2024 LordOfScripts

/*
Package bucketfs implements a Bit-Bucket Virtual File System which serves as
a hybrid bridge between vfs/dummy and vfs/memfs. It is much more suitable for
testing as well as for implementing Dry-Run capabilities than the default
DummyFS implementation.

# Use Case(s)

Let's say you wrote an application like WipeChromium which helps cleaning up
the immense amounts of information your Internet browser keeps. This is
information that not only occupies a lot of space (and grows uncontrollably)
but also may contain privacy-sensitive information.

During development you may have introduced bugs, imagine what would happen
if suddently due to a typo your application wipes out your entire HOME directory!
It is disastrous, right?

Or perhaps it is already well tested (deployment phase) but your end-user may
not be confident about what the application would do in his/her system or
her/his setup. Maybe they would feel at ease by first doing a Dry Run and see
console notifications about each of those sensitive filesystem operations like
Remove, RemoveAll, Mkdir, Rename, Symlink, ReadDir, etc.

## Advantages over DummyFS

The default DummyFS operates by throwing a user-specified error (or nil) on
every single operation. The filesystem itself does nothing else, not traceability.
The DummyFS does not fulfill the use case requirements.

* Is aware of Symbolic links, files & directories.

## Advantages over MemFS

With MemFS you have to replicate into memory. Each of those files and directories
may occupy significant space in memory. Now, I run on a Raspberry Pi with only
1GB of RAM, it is important to me.

Sometimes you only need to state that the filesystem object is there without
actually storing any data other than the name. The MemFS is an overkill for
the use case above.

# How it Works

As stated above, the BitBucketFS is a hybrid between a simplistic DummyFS and
a quite functional MemFS. This filesystem informs you what it is doing and
has the capability of pretending everything is fine, or being selective about
what is fine or not. It has two modes of operation: Silent & Hybrid.

## Silent Mode

In Silent Mode any supported filesystem operation would result in an informational
message printed on the console without producing any error whatsoever, regardless
of whether the filesystem object (file or directory) exists or not. You create
a VFS in Silent Mode like this:

	fs := bucketfs.Create()
	fs.Remove("/home/Documents")

And an operation like the one we used after creation will print:

	⚡ Remove /home/Documents

But your precious directory is still there! All operations in this mode do just
that, they print an informational message and pretend all is good. This is
quite useful for a Dry Run.

In summary, for silent mode ALL operations:

  - print the name of the operation and the primary parameter
  - return nil (no error) and possibly a dummy descriptor.
  - on Read*() operations nothing is actually read
  - on Write() operations nothing is actually written, only file size grows.
  - after Close() is called on a file descriptor, nothing is persisted; therefore,
    the size will be zero again.

## Hybrid Mode

Maybe you need something approaching more like the real OS but without the
overhead of MemFS. You create a BitBucketFS in hybrid mode like this:

	ErrBadThingHappened := errors.New("Bad thing happened!")
	fs := bucketfs.CreateWithError(ErrBadThingHappened)

Or alternatively, you can do it this way too with fluent API:

	ErrBadThingHappened := errors.New("Bad thing happened!")
	fs := bucketfs.Create().WithError(ErrBadThingHappened)

At this point the BitBucket filesystem is empty, it knows no files or directories.
Therefore any Rename, Remove*, Symlink, Open, ReadDir operation will result in
either an os.PathError or os.LinkError with a wrapped ErrBadThingHappened.

So, let's say you want it to recognized certain FAKE files and/or directories
as existent. Internally only a map with these names are stored.

	fs.WithDirectories([]string{"/home/.cache", "/home/Documents"})
	fs.WithFiles([]string{"/home/Documents/test.doc", /home/Documents/test1.txt"})

Those two methods are also part of the fluent API; therefore you could do
this too:

	ErrBadThingHappened := errors.New("Bad thing happened!")
	fs := bucketfs.Create().
	         WithError(ErrBadThingHappened).
	         WithDirectories([]string{...}).
	         WithFiles([]string{...})

Now that we have a BitBucketFS in hybrid mode with a few FAKE files & directories,
the supported filesystem operations would depend on whether the object exists
as FAKE. For example:

	fs.Mkdir("/home/Downloads")

would NOT produce an error and will add the named directory to the FAKE list.
However, this other operation:

	fs.Mkdir("/home/Documents")

would produce an os.PathError with a wrapped bucketfs.ErrDirExists error.

There are two versions of a BitBucketFS file node:

  - Minimal: is only aware whether the node is symbolic link, file or directory
  - Functional: is fully aware of the os.FileMode including (but not enforced)
    (Unix) permissions and size.
*/
package bucketfs
