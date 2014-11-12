package parser

import "testing"

func TestParseBoxr(t *testing.T) {
	boxrfile := `
  box: boxr/scratch
  name: test_app
  version: 0.0.1a
  services: 
    - boxr/redis
  build:
    - npm-install
  test:
    - npm test
  run: npm start
  `

	boxr, err := ParseBoxr(boxrfile)

	if err != nil {
		t.Errorf("%#v", err)
	}

	t.Logf("%#v", boxr)
}
