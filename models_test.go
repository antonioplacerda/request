package requests_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/antonioplacerda/requests"
)

func TestFloat(t *testing.T) {
	c := qt.New(t)

	c.Run("OK - string int", func(c *qt.C) {
		var res requests.Float
		err := res.UnmarshalJSON([]byte(`"1000"`))
		c.Assert(res, qt.Equals, requests.Float(1000))
		c.Assert(err, qt.IsNil)
	})

	c.Run("OK - string float", func(c *qt.C) {
		var res requests.Float
		err := res.UnmarshalJSON([]byte(`"32.89"`))
		c.Assert(res, qt.Equals, requests.Float(32.89))
		c.Assert(err, qt.IsNil)
	})

	c.Run("OK - int", func(c *qt.C) {
		var res requests.Float
		err := res.UnmarshalJSON([]byte(`1000`))
		c.Assert(res, qt.Equals, requests.Float(1000))
		c.Assert(err, qt.IsNil)
	})

	c.Run("OK - float", func(c *qt.C) {
		var res requests.Float
		err := res.UnmarshalJSON([]byte(`32.89`))
		c.Assert(res, qt.Equals, requests.Float(32.89))
		c.Assert(err, qt.IsNil)
	})

	c.Run("OK - negative", func(c *qt.C) {
		var res requests.Float
		err := res.UnmarshalJSON([]byte(`"-132.89"`))
		c.Assert(res, qt.Equals, requests.Float(-132.89))
		c.Assert(err, qt.IsNil)
	})

	c.Run("NOK - empty", func(c *qt.C) {
		var res requests.Float
		err := res.UnmarshalJSON([]byte(`""`))
		c.Assert(res, qt.Equals, requests.Float(0))
		c.Assert(err, qt.IsNil)
	})

	c.Run("NOK - invalid", func(c *qt.C) {
		var res requests.Float
		err := res.UnmarshalJSON([]byte(`"wops"`))
		c.Assert(res, qt.Equals, requests.Float(0))
		c.Assert(err, qt.ErrorMatches, `strconv.ParseFloat: parsing "wops": invalid syntax`)
	})
}
