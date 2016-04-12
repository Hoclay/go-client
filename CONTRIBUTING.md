# Contributing

Report issues via Github Issues or [Contact Us](https://handwriting.io/contact/)

If you'd like to make changes to this code, submit a Pull Request.

## Running Tests

The command to run all tests, vet and lint is:

    HANDWRITING_API_URL="https://key:secret@api.handwriting.io" TMPDIR="/tmp" make test

You will have to replace the key and secret in `HANDWRITING_API_URL` with your own.  Test tokens will work fine for this.

Any example that writes a file should use `os.TempDir()`.  After running tests you can inspect the files you've written in your `TMPDIR` directory.
