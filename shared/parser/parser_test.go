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
    steps:
      - npm-install
  test:
    steps:
      - npm test
  deploy:
    steps:
      - npm start
  `

	boxr, err := ParseBoxr(boxrfile)

	if err != nil {
		t.Errorf("%#v", err)
	}

	t.Logf("%#v", boxr)
}
