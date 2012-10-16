Testify - Thou shalt write tests
================================

Go code (golang) set of packages that provide many tools for testifying that your code will behave as you intend.

  * Easy assertions
  * Mocking
  * HTTP response trapping

Read the API documentation: http://go.pkgdoc.org/github.com/stretchrcom/testify

Assert package
--------------

The `assert` package provides some helpful methods that allow you to write better test code in Go.

See it in action:

    func TestSomething(t *testing.T) {
   
      assert.Equal(t, 123, 123, "they should be equal")

      assert.Nil(t, object)

      if assert.NotNil(t, object) {

        assert.Equal(t, "Something", object.Value)

      }

    }

  * Every assert func takes the `testing.T` object as the first argument.  This is how it writes the errors out through the normal `go test` capabilities.
  * Every assert func returns a bool indicating whether the assertion was successful or not, this is useful for if you want to go on making further assertions under certain conditions.