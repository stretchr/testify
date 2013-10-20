Here is a partial port of python 3 difflib package. Its main goal was to
make available a unified diff implementation, mostly for testing purposes.

The following class and functions have be ported:

* `SequenceMatcher`
* `unified_diff()`

Related doctests have been ported as well.

I have barely used to code yet so do not consider it being production-ready.
The API is likely to evolve too.
