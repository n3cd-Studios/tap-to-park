# Documentation

The documentation for routes can be generated with `npm run docs`

It can be accessed at [localhost:8080/docs/index.html](localhost:8080/docs/index.html).

### Documenting a route

Here is an example of a documentation comment:
```Go
// GetInvites godoc
//
// @Summary        Get the invites for your organization
// @Description    Get the invites associated with a User's organization based on their Bearer token
// @Tags        organization,invite
// @Accept        json
// @Produce        json
// @Param        size  query        number    false    "The size of a page"
// @Param        page  query        number    false    "The page"
// @Success        200    {array} []database.Invite
// @Failure        404 {string} string    "No invites were found for your organization."
// @Failure        500    {string} string    "Couldn't count all of the invites in the organization."
// @Failure        401    {string} string "Unauthorized."
// @Router        /organization/invites [get]
// @Security     BearerToken
func (*OrganizationRoutes) GetInvites(c *gin.Context)
```
The first line is just the name of the function `GetInvites` followed by `godoc`

The `@Summary` and `@Description` lines are just things to describe the route, summary is just the description but less in depth

The `@Tags` are a comma seperated list of metadata tags we use for the routes: organization, invite, spot, stripe, reservation, etc.

The `@Accept` and `@Produce` will most likely always be json, don't really have to worry about changing these.

The `@Param` is really important, it explains how the route is interacted with:
- For path params (`/route/to/{id}`, we would specify that like `@Param id path string true "The ID"`
- For query strings (`/route/to?hey=this`), we specify it like `@Param hey query string false "The hey param after the question"`
- Most of the others are located in the swaggo docs

The `@Success` is just what you return, for most of the routes, you'll specify an `{object}` or an `{array}` like the route above

The `@Failure` consists of a status code, the type and then an example of that type: `@Failure 401 {string} string "Unauthorized."`

The `@Router` might be the most important, just specifying the route and route type

The `@Security BearerToken` is for any protected route