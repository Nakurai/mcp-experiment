# Session2: The demo server

As we saw in [session 1](/docs/session1.md), the whole purpose of an MCP server is to bridge an LLM with a tool that it can use.

Let's build a trivial version of such a tool. We will make a tool that can receive messages from an outside source and send messages to people.

Once this experiment is done, the LLM should be able to talk to this tool via an MCP server, summarize our messages, suggest answers and send messages to people.

This demo server will also include a naive frontend, as if it was a regular website.

## Endpoints
The tool will be a demo server. Here are the endpoints:

- POST /api/login: initiate the Github Oauth web flow
    Once the flow is done, then the user will be created if they do not exist in our database yet, and a jwt will be generated to be used on the frontend.

- POST /messages: content,recipient,date / headers: jwt

- GET /messages: date / headers: access token
- GET /search: content, recipient, date / headers: access token

## Implementation details
To know more about how the server works, you can check out the [README](/demo-server/README.md) file of the demo server. In the context of this session, the high overview of the server's features is enough.