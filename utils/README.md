testify/utils
=============

gocheck2testify.sed
-------------------

This SED script converts "gocheck" style unit tests to "testify" style.

IMPORTANT: This assumes that your tests are already organized within test suites.

Assuming your test programs all end in \_test.go, then the command to use is:

    sed -r -ibak -f ../../stretchr/testify/utils/gocheck2testify.sed *_test.go

Oh, and be sure to uncomment one of the last two lines of the script if you prefer to use "this" or "s" for the suite instance variable instead of "suite".

Questions? Comments? Problems? Include @polyglot-jones in your GitHub issue/comment.