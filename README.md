Testify - Thou shalt write tests
================================

Go code (golang) set of packages that provide many tools for testifying that your code will behave as you intend.

  * Easy assertions
  * Mocking
  * HTTP response trapping

Read the API documentation: http://go.pkgdoc.org/github.com/stretchrcom/tetify

Assert package
--------------

The `assert` package provides some helpful methods that allow you to write better test code in Go.

Some examples:

   func TestSomething(t *testing.T) {
   
     assert.Equal(t, 123, 123, "they should be equal")
     assert.NotNil(t, object)
     assert.Nil(t, object)

   }