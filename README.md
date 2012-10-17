Testify - Thou Shalt Write Tests
================================

Go code (golang) set of packages that provide many tools for testifying that your code will behave as you intend.

  * Easy assertions
  * Mocking
  * HTTP response trapping

Read the API Documentation http://go.pkgdoc.org/github.com/stretchrcom/testify

For an introduction to writing test code in Go, see http://golang.org/doc/code.html#Testing

`assert` package
----------------

The `assert` package provides some helpful methods that allow you to write better test code in Go.

  * Prints friendly, easy to read failure descriptions
  * Allows for very readable code
  * Optionally annotate each assertion with a message

See it in action:

    func TestSomething(t *testing.T) {
   
   	  // assert equality
      assert.Equal(t, 123, 123, "they should be equal")

      // assert inequality
      assert.NotEqual(t, 123, 456, "they should not be equal")

      // assert for nil (good for errors)
      assert.Nil(t, object)

      // assert for not nil (good when you expect something)
      if assert.NotNil(t, object) {

      	// now we know that object isn't nil, we are safe to make
      	// further assertions without causing any errors
        assert.Equal(t, "Something", object.Value)

      }

    }

  * Every assert func takes the `testing.T` object as the first argument.  This is how it writes the errors out through the normal `go test` capabilities.
  * Every assert func returns a bool indicating whether the assertion was successful or not, this is useful for if you want to go on making further assertions under certain conditions.

`http` package
--------------

The `http` package contains test objects useful for testing code that relies on the `net/http` package.

`mock` package
--------------

The `mock` package provides a mechanism for easily writing mock objects that can be used in place of real objects when writing test code.

For more information on how to write mock code, check out the API documentation for the `mock` package.

------

Contributing
============

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please also submit a complete test function that demonstrates the issue.


Licence
=======
Copyright (c) 2012 Mat Ryer and Tyler Bunnell

Please consider promoting this project if you find it useful.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.