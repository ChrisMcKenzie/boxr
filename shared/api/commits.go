package api

import "github.com/Secret-Ironman/boxr/Godeps/_workspace/src/github.com/gin-gonic/gin"

func CommitsHookGet(c *gin.Context) {
	// Queue Process:
	// ==============
	// get information from sqlite about id
	// get changes from hook
	// pull repo
	// read boxr.(yaml|json)
	// build services
	// run build steps on defined box
	// run test steps on box
	// if pass run deploy steps
	// else notify failure

	c.String(200, "hook called")
}
